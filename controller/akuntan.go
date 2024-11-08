package controller

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Fungsi untuk menambahkan produk baru
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		var response model.Response
		response.Status = "Error: Bad Request"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Inisialisasi data produk baru dengan ObjectID untuk ID
	newProduct := model.Product{
		ID:          primitive.NewObjectID(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		Stock:       product.Stock,
		CreatedAt:   time.Now(),
	}

	// Insert produk ke dalam MongoDB
	_, err := atdb.InsertOneDoc(config.Mongoconn, "products", newProduct)
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
		"message": "Produk berhasil ditambahkan",
		"data":    newProduct,
	}
	at.WriteJSON(w, http.StatusCreated, response)
}



// Fungsi untuk mendapatkan daftar produk
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Ambil semua data produk dari MongoDB
	data, err := atdb.GetAllDoc[[]model.Product](config.Mongoconn, "products", primitive.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Data produk tidak ditemukan"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	if len(data) == 0 {
		var response model.Response
		response.Status = "Error: Data produk kosong"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Format hasil sebagai slice of map dengan ID, Name, Description, Price, Stock, dan CreatedAt untuk setiap produk
	var products []map[string]interface{}
	for _, product := range data {
		products = append(products, map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"category":    product.Category,
			"stock":       product.Stock,
			"createdAt":   product.CreatedAt,
		})
	}

	// Kirim data produk sebagai respon
	at.WriteJSON(w, http.StatusOK, products)
}


// Fungsi untuk mendapatkan detail produk berdasarkan ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL
	productID := r.URL.Query().Get("id")
	if productID == "" {
		var response model.Response
		response.Status = "Error: ID Produk tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Produk tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Ambil produk dari MongoDB
	var product model.Product
	filter := bson.M{"_id": objectID}
	err = config.Mongoconn.Collection("products").FindOne(context.TODO(), filter).Decode(&product)
	if err != nil {
		var response model.Response
		response.Status = "Error: Produk tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Kirim data produk sebagai respon
	response := map[string]interface{}{
		"status":  "success",
		"message": "Produk ditemukan",
		"data":    product,
	}
	at.WriteJSON(w, http.StatusOK, response)
}




// Fungsi untuk mengupdate produk berdasarkan ID
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL
	productID := r.URL.Query().Get("id")
	if productID == "" {
		var response model.Response
		response.Status = "Error: ID Produk tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Produk tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Decode data produk yang akan diupdate
	var requestBody struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Category    string  `json:"category"`
		Stock       int     `json:"stock"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Siapkan data untuk update
	updateData := bson.M{}
	if requestBody.Name != "" {
		updateData["name"] = requestBody.Name
	}
	if requestBody.Description != "" {
		updateData["description"] = requestBody.Description
	}
	if requestBody.Price != 0 {
		updateData["price"] = requestBody.Price
	}
	if requestBody.Category != "" {
		updateData["category"] = requestBody.Category
	}
	if requestBody.Stock != 0 {
		updateData["stock"] = requestBody.Stock
	}
	updateData["updatedAt"] = time.Now()

	// Update produk di MongoDB
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updateData}
	_, err = config.Mongoconn.Collection("products").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengupdate produk"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusNotModified, response)
		return
	}

	// Kirim respon sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Produk berhasil diupdate",
		"data":    updateData,
	}
	at.WriteJSON(w, http.StatusOK, response)
}



// Fungsi untuk menghapus produk berdasarkan ID
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL
	productID := r.URL.Query().Get("id")
	if productID == "" {
		var response model.Response
		response.Status = "Error: ID Produk tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Konversi productID dari string ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Produk tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Hapus data produk berdasarkan ID
	filter := bson.M{"_id": objectID}
	deleteResult, err := config.Mongoconn.Collection("products").DeleteOne(context.TODO(), filter)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal menghapus produk"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Periksa apakah ada produk yang dihapus
	if deleteResult.DeletedCount == 0 {
		var response model.Response
		response.Status = "Error: Produk tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Kirim respon sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Produk berhasil dihapus",
		"data":    deleteResult,
	}
	at.WriteJSON(w, http.StatusOK, response)
}




// Fungsi untuk mengekspor data produk ke CSV
func ExportProductsToCSV(w http.ResponseWriter, r *http.Request) {
	var products []model.Product

	// Ambil data produk dari MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.ProductCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var prod model.Product
		if err := cursor.Decode(&prod); err != nil {
			http.Error(w, "Error decoding product", http.StatusInternalServerError)
			return
		}
		products = append(products, prod)
	}

	// Tentukan header respons sebagai file CSV
	w.Header().Set("Content-Disposition", "attachment; filename=products.csv")
	w.Header().Set("Content-Type", "text/csv")

	// Buat writer CSV
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Tulis header CSV
	headers := []string{"ID", "Name", "Price", "Category", "Description", "Stock", "Created At", "Updated At"}
	if err := csvWriter.Write(headers); err != nil {
		http.Error(w, "Failed to write CSV headers", http.StatusInternalServerError)
		return
	}

	// Tulis data produk ke CSV
	for _, product := range products {
		row := []string{
			product.ID.Hex(),
			product.Name,
			formatPrice(product.Price),       // Format harga
			product.Category,
			product.Description,
			formatStock(product.Stock),       // Format stok
			product.CreatedAt.Format(time.RFC3339),
			product.UpdatedAt.Format(time.RFC3339),
		}
		if err := csvWriter.Write(row); err != nil {
			http.Error(w, "Failed to write product data to CSV", http.StatusInternalServerError)
			return
		}
	}
}

// Fungsi untuk format harga menjadi string
func formatPrice(price float64) string {
	return fmt.Sprintf("%.2f", price)
}

// Fungsi untuk format stok menjadi string
func formatStock(stock int) string {
	return fmt.Sprintf("%d", stock)
}

// controller pelanggan
// CreateCustomer handles creating a new customer
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCustomer.CreatedAt = time.Now()
	newCustomer.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.CustomerCollection.InsertOne(ctx, newCustomer)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer created successfully"})
}

// GetCustomers handles retrieving all customers
func GetCustomers(w http.ResponseWriter, r *http.Request) {
	var customers []model.Customer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.CustomerCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var cust model.Customer
		if err := cursor.Decode(&cust); err != nil {
			http.Error(w, "Error decoding customer", http.StatusInternalServerError)
			return
		}
		customers = append(customers, cust)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// GetCustomerByID handles retrieving a customer by ID
func GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var customer model.Customer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = config.CustomerCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&customer)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// UpdateCustomer handles updating a customer by ID
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	// Get the customer ID from URL query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	// Convert the string ID to ObjectID (assuming MongoDB)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var updatedCustomer model.Customer
	// Decode the JSON body to the updatedCustomer struct
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Context with a timeout for MongoDB operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Update the customer fields in MongoDB
	update := bson.M{
		"$set": bson.M{
			"name":      updatedCustomer.Name,
			"email":     updatedCustomer.Email,
			"phone":     updatedCustomer.Phone,
			"address":   updatedCustomer.Address,
			"updatedAt": time.Now(),
		},
	}

	// Perform the update operation in MongoDB
	_, err = config.CustomerCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer updated successfully"})
}


// DeleteCustomer handles deleting a customer by ID
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Get the customer ID from URL query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	// Convert the string ID to ObjectID (assuming MongoDB)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Context with a timeout for MongoDB operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete the customer document from the database
	_, err = config.CustomerCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer deleted successfully"})
}




// Handler Laporan
// Handler untuk membuat laporan keuangan
func CreateFinancialReport(w http.ResponseWriter, r *http.Request) {
	var newReport model.LaporanAkuntan

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&newReport); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse tanggal dari string ke time.Time
	startDate, err := time.Parse("2006-01-02", newReport.StartDate) // Mengambil dari objek
	if err != nil {
		http.Error(w, "Invalid start date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", newReport.EndDate) // Mengambil dari objek
	if err != nil {
		http.Error(w, "Invalid end date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Set waktu pembuatan
	newReport.StartDateTime = startDate // Pastikan field ini ada di model
	newReport.EndDateTime = endDate     // Pastikan field ini ada di model
	newReport.CreatedAt = time.Now()

	// Simpan ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newReport.ID = primitive.NewObjectID()
	_, err = config.ReportCollection.InsertOne(ctx, newReport)
	if err != nil {
		http.Error(w, "Failed to create financial report", http.StatusInternalServerError)
		return
	}

	// Kirim respon sukses
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Financial report created successfully"})
}

// Fungsi untuk mendapatkan laporan keuangan berdasarkan ID
func GetFinancialReportByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	var report model.LaporanAkuntan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = config.ReportCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&report)
	if err != nil {
		http.Error(w, "Financial report not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// Handler untuk mendapatkan semua laporan keuangan
func GetFinancialReports(w http.ResponseWriter, r *http.Request) {
	var reports []model.LaporanAkuntan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.ReportCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch financial reports", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var rep model.LaporanAkuntan
		if err := cursor.Decode(&rep); err != nil {
			http.Error(w, "Error decoding financial report", http.StatusInternalServerError)
			return
		}
		reports = append(reports, rep)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// Fungsi untuk menghapus laporan keuangan berdasarkan ID
func DeleteFinancialReport(w http.ResponseWriter, r *http.Request) {
	// Get the report ID from URL query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Report ID is required", http.StatusBadRequest)
		return
	}

	// Convert the string ID to ObjectID (assuming MongoDB)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	// Context with a timeout for MongoDB operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Delete the report document from the database
	_, err = config.ReportCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete financial report", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Financial report deleted successfully"})
}
