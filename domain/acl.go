package domain

func init() {
	relations := []DBRelation{
		{Related: "aclMappings"},
	}

	ModelDirectory.Register(Acl{}, "acls", relations)
}

type Acl struct {
	GormModel
	Resource    string       `json:"resource"`
	Method      string       `json:"method"`
	AclMappings []AclMapping `json:"aclMappings,omitempty"`
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
