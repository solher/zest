package resources

import (
	"time"

	"github.com/solher/zest/domain"
	"github.com/solher/zest/internalerrors"
	"github.com/solher/zest/usecases"
	"github.com/solher/zest/utils"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	usecases.DependencyDirectory.Register(NewUserGuestInter)
}

type UserGuestInter struct {
	repo      AbstractUserRepo
	userInter AbstractUserInter
}

func NewUserGuestInter(repo AbstractUserRepo, userInter AbstractUserInter) *UserGuestInter {
	return &UserGuestInter{repo: repo, userInter: userInter}
}

func (i *UserGuestInter) scopeModel(user *domain.User) error {
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}
	user.Account = domain.Account{}

	return nil
}

func (i *UserGuestInter) UpdateByID(id int, user *domain.User, context usecases.QueryContext) (*domain.User, error) {
	i.scopeModel(user)

	utils.Breakpoint()

	attributes := map[string]interface{}{
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Email":     user.Email,
	}

	user, err := i.repo.UpdateAttributesByID(id, attributes, context)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (i *UserGuestInter) UpdatePassword(id int, context usecases.QueryContext, oldPassword, newPassword string) (*domain.User, error) {
	user, err := i.repo.FindByID(id, context)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, internalerrors.NotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return nil, internalerrors.InvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 0)
	if err != nil {
		return nil, err
	}

	attributes := map[string]interface{}{
		"Password": string(hashedPassword),
	}

	user, err = i.repo.UpdateAttributesByID(id, attributes, usecases.QueryContext{})
	if err != nil {
		return nil, err
	}

	return user, nil
}
