package gocroot

import (
	"log"
	"github.com/gocroot/config"  // Import paket config untuk koneksi database
	"github.com/gocroot/route"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	// Mengecek apakah koneksi database berhasil
	if config.ErrorMongoconn != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", config.ErrorMongoconn)
	}

	// Mendaftarkan fungsi HTTP untuk Google Cloud Functions
	functions.HTTP("WebHook", route.URL)
}


