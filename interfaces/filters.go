package interfaces

type Filter struct {
	Fields  []string               `json:"fields"`
	Limit   int                    `json:"limit"`
	Order   string                 `json:"order"`
	Offset  int                    `json:"offset"`
	Where   map[string]interface{} `json:"where"`
	Include []interface{}          `json:"include"`
}

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

//[
// 	{
// 		"relation": "toto",
//		"filter": Filter,
//		"include": Include
// 	},
// 	"tutu",
//]
