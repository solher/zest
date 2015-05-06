package emails

import "github.com/Solher/auth-boilerplate-4/models"

// +gen repository:"Create"
type Email struct {
	models.GormModel
	ID     int    `json:"id,omitempty"`
	UserID int    `json:"userId,omitempty" sql:"index"`
	Email  string `json:"email,omitempty"`
}
