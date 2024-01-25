package session

import (
	"sync"
	"time"
)

// TODO: Decide if I want to use a session
type Manager struct {
	name        string
	lock        sync.Mutex
	sessions    map[string]interface{}
	key         sessionKey
	maxLifetime time.Duration
}

type sessionKey struct{}

func NewManager(name string, lifetime time.Duration) *Manager {
	return &Manager{
		name:        name,
		sessions:    make(map[string]interface{}),
		key:         sessionKey{},
		maxLifetime: lifetime,
	}
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifetime int64)
}

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionId() string
}
