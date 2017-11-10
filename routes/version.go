package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danackerson/battlefleet/app"
)

// VersionHandler now commented
func VersionHandler(w http.ResponseWriter, req *http.Request) {
	buildURL := "https://circleci.com/gh/danackerson/battlefleet/" + app.Version
	v := map[string]string{"version": buildURL, "build": app.Version}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}
