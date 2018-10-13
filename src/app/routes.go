package app

import (
	"github.com/fragmenta/mux"
	"github.com/fragmenta/mux/middleware/gzip"
	"github.com/fragmenta/mux/middleware/logrequest"
)

// SetupRoutes creates a new router and adds the routes for this app to it.
func SetupRoutes() *mux.Mux {

	router := mux.New()
	mux.SetDefault(router)

	// Add middleware
	router.AddMiddleware(logrequest.Middleware)
	router.AddMiddleware(gzip.Middleware)

	// Add the home page route
	router.Get("/", homeHandler)

	return router
}
