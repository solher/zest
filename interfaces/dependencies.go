package interfaces

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/solher/zest/apierrors"
	"github.com/solher/zest/domain"
	"github.com/solher/zest/usecases"
)

type (
	AbstractGormStore interface {
		Connect(adapter, url string) error
		Close() error
		GetDB() *gorm.DB
		BuildQuery(filter *usecases.Filter, ownerRelations []domain.DBRelation) (*gorm.DB, error)
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
