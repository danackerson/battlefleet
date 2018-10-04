package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danackerson/battlefleet/app"
)

// Version of application build
type Version struct {
	URL string
	Tag string
}

// GetVersionInfo does what it says on the box
func GetVersionInfo() *Version {
	urlVersion := app.Version

	if strings.HasPrefix(urlVersion, "vc") {
		urlVersion = strings.TrimLeft(urlVersion, "vc")
	}
	versionInfo := &Version{
		// strip "vc" off version
		URL: "https://circleci.com/gh/danackerson/battlefleet/" + urlVersion,
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
