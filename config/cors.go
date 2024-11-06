package config

import (
	"net/http"
)

// Daftar origins yang diizinkan
var Origins = []string{
	"https://www.bukupedia.co.id",
	"https://naskah.bukupedia.co.id",
	"https://bukupedia.co.id",
	"http://127.0.0.1:5500", // Origin lokal untuk pengujian
}

// Fungsi untuk memeriksa apakah origin diizinkan
func isAllowedOrigin(origin string) bool {
	for _, o := range Origins {
		if o == origin {
			return true
		}
	}
	return false
}

// Fungsi untuk mengatur header CORS
func SetAccessControlHeaders(w http.ResponseWriter, r *http.Request) bool {
    origin := r.Header.Get("Origin")

    // Pastikan origin yang diminta diizinkan
    if isAllowedOrigin(origin) || origin == "" {
        // Tangani preflight request (OPTIONS)
        if r.Method == http.MethodOptions {
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Login")
            w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS")
            w.Header().Set("Access-Control-Allow-Origin", origin)  // Set header origin yang diizinkan
            w.Header().Set("Access-Control-Max-Age", "3600")       // Cache preflight request selama 1 jam
            w.WriteHeader(http.StatusNoContent)                     // Tidak ada konten, status OK untuk preflight
            return true
        }

        // Tangani permintaan utama (GET, POST, PUT, DELETE, dll.)
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Origin", origin) // Set origin yang diizinkan
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, OPTIONS") // Pastikan metode lainnya diizinkan
        return false
    }

    // Jika origin tidak diizinkan, tidak ada header CORS yang diatur
    return false
}
