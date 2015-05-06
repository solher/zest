package interfaces

type Filter struct {
	Fields []string               `json:"fields"`
	Limit  int                    `json:"limit"`
	Order  string                 `json:"order"`
	Offset int                    `json:"offset"`
	Where  map[string]interface{} `json:"where"`
}

type GormFilter struct {
	Fields []string
	Limit  int
	Order  string
	Offset int
	Where  string
}
