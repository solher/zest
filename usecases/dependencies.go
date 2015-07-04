package usecases

import (
	"net/http"

	"github.com/solher/zest/apierrors"
	"github.com/solher/zest/domain"
)

type AbstractCacheStore interface {
	Add(key interface{}, value interface{}) error
	Remove(key interface{}) error
	Get(key interface{}) (interface{}, error)
	Purge() error
	MaxSize() int
}

type AbstractPermissionCacheInter interface {
	GetPermissionRoles(accountID int, resource, method string) ([]string, error)
	Refresh() error
	RefreshRole(accountID int) error
	RefreshAcl(aclID int) error
}

type AbstractSessionCacheInter interface {
	Add(authToken string, session domain.Session) error
	Remove(authToken string) error
	Get(authToken string) (domain.Session, error)
	Refresh() error
	RefreshSession(sessionID int) error
}

type AbstractAccountInter interface {
	GetGrantedRoles(accountID int, resource, method string) ([]string, error)
}

type AbstractRender interface {
	JSONError(w http.ResponseWriter, status int, apiError *apierrors.APIError, err error)
	JSON(w http.ResponseWriter, status int, object interface{})
}
