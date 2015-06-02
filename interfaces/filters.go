package interfaces

type GormFilter struct {
	Fields  []string
	Limit   int
	Order   string
	Offset  int
	Where   string
	Include []GormInclude
}

type GormInclude struct {
	Relation, Where string
}
