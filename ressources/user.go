package ressources

import "github.com/Solher/auth-scaffold/domain"

func init() {
	domain.ModelDirectory.Register(User{})
}

//go:generate gen -f
//+gen routes controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" interactor:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID" repository:"Create,CreateOne,Find,FindByID,Upsert,UpsertOne,DeleteAll,DeleteByID"
type User struct {
	domain.GormModel
	ID        int    `json:"id,omitempty" gorm:"primary_key"`
	AccountID int    `json:"accountId,omitempty" sql:"index"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Password  string `json:"-"`
	Email     string `json:"email,omitempty" sql:"unique"`
}
