package emails

import "github.com/Solher/auth-scaffold/domain"

func init() {
	domain.ModelDirectory.Register(Email{})
}

//go:generate gen -f
//+gen routes controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" controller_test interactor:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" interactor_test repository:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" repository_test
type Email struct {
	domain.GormModel
	ID     int    `json:"id,omitempty" gorm:"primary_key"`
	UserID int    `json:"userId,omitempty" sql:"index"`
	Email  string `json:"email,omitempty"`
}
