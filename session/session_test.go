package session

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type Session struct {
	name string
}

func FakeLogin(sess *Manager) {
	req, _ := http.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	sess.PutSession(rr, req, &Session{name: "testing"})
}

func TestSessionManager(t *testing.T) {
	sess := NewManager("test", time.Minute*5)

	var wg sync.WaitGroup
	numGoroutines := 20
	for i := 0; i <= numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for created := 0; created <= 50000; created++ {
				FakeLogin(sess)
			}
		}()
	}
	wg.Wait()
	fmt.Println("10 million sessions simulated")
}

func TestSessionManagerPut(t *testing.T) {
	sess := NewManager("test", time.Minute*5)

	req, _ := http.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	sess.PutSession(rr, req, &Session{name: "testing"})

	assert.Equal(t, Session{name: "testing"}, sess.GetSessionFromRequest(req))
}
