package domain

func init() {
	relations := []Relation{
		{Related: "roleMappings"},
		{Related: "aclMappings"},
	}

	ModelDirectory.Register(Role{}, "roles", relations)
}

//+gen hooks access controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type Role struct {
	GormModel
	Name         string        `json:"name,omitempty"`
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
	AclMappings  []AclMapping  `json:"aclMappings,omitempty"`
}

func (m *Role) SetRelatedID(idKey string, id int) {
}

func (m *Role) ScopeModel() error {
	return nil
}

func (m *Role) BeforeRender() {
	roleMappings := m.RoleMappings
	aclMappings := m.AclMappings

	for i := range roleMappings {
		(&roleMappings[i]).BeforeRender()
	}

	for i := range aclMappings {
		(&aclMappings[i]).BeforeRender()
	}
}
