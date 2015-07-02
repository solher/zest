package domain

type modelDirectory struct {
	Models    []interface{}
	Relations map[string][]DBRelation
}

type DBRelation struct {
	Ressource, Fk, Related string
}

func newModelDirectory() *modelDirectory {
	modelDir := &modelDirectory{}
	modelDir.Relations = make(map[string][]DBRelation)

	return modelDir
}

func (md *modelDirectory) Register(model interface{}, ressource string, relations []DBRelation) {
	md.Models = append(md.Models, model)

	for i := range relations {
		relations[i].Ressource = ressource
	}

	md.Relations[ressource] = relations
}

// WARNING: MIGHT NOT WORK WHEN MULTIPLE PATH ARE AVAILABLE.
func (md *modelDirectory) FindPathToOwner(ressource string) []DBRelation {
	relationPath := []DBRelation{}

	if ressource == "accounts"{
		return relationPath
	}

	relationPath = md.findPathToOwner("", ressource, relationPath)

	return relationPath
}

func (md *modelDirectory) findPathToOwner(lastRessource, ressource string, relationPath []DBRelation) []DBRelation {
	relations := md.Relations[ressource]

	for _, relation := range relations {
		if relation.Related == lastRessource && relation.Fk != "" {
			relationPath = append(relationPath, relation)
		}
	}

	for _, relation := range relations {
		if relation.Fk == "accountId" {
			return relationPath
		} else if !containsRessource(relationPath, relation.Related) {
			relationPathTmp := relationPath
			if relation.Fk != "" {
				relationPathTmp = append(relationPath, relation)
			}
			relationPath := md.findPathToOwner(ressource, relation.Related, relationPathTmp)
			if relationPath != nil {
				return relationPath
			}
		}
	}

	return nil
}

func containsRessource(relations []DBRelation, ressource string) bool {
	for _, relation := range relations {
		if relation.Related == ressource {
			return true
		}
	}
	return false
}

var ModelDirectory *modelDirectory = newModelDirectory()
