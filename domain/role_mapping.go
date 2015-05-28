package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "accounts", Fk: "accountId"},
		{Related: "roles", Fk: "roleId"},
	}

	ModelDirectory.Register(RoleMapping{}, "roleMappings", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type RoleMapping struct {
	GormModel
	AccountID int     `json:"accountId,omitempty" sql:"index"`
	Account   Account `json:"account,omitempty"`
	RoleID    int     `json:"roleId,omitempty" sql:"index"`
	Role      Role    `json:"role,omitempty"`
}

func (m *RoleMapping) ScopeModel(accountID int) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt

	if accountID != 0 {
		m.AccountID = accountID
	}
}
