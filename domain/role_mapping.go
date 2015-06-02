package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "accounts", Fk: "accountId"},
		{Related: "roles", Fk: "roleId"},
	}

	ModelDirectory.Register(RoleMapping{}, "roleMappings", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type RoleMapping struct {
	GormModel
	AccountID int     `json:"accountId,omitempty" sql:"index"`
	Account   Account `json:"account,omitempty"`
	RoleID    int     `json:"roleId,omitempty" sql:"index"`
	Role      Role    `json:"role,omitempty"`
}

func (m *RoleMapping) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "roleID":
		m.RoleID = id
	case "accountID":
		m.AccountID = id
	}
}

func scopeRoleMapping(m *RoleMapping) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}

func (m *RoleMapping) BeforeRender() error {
	m.Role.BeforeRender()
	m.Account.BeforeRender()
	return nil
}

func (m *RoleMapping) BeforeActionCreate() error {
	scopeRoleMapping(m)
	return nil
}

func (m *RoleMapping) AfterActionCreate() error {
	return nil
}

func (m *RoleMapping) BeforeActionUpdate() error {
	scopeRoleMapping(m)
	return nil
}

func (m *RoleMapping) AfterActionUpdate() error {
	return nil
}

func (m *RoleMapping) BeforeActionDelete() error {
	return nil
}

func (m *RoleMapping) AfterActionDelete() error {
	return nil
}
