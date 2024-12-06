package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Fungsi untuk menambahkan transaksi penjualan baru
func CreateSalesTransaction(respw http.ResponseWriter, req *http.Request) {
	var transaction model.SalesTransaction
	if err := json.NewDecoder(req.Body).Decode(&transaction); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	transaction.TransactionDate = time.Now()

	_, err := atdb.InsertOneDoc(config.Mongoconn, "transaksi_penjualan", transaction)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Transaksi berhasil ditambahkan",
		"data":    transaction,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Fungsi untuk mendapatkan semua transaksi penjualan
func GetSalesTransactions(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.SalesTransaction](config.Mongoconn, "transaksi_penjualan", bson.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data transaksi tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, data)
}

// Fungsi untuk mendapatkan transaksi penjualan berdasarkan ID
func GetSalesTransactionByID(respw http.ResponseWriter, req *http.Request) {
	transactionID := req.URL.Query().Get("id")
	if transactionID == "" {
		var respn model.Response
		respn.Status = "Error: ID Transaksi tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Transaksi tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var transaction model.SalesTransaction
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.SalesTransaction](config.Mongoconn, "transaksi_penjualan", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Transaksi tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, transaction)
}

// Fungsi untuk mengupdate transaksi penjualan berdasarkan ID
func UpdateSalesTransaction(respw http.ResponseWriter, req *http.Request) {
	transactionID := req.URL.Query().Get("id")
	if transactionID == "" {
		var respn model.Response
		respn.Status = "Error: ID Transaksi tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Transaksi tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var updatedTransaction model.SalesTransaction
	if err := json.NewDecoder(req.Body).Decode(&updatedTransaction); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	update := bson.M{"$set": updatedTransaction}
	filter := bson.M{"_id": objectID}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "transaksi_penjualan", filter, update)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate transaksi"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Transaksi berhasil diupdate",
		"data":    updatedTransaction,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Fungsi untuk menghapus transaksi penjualan berdasarkan ID
func DeleteSalesTransaction(respw http.ResponseWriter, req *http.Request) {
	transactionID := req.URL.Query().Get("id")
	if transactionID == "" {
		var respn model.Response
		respn.Status = "Error: ID Transaksi tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Transaksi tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "transaksi_penjualan", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus transaksi"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Transaksi berhasil dihapus",
		"data":    deleteResult,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
