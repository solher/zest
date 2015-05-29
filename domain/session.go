package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "accounts", Fk: "accountId"},
	}

	ModelDirectory.Register(Session{}, "sessions", relations)
}

var SessionRelated *Relation

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type Session struct {
	GormModel
	AccountID int       `json:"accountId,omitempty" sql:"index"`
	Account   Account   `json:"account,omitempty"`
	AuthToken string    `json:"authToken,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Agent     string    `json:"agent,omitempty"`
	ValidTo   time.Time `json:"validTo,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
}

func (m *Session) ScopeModel(accountID int) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
	m.AccountID = accountID
}
