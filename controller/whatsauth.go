package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/report"
	"github.com/gocroot/helper/whatsauth"
	"github.com/gocroot/model"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
)

func GetHome(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	resp.Response = at.GetIPaddress()
	at.WriteJSON(respw, http.StatusOK, resp)
}

func PostInboxNomor(respw http.ResponseWriter, req *http.Request) {
    var resp itmodel.Response
    var msg itmodel.IteungMessage
    httpstatus := http.StatusUnauthorized
    resp.Response = "Wrong Secret or Phone Number Not Registered"

    // Ambil nomor WA dari request
    waphonenumber := at.GetParam(req)

    // Validasi apakah nomor HP terdaftar di database
    var user struct {
        NoHP string `bson:"no_hp"`
    }
	err := config.Mongoconn.Client().Database("akuntan").Collection("user").FindOne(req.Context(), bson.M{"no_hp": waphonenumber}).Decode(&user)
    if err != nil {
        // Jika nomor tidak ditemukan, balikan respon
        resp.Response = "Phone number not registered"
        at.WriteJSON(respw, httpstatus, resp)
        return
    }

    // Dapatkan profil aplikasi
    prof, err := whatsauth.GetAppProfile(waphonenumber, config.Mongoconn)
    if err != nil {
        resp.Response = err.Error()
        httpstatus = http.StatusServiceUnavailable
        at.WriteJSON(respw, httpstatus, resp)
        return
    }

    // Cek Secret dari header
    if at.GetSecretFromHeader(req) == prof.Secret {
        // Decode pesan dari body
        err := json.NewDecoder(req.Body).Decode(&msg)
        if err != nil {
            resp.Response = err.Error()
        } else {
            // Proses WebHook
            resp, err = whatsauth.WebHook(prof, msg, config.Mongoconn)
            if err != nil {
                resp.Response = err.Error()
            }
        }
    }
    at.WriteJSON(respw, httpstatus, resp)
}

// jalan setiap jam 3 pagi
func GetNewToken(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	httpstatus := http.StatusServiceUnavailable

	var wg sync.WaitGroup
	wg.Add(3) // Menambahkan jumlah goroutine yang akan dijalankan

	// Mutex untuk mengamankan akses ke variabel resp dan httpstatus
	var mu sync.Mutex
	// Variabel untuk menyimpan kesalahan terakhir
	var lastErr error

	// 1. Refresh token
	go func() {
		defer wg.Done() // Memanggil wg.Done() setelah fungsi selesai
		profs, err := atdb.GetAllDoc[[]model.Profile](config.Mongoconn, "profile", bson.M{})
		if err != nil {
			mu.Lock()
			lastErr = err
			resp.Response = err.Error()
			mu.Unlock()
			return
		}
		for _, prof := range profs {
			dt := &itmodel.WebHook{
				URL:    prof.URL,
				Secret: prof.Secret,
			}
			res, err := whatsauth.RefreshToken(dt, prof.Phonenumber, config.WAAPIGetToken, config.Mongoconn)
			if err != nil {
				mu.Lock()
				lastErr = err
				resp.Response = err.Error()
				httpstatus = http.StatusInternalServerError
				mu.Unlock()
				continue // Lanjutkan ke iterasi berikutnya
			}
			mu.Lock()
			resp.Response = at.Jsonstr(res.ModifiedCount)
			httpstatus = http.StatusOK
			mu.Unlock()
		}
	}()

	// 2. Menjalankan fungsi RekapMeetingKemarin dalam goroutine
	go func() {
		defer wg.Done() // Memanggil wg.Done() setelah fungsi selesai
		if err := report.RekapMeetingKemarin(config.Mongoconn); err != nil {
			mu.Lock()
			lastErr = err
			resp.Response = err.Error()
			httpstatus = http.StatusInternalServerError
			mu.Unlock()
		}
	}()

	// 3. Menjalankan fungsi RekapPagiHari dalam goroutine
	go func() {
		defer wg.Done() // Memanggil wg.Done() setelah fungsi selesai
		if err := report.RekapPagiHari(config.Mongoconn); err != nil {
			mu.Lock()
			lastErr = err
			resp.Response = err.Error()
			httpstatus = http.StatusInternalServerError
			mu.Unlock()
		}
	}()

	wg.Wait() // Menunggu sampai semua goroutine selesai

	// Menggunakan status yang benar dari kesalahan terakhir jika ada
	if lastErr != nil {
		at.WriteJSON(respw, httpstatus, resp)
	} else {
		at.WriteJSON(respw, http.StatusOK, resp)
	}
}

func NotFound(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	resp.Response = "Not Found"
	at.WriteJSON(respw, http.StatusNotFound, resp)
}
