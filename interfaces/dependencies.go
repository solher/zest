package interfaces

import (
	"net/http"

	"github.com/Solher/zest/apierrors"
	"github.com/Solher/zest/domain"
	"github.com/Solher/zest/usecases"
	"github.com/jinzhu/gorm"
)

type (
	AbstractGormStore interface {
		Connect(adapter, url string) error
		Close() error
		GetDB() *gorm.DB
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
