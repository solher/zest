package domain

import "time"

//go:generate ./_gen.sh

type GormModel struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	ID        int       `json:"id,omitempty" gorm:"primary_key"`
}
