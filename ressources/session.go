package ressources

import (
	"time"

	"github.com/Solher/auth-scaffold/domain"
)

func init() {
	domain.ModelDirectory.Register(Session{})
}

//go:generate gen -f
//+gen routes controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type Session struct {
	domain.GormModel
	AccountID int       `json:"accountId,omitempty" sql:"index"`
	Account   Account   `json:"account,omitempty"`
	AuthToken string    `json:"authToken,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Agent     string    `json:"agent,omitempty"`
	ValidTo   time.Time `json:"validTo,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
}

func (m *Session) ScopeModel() {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}
