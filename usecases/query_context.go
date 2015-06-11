package usecases

import "github.com/solher/zest/domain"

type QueryContext struct {
	Filter         *Filter
	OwnerRelations []domain.DBRelation
}

func NewQueryContext(filter *Filter, ownerRelations []domain.DBRelation) *QueryContext {
	return &QueryContext{Filter: filter, OwnerRelations: ownerRelations}
}

type Filter struct {
	Fields  []string               `json:"fields"`
	Limit   int                    `json:"limit"`
	Order   string                 `json:"order"`
	Offset  int                    `json:"offset"`
	Where   map[string]interface{} `json:"where"`
	Include []interface{}          `json:"include"`
}
