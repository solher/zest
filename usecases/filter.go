package usecases

type AbstractFilter interface {
	Fields() []string
	Limit() int
	Order() string
	Offset() int
	Where() map[string]interface{}
	Include() []interface{}
	SetFields(fields []string)
	SetLimit(limit int)
	SetOrder(order string)
	SetOffset(offset int)
	SetWhere(where map[string]interface{})
	SetInclude(include []interface{})
}

type Filter struct {
	fields  []string
	limit   int
	order   string
	offset  int
	where   map[string]interface{}
	include []interface{}
}

func (f *Filter) Fields() []string {
	return f.fields
}

func (f *Filter) Limit() int {
	return f.limit
}

func (f *Filter) Order() string {
	return f.order
}

func (f *Filter) Offset() int {
	return f.offset
}

func (f *Filter) Where() map[string]interface{} {
	return f.where
}

func (f *Filter) Include() []interface{} {
	return f.include
}

func (f *Filter) SetFields(fields []string) {
	f.fields = fields
}

func (f *Filter) SetLimit(limit int) {
	f.limit = limit
}

func (f *Filter) SetOrder(order string) {
	f.order = order
}

func (f *Filter) SetOffset(offset int) {
	f.offset = offset
}

func (f *Filter) SetWhere(where map[string]interface{}) {
	f.where = where
}

func (f *Filter) SetInclude(include []interface{}) {
	f.include = include
}
