package domain

type modelDirectory struct {
	Models    []interface{}
	Relations map[string][]DBRelation
}

type DBRelation struct {
	Resource, Fk, Related string
}

func newModelDirectory() *modelDirectory {
	modelDir := &modelDirectory{}
	modelDir.Relations = make(map[string][]DBRelation)

	return modelDir
}

func (md *modelDirectory) Register(model interface{}, resource string, relations []DBRelation) {
	md.Models = append(md.Models, model)

	for i := range relations {
		relations[i].Resource = resource
	}

	md.Relations[resource] = relations
}

// WARNING: MIGHT NOT WORK WHEN MULTIPLE PATH ARE AVAILABLE.
func (md *modelDirectory) FindPathToOwner(resource string) []DBRelation {
	relationPath := []DBRelation{}

	if resource == "accounts" {
		return relationPath
	}

	relationPath = md.findPathToOwner("", resource, relationPath)

	return relationPath
}

func (md *modelDirectory) findPathToOwner(lastResource, resource string, relationPath []DBRelation) []DBRelation {
	relations := md.Relations[resource]

	for _, relation := range relations {
		if relation.Related == lastResource && relation.Fk != "" {
			relationPath = append(relationPath, relation)
		}
	}

	for _, relation := range relations {
		if relation.Fk == "accountId" {
			return relationPath
		} else if !containsResource(relationPath, relation.Related) {
			relationPathTmp := relationPath
			if relation.Fk != "" {
				relationPathTmp = append(relationPath, relation)
			}
			relationPath := md.findPathToOwner(resource, relation.Related, relationPathTmp)
			if relationPath != nil {
				return relationPath
			}
		}
	}

	return nil
}

func containsResource(relations []DBRelation, resource string) bool {
	for _, relation := range relations {
		if relation.Related == resource {
			return true
		}
	}
	return false
}

var ModelDirectory *modelDirectory = newModelDirectory()
