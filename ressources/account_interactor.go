package ressources

import (
	"time"

	"github.com/Solher/auth-scaffold/interfaces"
	"github.com/Solher/auth-scaffold/internalerrors"
	"github.com/Solher/auth-scaffold/utils"
	"golang.org/x/crypto/bcrypt"
)

type AbstractAccountRepo interface {
	Create(accounts []Account) ([]Account, error)
	CreateOne(account *Account) (*Account, error)
	Find(filter *interfaces.Filter) ([]Account, error)
	FindByID(id int, filter *interfaces.Filter) (*Account, error)
	Upsert(accounts []Account) ([]Account, error)
	UpsertOne(account *Account) (*Account, error)
	DeleteAll(filter *interfaces.Filter) error
	DeleteByID(id int) error
}

type AccountInter struct {
	repo        AbstractAccountRepo
	userRepo    AbstractUserRepo
	sessionRepo AbstractSessionRepo
}

func NewAccountInter(repo AbstractAccountRepo, userRepo AbstractUserRepo, sessionRepo AbstractSessionRepo) *AccountInter {
	return &AccountInter{repo: repo, userRepo: userRepo, sessionRepo: sessionRepo}
}

func (i *AccountInter) Signin(ip, userAgent string, credentials *Credentials) (*Session, error) {
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

	session := &Session{
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

func (i *AccountInter) Signout(currentSession *Session) error {
	err := i.sessionRepo.DeleteByID(currentSession.ID)
	if err != nil {
		return err
	}

	return nil
}
func (i *AccountInter) Signup(user *User) (*Account, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	account := &Account{
		Users: []User{*user},
	}

	account, err = i.repo.CreateOne(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
func (i *AccountInter) Current(currentSession *Session) (*Account, error) {
	filter := &interfaces.Filter{
		Include: []interface{}{"users"},
	}

	account, err := i.repo.FindByID(currentSession.AccountID, filter)
	if err != nil {
		return nil, err
	}

	account.Sessions = []Session{*currentSession}

	return account, nil
}
