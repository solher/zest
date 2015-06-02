package domain

func init() {
	relations := []Relation{
		{Related: "roles", Fk: "roleId"},
		{Related: "acls", Fk: "aclId"},
	}

	ModelDirectory.Register(AclMapping{}, "aclMappings", relations)
}

//+gen hooks access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
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

func (m *AclMapping) ScopeModel() error {
	return nil
}

func (m *AclMapping) BeforeRender() error {
	m.Role.BeforeRender()
	m.Acl.BeforeRender()
	return nil
}
