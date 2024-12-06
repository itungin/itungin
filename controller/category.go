package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCategory(respw http.ResponseWriter, req *http.Request) {
	var category model.Category
	if err := json.NewDecoder(req.Body).Decode(&category); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	newCategory := model.Category{
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := atdb.InsertOneDoc(config.Mongoconn, "kategori", newCategory)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Kategori berhasil ditambahkan",
		"data":    newCategory,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func GetAllCategory(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Category](config.Mongoconn, "kategori", primitive.M{})
	if err != nil || len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data kategori tidak ditemukan atau kosong"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, data)
}

func GetCategoryByID(respw http.ResponseWriter, req *http.Request) {
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var category model.Category
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Category](config.Mongoconn, "kategori", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, category)
}

func UpdateCategory(respw http.ResponseWriter, req *http.Request) {
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	updateData := bson.M{}
	if requestBody.Name != "" {
		updateData["name"] = requestBody.Name
	}
	if requestBody.Description != "" {
		updateData["description"] = requestBody.Description
	}
	updateData["updatedAt"] = time.Now()

	update := bson.M{"$set": updateData}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "kategori", bson.M{"_id": objectID}, update)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate category"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, updateData)
}

func DeleteCategory(respw http.ResponseWriter, req *http.Request) {
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "kategori", bson.M{"_id": objectID})
	if err != nil || deleteResult.DeletedCount == 0 {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus category atau category tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Category berhasil dihapus",
	})
}
