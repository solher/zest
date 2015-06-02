package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	relations := []Relation{
		{Related: "accounts", Fk: "accountId"},
	}

	ModelDirectory.Register(User{}, "users", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type User struct {
	GormModel
	AccountID int    `json:"accountId,omitempty" sql:"index"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Password  string `json:"-"`
	Email     string `json:"email,omitempty" sql:"unique"`
}

func scopeUser(m *User) error {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt

	if m.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), 0)
		if err != nil {
			return err
		}

		m.Password = string(hashedPassword)
	}

	return nil
}

func (m *User) ValidateCreate() error {
	return nil
}

func (m *User) ValidateUpdate() error {
	return nil
}

func (m *User) ValidateDelete() error {
	return nil
}

func (m *User) BeforeActionCreate() error {
	err := scopeUser(m)
	if err != nil {
		return err
	}

	return nil
}

func (m *User) AfterActionCreate() error {
	return nil
}

func (m *User) BeforeActionUpdate() error {
	err := scopeUser(m)
	if err != nil {
		return err
	}

	return nil
}

func (m *User) AfterActionUpdate() error {
	return nil
}

func (m *User) BeforeActionDelete() error {
	return nil
}

func (m *User) AfterActionDelete() error {
	return nil
}
