package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"github.com/gocroot/helper/at"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Fungsi untuk menambahkan kategori baru
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory model.Category
	if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
		var response model.Response
		response.Status = "Error: Bad Request"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Set waktu pembuatan kategori
	newCategory.CreatedAt = time.Now()

	// Insert kategori ke MongoDB
	_, err := config.CategoryCollection.InsertOne(context.Background(), newCategory)
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
		"message": "Kategori berhasil dibuat",
		"data":    newCategory,
	}
	at.WriteJSON(w, http.StatusCreated, response)
}

// Fungsi untuk mendapatkan daftar kategori
func GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []model.Category

	// Ambil data dari MongoDB
	cursor, err := config.CategoryCollection.Find(context.Background(), bson.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengambil data kategori"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}
	defer cursor.Close(context.Background())

	// Decode hasil pencarian kategori
	for cursor.Next(context.Background()) {
		var category model.Category
		if err := cursor.Decode(&category); err != nil {
			var response model.Response
			response.Status = "Error: Gagal mendekode kategori"
			at.WriteJSON(w, http.StatusInternalServerError, response)
			return
		}
		categories = append(categories, category)
	}

	// Kirim data kategori sebagai respon
	at.WriteJSON(w, http.StatusOK, categories)
}

// Fungsi untuk mendapatkan detail kategori berdasarkan ID
func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category model.Category
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = config.CategoryCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// Fungsi untuk mengupdate kategori berdasarkan ID
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var updatedCategory model.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":        updatedCategory.Name,
			"description": updatedCategory.Description,
			"updatedAt":   time.Now(),
		},
	}

	_, err = config.CategoryCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category updated successfully"})
}

// Fungsi untuk menghapus kategori berdasarkan ID
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.CategoryCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted successfully"})
}
