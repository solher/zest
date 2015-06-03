package domain

func init() {
	relations := []Relation{
		{Related: "aclMappings"},
	}

	ModelDirectory.Register(Acl{}, "acls", relations)
}

//+gen routes controller:"Create,Find,FindByID,Upsert,UpdateByID,DeleteAll,DeleteByID,Related,RelatedOne" repository:"Create,CreateOne,Find,FindByID,Update,UpdateByID,DeleteAll,DeleteByID,Raw" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,UpdateByID,DeleteAll,DeleteByID"
type Acl struct {
	GormModel
	Ressource   string       `json:"ressource,omitempty"`
	Method      string       `json:"method,omitempty"`
	AclMappings []AclMapping `json:"aclMappings,omitempty"`
}

func (m *Acl) SetRelatedID(idKey string, id int) {
}

func (m *Acl) ScopeModel() error {
	return nil
}

func (m *Acl) BeforeRender() {
	aclMappings := m.AclMappings

	for i := range aclMappings {
		(&aclMappings[i]).BeforeRender()
	}
}
