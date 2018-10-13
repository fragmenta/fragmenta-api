package app

import (
	"net/http"
	"os"

	"github.com/fragmenta/view"
)

var hostName string

// HandleShowHome servers a json response for a health check, including the hostname
func homeHandler(w http.ResponseWriter, r *http.Request) error {

	if hostName == "" {
		hostName, _ = os.Hostname()
	}

	view := view.NewRenderer(w, r)
	view.AddKey("host", hostName)
	view.AddKey("message", "API")
	view.Template("app/views/health.json.got")
	return view.Render()
}
