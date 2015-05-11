package ressources

import "github.com/Solher/auth-scaffold/domain"

func init() {
	domain.ModelDirectory.Register(Account{})
}

//go:generate gen -f
//+gen repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type Account struct {
	domain.GormModel
	ID          int       `json:"id,omitempty" gorm:"primary_key"`
	Users       []User    `json:"users,omitempty"`
	Sessions    []Session `json:"sessions,omitempty"`
	IsAdmin     bool      `json:"isAdmin,omitempty"`
	IsActivated bool      `json:"isActivated,omitempty"`
}
