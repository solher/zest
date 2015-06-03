package usecases

import "github.com/Solher/zest/domain"

type AbstractCacheStore interface {
	Add(key interface{}, value interface{}) error
	Remove(key interface{}) error
	Get(key interface{}) (interface{}, error)
	Purge() error
	MaxSize() int
}

type AbstractPermissionCacheInter interface {
	GetPermissionRoles(accountID int, ressource, method string) ([]string, error)
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
