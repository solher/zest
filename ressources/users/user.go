package users

import (
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/ressources/emails"
)

func init() {
	domain.Models.Register(User{})
}

//go:generate gen -f
//+gen routes controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" controller_test interactor:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" interactor_test repository:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID"
type User struct {
	domain.GormModel
	ID        int            `json:"id,omitempty" gorm:"primary_key"`
	FirstName string         `json:"firstName,omitempty"`
	LastName  string         `json:"lastName,omitempty"`
	Password  string         `json:"password,omitempty"`
	Emails    []emails.Email `json:"emails,omitempty"`
}
