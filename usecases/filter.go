package usecases

type Filter struct {
	Fields  []string               `json:"fields"`
	Limit   int                    `json:"limit"`
	Order   string                 `json:"order"`
	Offset  int                    `json:"offset"`
	Where   map[string]interface{} `json:"where"`
	Include []interface{}          `json:"include"`
}
