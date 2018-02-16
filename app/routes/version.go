package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danackerson/battlefleet/app"
)

// Version of application build
type Version struct {
	URL string
	Tag string
}

// GetVersionInfo does what it says on the box
func GetVersionInfo() *Version {
	versionInfo := &Version{
		URL: "https://circleci.com/gh/danackerson/battlefleet/" + app.Version,
		Tag: app.Version,
	}
	return versionInfo
}

// VersionHandler now commented
func VersionHandler(w http.ResponseWriter, req *http.Request) {
	versionInfo := GetVersionInfo()
	//dumpRequest(req)
	v := map[string]string{"version": versionInfo.URL, "build": versionInfo.Tag}

	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write(data)
}
