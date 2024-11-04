package main

import (
	"net/http"
	"log"

	"github.com/gocroot/route"
)

func main() {
	// Memastikan aplikasi terus berjalan dengan server HTTP
	http.HandleFunc("/", route.URL)
	log.Println("Server berjalan di port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
