package app

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/fragmenta/mux/log"
	"github.com/fragmenta/query"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/view"
)

// Setup sets up our application.
func Setup() {

	// Setup the log package, exit on failure
	err := SetupLog(config.Get("log"))
	if err != nil {
		fmt.Printf("app: failed to set up log %s", err)
		os.Exit(1)
	}

	// Now that we have a log, log server startup and completon of setup
	log.Printf("Starting server on port:%s", config.Get("port"))
	defer log.Timef(time.Now(), "Finished loading server")

	// Setup our view templates
	err = SetupView()
	if err != nil {
		log.Printf("app: failed to read views: %s", err)
		os.Exit(1)
	}
	// Setup our database
	err = SetupDatabase()
	if err != nil {
		log.Printf("app: failed to read database: %s config:%s", err, config.Get("db"))
		os.Exit(1)
	}

	// Set up our app routes
	SetupRoutes()
}

// SetupLog sets up the default loggers - one to stdout, one to a file
func SetupLog(logPath string) error {
	defer log.Timef(time.Now(), "Finished log setup")

	// Set up a stderr logger with time prefix
	logger, err := log.NewStdErr()
	if err != nil {
		return err
	}
	log.Add(logger)

	// Set up a file logger pointing at the right location for this config.
	fileLog, err := log.NewFile(logPath)
	if err != nil {
		return err
	}
	log.Add(fileLog)

	return nil
}

// SetupView sets up the view package by loadind templates.
func SetupView() error {
	defer log.Timef(time.Now(), "Finished loading templates")

	// JSON escape function
	view.Helpers["json"] = func(t string) template.HTML {
		// Escape mandatory characters
		t = strings.Replace(t, "\r", " ", -1)
		t = strings.Replace(t, "\n", " ", -1)
		t = strings.Replace(t, "\t", " ", -1)
		t = strings.Replace(t, "\\", "\\\\", -1)
		t = strings.Replace(t, "\"", "\\\"", -1)
		// Because we use html/template escape as temlate.HTML
		return template.HTML(t)
	}

	view.Production = config.Production()
	err := view.LoadTemplates()
	if err != nil {
		return err
	}

	return nil
}

// SetupDatabase sets up the query database with the settings from config.
func SetupDatabase() error {
	defer log.Timef(time.Now(), "Finished loading database %s with user %s", config.Get("db"), config.Get("db_user"))

	options := map[string]string{
		"adapter":  config.Get("db_adapter"),
		"user":     config.Get("db_user"),
		"password": config.Get("db_pass"),
		"db":       config.Get("db"),
	}

	// If host and port supplied in config, apply them - these may or may not be present
	if len(config.Get("db_host")) > 0 {
		options["host"] = config.Get("db_host")
	}

	if len(config.Get("db_port")) > 0 {
		options["port"] = config.Get("db_port")
	}

	// Ask query to open the database
	err := query.OpenDatabase(options)

	if err != nil {
		return err
	}

	query.SetMaxOpenConns(100)
	return nil
}
