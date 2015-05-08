package domain

import "time"

type GormModel struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type modelDirectory []interface{}

func (md modelDirectory) Register(model interface{}) {
	md = append(md, model)
}

var Models modelDirectory = modelDirectory{}
