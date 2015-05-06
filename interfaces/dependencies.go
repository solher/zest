package interfaces

import (
	"net/http"

	"github.com/Solher/auth-scaffold/apierrors"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

type (
	GormStore interface {
		Connect(adapter, url string) error
		Close() error
		GetDB() *gorm.DB
		MigrateTables(tables []interface{}) error
		ReinitTables(tables []interface{}) error
		BuildQuery(filter *Filter) (*gorm.DB, error)
	}

	Render interface {
		JSONError(w http.ResponseWriter, status int, apiError *apierrors.APIError, err error)
		JSON(w http.ResponseWriter, status int, object interface{})
	}

	Router interface {
		Handle(method string, path string, handle Handle)
	}

	Handle func(http.ResponseWriter, *http.Request, httprouter.Params)
)
