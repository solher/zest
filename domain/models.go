package domain

import "time"

//go:generate ./_gen.sh

type GormModel struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	ID        int       `json:"id,omitempty" gorm:"primary_key"`
}

type Related struct {
	ModelName string
	Next      *Related
}

func (r *Related) Add(modelName string) *Related {
	r.Next = &Related{ModelName: modelName}

	return r.Next
}

type modelDirectory struct {
	Models []interface{}
}

func (md *modelDirectory) Register(model interface{}) {
	md.Models = append(md.Models, model)
}

var ModelDirectory *modelDirectory = &modelDirectory{}
