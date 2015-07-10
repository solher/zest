package resources

import (
	"time"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"

	"github.com/solher/zest/internalerrors"
	"github.com/solher/zest/utils"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	usecases.DependencyDirectory.Register(NewAccountGuestInter)
}

type AccountGuestInter struct {
	repo                 AbstractAccountRepo
	userInter            AbstractUserInter
	sessionInter         AbstractSessionInter
	sessionCacheInter    usecases.AbstractSessionCacheInter
	permissionCacheInter usecases.AbstractPermissionCacheInter
}

func NewAccountGuestInter(repo AbstractAccountRepo, userInter AbstractUserInter, sessionInter AbstractSessionInter,
	sessionCacheInter *usecases.SessionCacheInter, permissionCacheInter *usecases.PermissionCacheInter) *AccountGuestInter {

	return &AccountGuestInter{repo: repo, userInter: userInter, sessionInter: sessionInter,
		sessionCacheInter: sessionCacheInter, permissionCacheInter: permissionCacheInter}
}

func (i *AccountGuestInter) Signin(ip, userAgent string, credentials *Credentials) (*domain.Session, error) {
	filter := &usecases.Filter{
		Limit: 1,
		Where: map[string]interface{}{"email": credentials.Email},
	}

	users, err := i.userInter.Find(usecases.QueryContext{Filter: filter})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, internalerrors.NotFound
	}
	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return nil, internalerrors.NotFound
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

	session, err = i.sessionInter.CreateOne(session)
	if err != nil {
		return nil, err
	}

	err = i.sessionCacheInter.Add(session.AuthToken, *session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (i *AccountGuestInter) Signout(currentSession *domain.Session) error {
	authToken := currentSession.AuthToken

	err := i.sessionInter.DeleteByID(currentSession.ID, usecases.QueryContext{})
	if err != nil {
		return err
	}

	err = i.sessionCacheInter.Remove(authToken)
	if err != nil {
		return err
	}

	return nil
}

func (i *AccountGuestInter) Signup(user *domain.User) (*domain.Account, error) {
	account, err := i.repo.CreateOne(&domain.Account{})
	if err != nil {
		return nil, err
	}

	user.AccountID = account.ID

	user, err = i.userInter.CreateOne(user)
	if err != nil {
		return nil, err
	}

	account.Users = []domain.User{*user}

	return account, nil
}

func (i *AccountGuestInter) CurrentSessionFromToken(authToken string) (*domain.Session, error) {
	session, err := i.sessionCacheInter.Get(authToken)

	if err != nil {
		filter := &usecases.Filter{
			Limit: 1,
			Where: map[string]interface{}{"authToken": authToken},
		}

		sessions, err := i.sessionInter.Find(usecases.QueryContext{Filter: filter})
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

func (i *AccountGuestInter) GetGrantedRoles(accountID int, resource, method string) ([]string, error) {
	roleNames, err := i.permissionCacheInter.GetPermissionRoles(accountID, resource, method)
	if err != nil {
		return nil, err
	}

	return roleNames, nil
}
