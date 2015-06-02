package interfaces

import (
	"net/http"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/Solher/auth-scaffold/domain"
	"github.com/Solher/auth-scaffold/usecases"
	"github.com/jinzhu/gorm"
)

type (
	AbstractGormStore interface {
		Connect(adapter, url string) error
		Close() error
		GetDB() *gorm.DB
		MigrateTables(tables []interface{}) error
		ReinitTables(tables []interface{}) error
		BuildQuery(filter *usecases.Filter, ownerRelations []domain.Relation) (*gorm.DB, error)
	}

	AbstractAccountInter interface {
		GetGrantedRoles(accountID int, ressource, method string) ([]string, error)
	}

	AbstractRender interface {
		JSONError(w http.ResponseWriter, status int, apiError *apierrors.APIError, err error)
		JSON(w http.ResponseWriter, status int, object interface{})
	}

	AbstractCacheStore interface {
		Add(key interface{}, value interface{}) error
		Remove(key interface{}) error
		Get(key interface{}) (interface{}, error)
		Purge() error
	}
)
