package main

import (
	"log"
	"net/http"

	"github.com/gocroot/config"  // Pastikan untuk mengimpor paket config
	"github.com/gocroot/route"
)

func main() {
	// Inisialisasi koneksi database
	config.InitMongoDB()

	// Mengecek apakah koneksi database berhasil
	if config.ErrorMongoconn != nil {
		log.Fatalf("Gagal terhubung ke MongoDB: %v", config.ErrorMongoconn)
	}

	// Memastikan aplikasi terus berjalan dengan server HTTP
	http.HandleFunc("/", route.URL)
	log.Println("Server berjalan di port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
