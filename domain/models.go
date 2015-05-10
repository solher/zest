package domain

import "time"

type GormModel struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type modelDirectory struct {
	Models []interface{}
}

func (md *modelDirectory) Register(model interface{}) {
	md.Models = append(md.Models, model)
}

var ModelDirectory *modelDirectory = &modelDirectory{}
