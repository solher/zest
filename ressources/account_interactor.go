package ressources

import (
	"database/sql"
	"time"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/Solher/auth-scaffold/utils"
	"golang.org/x/crypto/bcrypt"
)

type AbstractAccountRepo interface {
	Create(accounts []domain.Account) ([]domain.Account, error)
	CreateOne(account *domain.Account) (*domain.Account, error)
	Find(filter *interfaces.Filter) ([]domain.Account, error)
	FindByID(id int, filter *interfaces.Filter) (*domain.Account, error)
	Upsert(accounts []domain.Account) ([]domain.Account, error)
	UpsertOne(account *domain.Account) (*domain.Account, error)
	DeleteAll(filter *interfaces.Filter) error
	DeleteByID(id int) error
	Raw(query string, values ...interface{}) (*sql.Rows, error)
}

type AccountInter struct {
	repo        AbstractAccountRepo
	userRepo    AbstractUserRepo
	sessionRepo AbstractSessionRepo
}

func NewAccountInter(repo AbstractAccountRepo, userRepo AbstractUserRepo, sessionRepo AbstractSessionRepo) *AccountInter {
	return &AccountInter{repo: repo, userRepo: userRepo, sessionRepo: sessionRepo}
}

func (i *AccountInter) Signin(ip, userAgent string, credentials *Credentials) (*domain.Session, error) {
	filter := &interfaces.Filter{
		Limit: 1,
		Where: map[string]interface{}{"email": credentials.Email},
	}

	users, err := i.userRepo.Find(filter)
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

	return session, nil
}

func (i *AccountInter) Signout(currentSession *domain.Session) error {
	err := i.sessionRepo.DeleteByID(currentSession.ID)
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
	filter := &interfaces.Filter{
		Include: []interface{}{"users"},
	}

	account, err := i.repo.FindByID(currentSession.AccountID, filter)
	if err != nil {
		return nil, err
	}

	currentSession.Account = domain.Account{}
	account.Sessions = []domain.Session{*currentSession}

	return account, nil
}

func (i *AccountInter) CurrentSessionFromToken(authToken string) (*domain.Session, error) {
	filter := &interfaces.Filter{
		Limit:   1,
		Where:   map[string]interface{}{"authToken": authToken},
		Include: []interface{}{"account"},
	}

	sessions, err := i.sessionRepo.Find(filter)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 1 {
		session := sessions[0]

		if session.ValidTo.After(time.Now()) {
			return &session, nil
		}
	}

	return nil, nil
}

func (i *AccountInter) GetGrantedRoles(accountID int, ressource, method string) ([]string, error) {
	var rows *sql.Rows
	var err error
	roleNames := []string{}

	if accountID == 0 {
		rows, err = i.repo.Raw(`
			SELECT DISTINCT roles.name
			FROM roles, acl_mappings, acls
			INNER JOIN acl_mappings AS am ON am.role_id = roles.id AND am.acl_id = acls.id
			WHERE roles.name IN ('Guest', 'Anyone') AND acls.ressource = ? AND acls.method = ?
			`, ressource, method)
	} else {
		rows, err = i.repo.Raw(`
			SELECT DISTINCT roles.name
			FROM roles, acl_mappings, acls
			INNER JOIN acl_mappings AS am ON am.role_id = roles.id AND am.acl_id = acls.id
			WHERE roles.name IN ('Authenticated', 'Anyone') AND acls.ressource = ? AND acls.method = ?
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
			FROM role_mappings, roles, acl_mappings, acls
			INNER JOIN role_mappings AS rm ON rm.role_id = roles.id
			INNER JOIN acl_mappings AS am ON am.role_id = roles.id AND am.acl_id = acls.id
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

	return roleNames, nil
}
