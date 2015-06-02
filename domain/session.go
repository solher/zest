package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "accounts", Fk: "accountId"},
	}

	ModelDirectory.Register(Session{}, "sessions", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
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

func scopeSession(m *Session) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
	m.DeletedAt = m.CreatedAt
}

func (m *Session) ValidateCreate() error {
	return nil
}

func (m *Session) ValidateUpdate() error {
	return nil
}

func (m *Session) ValidateDelete() error {
	return nil
}

func (m *Session) BeforeCreate() error {
	scopeSession(m)
	return nil
}

func (m *Session) AfterCreate() error {
	return nil
}

func (m *Session) BeforeUpdate() error {
	scopeSession(m)
	return nil
}

func (m *Session) AfterUpdate() error {
	return nil
}

func (m *Session) BeforeDelete() error {
	return nil
}

func (m *Session) AfterDelete() error {
	return nil
}
