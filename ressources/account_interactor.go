package ressources

import (
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
