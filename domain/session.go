package domain

import "time"

func init() {
	relations := []DBRelation{
		{Related: "Accounts", Fk: "AccountId"},
	}

	ModelDirectory.Register(Session{}, "sessions", relations)
}

type Session struct {
	GormModel
	AuthToken string    `json:"authToken,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Agent     string    `json:"agent,omitempty"`
	ValidTo   time.Time `json:"validTo,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
	AccountID int       `json:"AccountId,omitempty" sql:"index"`
	Account   Account   `json:"Account,omitempty"`
}

func (m *Session) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "AccountID":
		m.AccountID = id
	}
}

func (m *Session) BeforeRender() {
	m.Account.BeforeRender()
}
