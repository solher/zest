package domain

import "time"

func init() {
	relations := []Relation{
		{Related: "roles", Fk: "roleId"},
		{Related: "acls", Fk: "aclId"},
	}

	ModelDirectory.Register(AclMapping{}, "aclMappings", relations)
}

var AclMappingRelated *Relation

//+gen access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type AclMapping struct {
	GormModel
	RoleID int  `json:"roleId,omitempty" sql:"index"`
	Role   Role `json:"role,omitempty"`
	AclID  int  `json:"aclId,omitempty" sql:"index"`
	Acl    Acl  `json:"acl,omitempty"`
}

func (m *AclMapping) ScopeModel(roleID int) {
	m.ID = 0
	m.CreatedAt = time.Time{}
	m.UpdatedAt = m.CreatedAt

	if roleID != 0 {
		m.RoleID = roleID
	}
}
