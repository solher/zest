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

func scopeAclMapping(m *AclMapping) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt
}

func (m *AclMapping) ValidateCreate() error {
	return nil
}

func (m *AclMapping) ValidateUpdate() error {
	return nil
}

func (m *AclMapping) ValidateDelete() error {
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
