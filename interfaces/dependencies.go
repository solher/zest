package interfaces

import (
	"net/http"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/domain"
	"github.com/jinzhu/gorm"
)

type (
	AbstractGormStore interface {
		Connect(adapter, url string) error
		Close() error
		GetDB() *gorm.DB
		MigrateTables(tables []interface{}) error
		ReinitTables(tables []interface{}) error
		BuildQuery(filter *Filter, ownerRelations []domain.Relation) (*gorm.DB, error)
	}

	AbstractAccountInter interface {
		GetGrantedRoles(accountID int, ressource, method string) ([]string, error)
	}

	AbstractRender interface {
		JSONError(w http.ResponseWriter, status int, apiError *apierrors.APIError, err error)
		JSON(w http.ResponseWriter, status int, object interface{})
	}
)
