package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Decode data JSON dari request body
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		var response model.Response
		response.Status = "Error: Bad Request"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Validasi data yang diperlukan
	if user.Name == "" || user.Email == "" || user.Password == "" || user.Nohp == "" {
		var response model.Response
		response.Status = "Error: Data tidak lengkap"
		response.Response = "Field Name, Email, Password, dan No HP wajib diisi"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Cek apakah email sudah terdaftar
	filter := bson.M{"email": user.Email}
	var existingUser model.User

	// Gunakan salah satu pendekatan (fungsi helper atau manual)
	// err := atdb.GetOneDoc(config.Mongoconn, "users", filter, &existingUser) // Jika fungsi helper sudah diperbaiki
	err := config.Mongoconn.Collection("users").FindOne(context.TODO(), filter).Decode(&existingUser) // Jika manual

	if err == nil { // Jika tidak ada error, artinya email ditemukan
		var response model.Response
		response.Status = "Error: Email sudah digunakan"
		response.Response = "Pengguna dengan email tersebut sudah terdaftar"
		at.WriteJSON(w, http.StatusConflict, response)
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengenkripsi password"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Inisialisasi data pengguna baru
	newUser := model.User{
		ID:       primitive.NewObjectID(),
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
		Nohp:     user.Nohp,
	}

	// Simpan pengguna ke database
	_, err = atdb.InsertOneDoc(config.Mongoconn, "users", newUser)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal Insert Database"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Kirim respon sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Pengguna berhasil didaftarkan",
		"data": map[string]interface{}{
			"id":    newUser.ID.Hex(),
			"name":  newUser.Name,
			"email": newUser.Email,
			"no_hp": newUser.Nohp,
		},
	}
	at.WriteJSON(w, http.StatusCreated, response)
}

