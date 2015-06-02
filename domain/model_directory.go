package domain

type modelDirectory struct {
	Models    []interface{}
	Relations map[string][]Relation
}

type Relation struct {
	Ressource, Fk, Related string
}

func newModelDirectory() *modelDirectory {
	modelDir := &modelDirectory{}
	modelDir.Relations = make(map[string][]Relation)

	return modelDir
}

func (md *modelDirectory) Register(model interface{}, ressource string, relations []Relation) {
	md.Models = append(md.Models, model)

	for i := range relations {
		relations[i].Ressource = ressource
	}

	md.Relations[ressource] = relations
}

// WARNING: MIGHT NOT WORK WHEN MULTIPLE PATH ARE AVAILABLE.
func (md *modelDirectory) FindPathToOwner(ressource string) []Relation {
	relationPath := []Relation{}

	relationPath = md.findPathToOwner("", ressource, relationPath)

	return relationPath
}

func (md *modelDirectory) findPathToOwner(lastRessource, ressource string, relationPath []Relation) []Relation {
	relations := md.Relations[ressource]

	for _, relation := range relations {
		if relation.Related == lastRessource && relation.Fk != "" {
			relationPath = append(relationPath, relation)
		}
	}

	for _, relation := range relations {
		if relation.Fk == "accountId" {
			return append(relationPath, relation)
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

func containsRessource(relations []Relation, ressource string) bool {
	for _, relation := range relations {
		if relation.Related == ressource {
			return true
		}
	}
	return false
}

var ModelDirectory *modelDirectory = newModelDirectory()
