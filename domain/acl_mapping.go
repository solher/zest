package domain

import "time"

func init() {
	ModelDirectory.Register(AclMapping{})
	AclMappingRelated = &Related{ModelName: "role"}
	AclMappingRelated.Add("roleMapping")
}

var AclMappingRelated *Related

//+gen access controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
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
