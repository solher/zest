package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "roles", Fk: "roleId"},
		{Related: "acls", Fk: "aclId"},
	}

	ModelDirectory.Register(AclMapping{}, "aclMappings", relations)
}

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type AclMapping struct {
	GormModel
	RoleID int  `json:"roleId,omitempty" sql:"index"`
	Role   Role `json:"role,omitempty"`
	AclID  int  `json:"aclId,omitempty" sql:"index"`
	Acl    Acl  `json:"acl,omitempty"`
}

func (m *AclMapping) SetRelatedID(idKey string, id int) {
	switch idKey {
	case "roleID":
		m.RoleID = id
	case "aclID":
		m.AclID = id
	}
}

func scopeAclMapping(m *AclMapping) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}

func (m *AclMapping) BeforeRender() error {
	m.Role.BeforeRender()
	m.Acl.BeforeRender()
	return nil
}

func (m *AclMapping) BeforeActionCreate() error {
	scopeAclMapping(m)
	return nil
}

func (m *AclMapping) AfterActionCreate() error {
	return nil
}

func (m *AclMapping) BeforeActionUpdate() error {
	scopeAclMapping(m)
	return nil
}

func (m *AclMapping) AfterActionUpdate() error {
	return nil
}

func (m *AclMapping) BeforeActionDelete() error {
	return nil
}

func (m *AclMapping) AfterActionDelete() error {
	return nil
}
