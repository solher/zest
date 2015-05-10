package sessions

import (
	"time"

	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/ressources/users"
)

func init() {
	domain.ModelDirectory.Register(Session{})
}

//go:generate gen -f
//+gen routes controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" controller_test interactor:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" interactor_test repository:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" repository_test
type Session struct {
	domain.GormModel
	ID        int        `json:"id,omitempty" gorm:"primary_key"`
	UserID    int        `json:"userId,omitempty" sql:"index"`
	User      users.User `json:"user,omitempty"`
	AuthToken string     `json:"authToken,omitempty"`
	IP        string     `json:"ip,omitempty"`
	Agent     string     `json:"agent,omitempty"`
	ValidTo   time.Time  `json:"validTo,omitempty"`
	DeletedAt time.Time  `json:"deletedAt,omitempty"`
}
