package gocroot

import (
	"github.com/gocroot/route"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {

	// Mendaftarkan fungsi HTTP untuk Google Cloud Functions
	functions.HTTP("WebHook", route.URL)
}


