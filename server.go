package main

import (
	"fmt"
	"os"

	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"

	"github.com/fragmenta/fragmenta-api/src/app"
)

// Main entrypoint for the server which performs bootstrap, setup
// then runs the server. Most setup is delegated to the src/app pkg.
func main() {

	// Setup our server
	server, err := SetupServer()
	if err != nil {
		fmt.Printf("server: error setting up %s\n", err)
		return
	}

	// If server requests https port, use tls instead
	if server.Port() == 443 {

		// Redirect http traffic to https
		server.StartRedirectAll(80, server.Config("root_url"))

		// Serve https directly with our certs
		err = server.StartTLS(server.Config("tls_cert"), server.Config("tls_key"))
		if err != nil {
			fmt.Printf("Error starting server %s", err)
			os.Exit(1)
		}

	} else {
		// Start the server - used to just do this
		err = server.Start()
		if err != nil {
			fmt.Printf("Error starting server %s", err)
			os.Exit(1)
		}
	}

}

// SetupServer creates a new server, and delegates setup to the app pkg.
func SetupServer() (*server.Server, error) {

	// Setup server
	s, err := server.New()
	if err != nil {
		return nil, err
	}

	// Load the appropriate config
	c := config.New()
	err = c.Load("secrets/fragmenta.json")
	if err != nil {
		return nil, err
	}
	config.Current = c

	// Check environment variable to see if we are in production mode
	if os.Getenv("FRAG_ENV") == "production" {
		config.Current.Mode = config.ModeProduction
	}

	// Call the app to perform additional setup
	app.Setup()

	return s, nil
}
