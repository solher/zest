package ressources

import "github.com/Solher/auth-scaffold/domain"

func init() {
	domain.ModelDirectory.Register(Account{})
}

//go:generate gen -f
//+gen routes repository:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID"
type Account struct {
	domain.GormModel
	ID          int       `json:"id,omitempty" gorm:"primary_key"`
	Users       []User    `json:"users,omitempty"`
	Sessions    []Session `json:"sessions,omitempty"`
	IsAdmin     bool      `json:"isAdmin,omitempty"`
	IsActivated bool      `json:"isActivated,omitempty"`
}
