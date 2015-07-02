package domain

import "time"

func init() {
	relations := []DBRelation{
		{Related: "accounts", Fk: "accountId"},
	}

	ModelDirectory.Register(Session{}, "sessions", relations)
}

type Session struct {
	GormModel
	AuthToken string    `json:"authToken"`
	IP        string    `json:"ip"`
	Agent     string    `json:"agent"`
	ValidTo   time.Time `json:"validTo"`
	DeletedAt time.Time `json:"deletedAt"`
	AccountID int       `json:"accountId" sql:"index"`
	Account   Account   `json:"account,omitempty"`
}

func (m *Session) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "accountID":
		m.AccountID = id
	}
}

func (m *Session) BeforeRender() {
	m.Account.BeforeRender()
}
