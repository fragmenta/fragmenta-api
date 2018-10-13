package app

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fragmenta/server/config"
)

// TestRouter tests our routes are functioning correctly.
func TestRouter(t *testing.T) {

	// chdir into the root to load config/assets to test this code
	err := os.Chdir("../../")
	if err != nil {
		t.Errorf("Chdir error: %s", err)
	}

	c := config.New()
	c.Load("secrets/fragmenta.json")
	c.Mode = config.ModeTest
	config.Current = c

	// First, set up the logger
	err = SetupLog()
	if err != nil {
		t.Fatalf("app: failed to set up log %s", err)
	}

	// Setup our view templates
	SetupView()

	// Setup our database
	SetupDatabase()

	// Setup our router and handlers
	router := SetupRoutes()

	// Test serving the route / which should always exist
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// Test code on response
	if w.Code != http.StatusOK {
		t.Fatalf("app: error code on / expected:%d got:%d", http.StatusOK, w.Code)
	}
}
