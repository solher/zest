package ressources

import (
	"database/sql"
	"time"

	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/usecases"

	"github.com/Solher/zest/internalerrors"
	"github.com/Solher/zest/utils"
	"golang.org/x/crypto/bcrypt"
)

type AbstractAccountRepo interface {
	Create(accounts []domain.Account) ([]domain.Account, error)
	CreateOne(account *domain.Account) (*domain.Account, error)
	Find(filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.Account, error)
	FindByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Account, error)
	Update(accounts []domain.Account, filter *usecases.Filter, ownerRelations []domain.Relation) ([]domain.Account, error)
	UpdateByID(id int, account *domain.Account, filter *usecases.Filter, ownerRelations []domain.Relation) (*domain.Account, error)
	DeleteAll(filter *usecases.Filter, ownerRelations []domain.Relation) error
	DeleteByID(id int, filter *usecases.Filter, ownerRelations []domain.Relation) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type AccountInter struct {
	repo                 AbstractAccountRepo
	userRepo             AbstractUserRepo
	sessionRepo          AbstractSessionRepo
	sessionCacheInter    *usecases.SessionCacheInter
	permissionCacheInter *usecases.PermissionCacheInter
}

func NewAccountInter(repo AbstractAccountRepo, userRepo AbstractUserRepo, sessionRepo AbstractSessionRepo,
	sessionCacheInter *usecases.SessionCacheInter, permissionCacheInter *usecases.PermissionCacheInter) *AccountInter {

	return &AccountInter{repo: repo, userRepo: userRepo, sessionRepo: sessionRepo,
		sessionCacheInter: sessionCacheInter, permissionCacheInter: permissionCacheInter}
}

func (i *AccountInter) Signin(ip, userAgent string, credentials *Credentials) (*domain.Session, error) {
	filter := &usecases.Filter{
		Limit: 1,
		Where: map[string]interface{}{"email": credentials.Email},
	}

	users, err := i.userRepo.Find(filter, nil)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, internalerrors.RessourceNotFound
	}
	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return nil, err
	}

	authToken := utils.RandStr(64, "alphanum")

	var validTo time.Time
	if credentials.RememberMe == true {
		validTo = time.Now().Add(365 * 24 * time.Hour)
	} else {
		validTo = time.Now().Add(24 * time.Hour)
	}

	session := &domain.Session{
		AccountID: user.AccountID,
		AuthToken: authToken,
		IP:        ip,
		Agent:     userAgent,
		ValidTo:   validTo,
	}

	session, err = i.sessionRepo.CreateOne(session)
	if err != nil {
		return nil, err
	}

	err = i.sessionCacheInter.Add(session.AuthToken, *session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (i *AccountInter) Signout(currentSession *domain.Session) error {
	authToken := currentSession.AuthToken

	err := i.sessionRepo.DeleteByID(currentSession.ID, nil, nil)
	if err != nil {
		return err
	}

	err = i.sessionCacheInter.Remove(authToken)
	if err != nil {
		return err
	}

	return nil
}
func (i *AccountInter) Signup(user *domain.User) (*domain.Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	account := &domain.Account{
		Users: []domain.User{*user},
	}

	account, err = i.repo.CreateOne(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
func (i *AccountInter) Current(currentSession *domain.Session) (*domain.Account, error) {
	filter := &usecases.Filter{
		Include: []interface{}{"users"},
	}

	account, err := i.repo.FindByID(currentSession.AccountID, filter, nil)
	if err != nil {
		return nil, err
	}

	currentSession.Account = domain.Account{}
	account.Sessions = []domain.Session{*currentSession}

	return account, nil
}

func (i *AccountInter) CurrentSessionFromToken(authToken string) (*domain.Session, error) {
	session, err := i.sessionCacheInter.Get(authToken)

	if err != nil {
		filter := &usecases.Filter{
			Limit: 1,
			Where: map[string]interface{}{"authToken": authToken},
		}

		sessions, err := i.sessionRepo.Find(filter, nil)
		if err != nil {
			return nil, err
		}

		if len(sessions) == 1 {
			session = sessions[0]

			err = i.sessionCacheInter.Add(session.AuthToken, session)
			if err != nil {
				return nil, err
			}
		}
	}

	if session.ValidTo.After(time.Now()) {
		return &session, nil
	}

	return nil, nil
}

func (i *AccountInter) GetGrantedRoles(accountID int, ressource, method string) ([]string, error) {
	var rows *sql.Rows

	roleNames, err := i.permissionCacheInter.GetPermissionRoles(accountID, ressource, method)

	if err != nil {
		if accountID == 0 {
			rows, err = i.repo.Raw(`
				SELECT DISTINCT roles.name
				FROM roles, acls
				INNER JOIN acl_mappings ON acl_mappings.role_id = roles.id AND acl_mappings.acl_id = acls.id
				WHERE roles.name IN ('Guest', 'Anyone') AND acls.ressource = ? AND acls.method = ?
				`, ressource, method)
		} else {
			rows, err = i.repo.Raw(`
				SELECT DISTINCT roles.name
				FROM roles, acls
				INNER JOIN acl_mappings ON acl_mappings.role_id = roles.id AND acl_mappings.acl_id = acls.id
				WHERE roles.name IN ('Authenticated', 'Owner', 'Anyone') AND acls.ressource = ? AND acls.method = ?
				`, ressource, method)
		}

		defer rows.Close()

		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var roleName string
			rows.Scan(&roleName)
			roleNames = append(roleNames, roleName)
		}

		if len(roleNames) == 0 {
			rows, err = i.repo.Raw(`
				SELECT DISTINCT roles.name
				FROM roles, acls
				INNER JOIN role_mappings ON role_mappings.role_id = roles.id
				INNER JOIN acl_mappings ON acl_mappings.role_id = roles.id AND acl_mappings.acl_id = acls.id
				WHERE role_mappings.account_id = ? AND acls.ressource = ? AND acls.method = ?
				`, accountID, ressource, method)
		}

		defer rows.Close()

		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var roleName string
			rows.Scan(&roleName)
			roleNames = append(roleNames, roleName)
		}
	}

	return roleNames, nil
}
