package domain

func init() {
	relations := []DBRelation{
		{Related: "AclMappings"},
	}

	ModelDirectory.Register(Acl{}, "acls", relations)
}

type Acl struct {
	GormModel
	Ressource   string       `json:"ressource,omitempty"`
	Method      string       `json:"method,omitempty"`
	AclMappings []AclMapping `json:"AclMappings,omitempty"`
}

func (m *Acl) SetRelatedID(idKey string, id int) {
	switch idKey {
	}
}

func (m *Acl) BeforeRender() {
	for i := range m.AclMappings {
		(&m.AclMappings[i]).BeforeRender()
	}
}
