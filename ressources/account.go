package ressources

import (
	"time"

	"github.com/Solher/auth-scaffold/domain"
)

func init() {
	domain.ModelDirectory.Register(Account{})
}

//go:generate gen -f
//+gen repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type Account struct {
	domain.GormModel
	Users       []User    `json:"users,omitempty"`
	Sessions    []Session `json:"sessions,omitempty"`
	IsAdmin     bool      `json:"isAdmin,omitempty"`
	IsActivated bool      `json:"isActivated,omitempty"`
}

func (m *Account) ScopeModel() {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}
