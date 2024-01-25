package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

// TODO: Decide if I want to use a session
type Manager struct {
	name     string
	lock     sync.RWMutex
	sessions map[string]interface{}
	key      sessionKey
	lifetime time.Duration
}

type sessionKey struct{}

func NewManager(name string, lifetime time.Duration) *Manager {
	return &Manager{
		name:     name,
		sessions: make(map[string]interface{}),
		key:      sessionKey{},
		lifetime: lifetime,
	}
}

func randBase64String(entropyBytes int) string {
	b := make([]byte, entropyBytes)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (m *Manager) PutSession(w http.ResponseWriter, r *http.Request, sess interface{}) {
	cookieValue := randBase64String(33)

	time.AfterFunc(m.lifetime, func() {
		m.lock.Lock()
		delete(m.sessions, cookieValue)
		m.lock.Unlock()
	})

	m.lock.Lock()
	m.sessions[cookieValue] = sess
	m.lock.Unlock()

	cookie := &http.Cookie{
		Name:     m.name,
		Value:    cookieValue,
		Expires:  time.Now().Add(m.lifetime),
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func (m *Manager) DeleteSession(r *http.Request) {
	cookie, err := r.Cookie(m.name)
	if err != nil {
		return
	}
	m.lock.Lock()
	delete(m.sessions, cookie.Value)
	m.lock.Unlock()
}

func (m *Manager) GetSessionFromRequest(r *http.Request) interface{} {
	cookie, err := r.Cookie(m.name)
	if err != nil {
		return nil
	}
	m.lock.RLock()
	sess := m.sessions[cookie.Value]
	m.lock.RUnlock()

	return sess
}

func (m *Manager) LoadSessionIntoContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := m.GetSessionFromRequest(r)
		if sess == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), m.key, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Manager) GetSessionFromContext(r *http.Request) interface{} {
	return r.Context().Value(m.key)
}
