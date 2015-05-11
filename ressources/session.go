package ressources

import (
	"time"

	"github.com/Solher/auth-scaffold/domain"
)

func init() {
	domain.ModelDirectory.Register(Session{})
}

//go:generate gen -f
//+gen routes controller:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" interactor:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID" repository:"Create,Find,FindByID,Upsert,DeleteAll,DeleteByID"
type Session struct {
	domain.GormModel
	ID        int       `json:"id,omitempty" gorm:"primary_key"`
	AccountID int       `json:"accountId,omitempty" sql:"index"`
	Account   Account   `json:"account,omitempty"`
	AuthToken string    `json:"authToken,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Agent     string    `json:"agent,omitempty"`
	ValidTo   time.Time `json:"validTo,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
}
