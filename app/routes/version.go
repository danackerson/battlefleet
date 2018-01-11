package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danackerson/battlefleet/app"
)

type version struct {
	URL string
	Tag string
}

// GetVersionURL does what it says on the box
func GetVersionInfo() *version {
	versionInfo := &version{
		URL: "https://circleci.com/gh/danackerson/battlefleet/" + app.Version,
		Tag: app.Version,
	}
	return versionInfo
}

// VersionHandler now commented
func VersionHandler(w http.ResponseWriter, req *http.Request) {
	versionInfo := GetVersionInfo()
	v := map[string]string{"version": versionInfo.URL, "build": versionInfo.Tag}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}
