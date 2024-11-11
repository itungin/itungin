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

	cursor, err := config.Mongoconn.Collection("products").Find(ctx, bson.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengambil data produk"
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var prod model.Product
		if err := cursor.Decode(&prod); err != nil {
			var response model.Response
			response.Status = "Error: Gagal mendekode data produk"
			at.WriteJSON(w, http.StatusInternalServerError, response)
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
		var response model.Response
		response.Status = "Error: Gagal menulis header CSV"
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Tulis data produk ke CSV
	for _, product := range products {
		row := []string{
			product.ID.Hex(),
			product.Name,
			formatPrice(product.Price),        // Format harga
			product.Category,
			product.Description,
			formatStock(product.Stock),        // Format stok
			product.CreatedAt.Format(time.RFC3339),
			product.UpdatedAt.Format(time.RFC3339),
		}
		if err := csvWriter.Write(row); err != nil {
			var response model.Response
			response.Status = "Error: Gagal menulis data produk ke CSV"
			at.WriteJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// Kirimkan respons sukses
	var response model.Response
	response.Status = "Success: Produk berhasil diekspor ke CSV"
	at.WriteJSON(w, http.StatusOK, response)
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
	var customer model.Customer

	// Decode data pelanggan dari body permintaan
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		var response model.Response
		response.Status = "Error: Bad Request"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Inisialisasi data pelanggan baru dengan ObjectID untuk ID
	newCustomer := model.Customer{
		ID:        primitive.NewObjectID(),
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Address:   customer.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert pelanggan ke dalam MongoDB
	_, err := atdb.InsertOneDoc(config.Mongoconn, "customers", newCustomer)
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
		"message": "Pelanggan berhasil ditambahkan",
		"data":    newCustomer,
	}
	at.WriteJSON(w, http.StatusCreated, response)
}


// GetCustomers handles retrieving all customers
func GetCustomers(w http.ResponseWriter, r *http.Request) {
	// Ambil semua data pelanggan dari MongoDB
	data, err := atdb.GetAllDoc[[]model.Customer](config.Mongoconn, "customers", primitive.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Data pelanggan tidak ditemukan"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	if len(data) == 0 {
		var response model.Response
		response.Status = "Error: Data pelanggan kosong"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Format hasil sebagai slice of map dengan ID, Name, Email, Phone, Address, CreatedAt, dan UpdatedAt untuk setiap pelanggan
	var customers []map[string]interface{}
	for _, customer := range data {
		customers = append(customers, map[string]interface{}{
			"id":        customer.ID,
			"name":      customer.Name,
			"email":     customer.Email,
			"phone":     customer.Phone,
			"address":   customer.Address,
			"createdAt": customer.CreatedAt,
			"updatedAt": customer.UpdatedAt,
		})
	}

	// Kirim data pelanggan sebagai respon
	at.WriteJSON(w, http.StatusOK, customers)
}



// GetCustomerByID handles retrieving a customer by ID
func GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL
	customerID := r.URL.Query().Get("id")
	if customerID == "" {
		var response model.Response
		response.Status = "Error: ID Pelanggan tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Konversi ID dari string ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Pelanggan tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Ambil data pelanggan dari MongoDB
	var customer model.Customer
	filter := bson.M{"_id": objectID}
	err = config.Mongoconn.Collection("customers").FindOne(context.TODO(), filter).Decode(&customer)
	if err != nil {
		var response model.Response
		response.Status = "Error: Pelanggan tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Kirim data pelanggan sebagai respon
	response := map[string]interface{}{
		"status":  "success",
		"message": "Pelanggan ditemukan",
		"data":    customer,
	}
	at.WriteJSON(w, http.StatusOK, response)
}


// UpdateCustomer handles updating a customer by ID
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL
	customerID := r.URL.Query().Get("id")
	if customerID == "" {
		var response model.Response
		response.Status = "Error: ID Pelanggan tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Konversi ID dari string ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Pelanggan tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Decode data pelanggan yang akan diupdate
	var requestBody struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
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
	if requestBody.Email != "" {
		updateData["email"] = requestBody.Email
	}
	if requestBody.Phone != "" {
		updateData["phone"] = requestBody.Phone
	}
	if requestBody.Address != "" {
		updateData["address"] = requestBody.Address
	}
	updateData["updatedAt"] = time.Now()

	// Update pelanggan di MongoDB
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updateData}
	_, err = config.Mongoconn.Collection("customers").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengupdate pelanggan"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusNotModified, response)
		return
	}

	// Kirim respon sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Pelanggan berhasil diupdate",
		"data":    updateData,
	}
	at.WriteJSON(w, http.StatusOK, response)
}



// DeleteCustomer handles deleting a customer by ID
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL
	customerID := r.URL.Query().Get("id")
	if customerID == "" {
		var response model.Response
		response.Status = "Error: ID Pelanggan tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Konversi customerID dari string ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Pelanggan tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Hapus data pelanggan berdasarkan ID
	filter := bson.M{"_id": objectID}
	deleteResult, err := config.Mongoconn.Collection("customers").DeleteOne(context.TODO(), filter)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal menghapus pelanggan"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Periksa apakah ada pelanggan yang dihapus
	if deleteResult.DeletedCount == 0 {
		var response model.Response
		response.Status = "Error: Pelanggan tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Kirim respon sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Pelanggan berhasil dihapus",
		"data":    deleteResult,
	}
	at.WriteJSON(w, http.StatusOK, response)
}






// Handler Laporan
// Handler untuk membuat laporan keuangan
func CreateFinancialReport(w http.ResponseWriter, r *http.Request) {
	var report model.LaporanAkuntan

	// Decode data laporan dari body permintaan
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		var response model.Response
		response.Status = "Error: Bad Request"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Parse tanggal dari string ke time.Time
	startDate, err := time.Parse("2006-01-02", report.StartDate)
	if err != nil {
		var response model.Response
		response.Status = "Error: Invalid start date format. Use YYYY-MM-DD"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}
	endDate, err := time.Parse("2006-01-02", report.EndDate)
	if err != nil {
		var response model.Response
		response.Status = "Error: Invalid end date format. Use YYYY-MM-DD"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Inisialisasi data laporan baru dengan ObjectID untuk ID dan mengisi income, expenses, profit
	newReport := model.LaporanAkuntan{
		ID:             primitive.NewObjectID(),
		StartDate:      report.StartDate,
		EndDate:        report.EndDate,
		StartDateTime:  startDate,
		EndDateTime:    endDate,
		Income:         report.Income,  // Pastikan income dikirim dalam JSON
		Expenses:       report.Expenses, // Pastikan expenses dikirim dalam JSON
		Profit:         report.Profit, // Pastikan profit dikirim dalam JSON
		CreatedAt:      time.Now(),
	}

	// Insert laporan ke dalam MongoDB
	_, err = atdb.InsertOneDoc(config.Mongoconn, "financial_reports", newReport)
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
		"message": "Laporan keuangan berhasil dibuat",
		"data":    newReport,
	}
	at.WriteJSON(w, http.StatusCreated, response)
}



// Fungsi untuk mendapatkan laporan keuangan berdasarkan ID
func GetFinancialReportByID(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL
	reportID := r.URL.Query().Get("id")
	if reportID == "" {
		var response model.Response
		response.Status = "Error: ID Laporan tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Konversi ID dari string ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Laporan tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Ambil data laporan keuangan dari MongoDB
	var report model.LaporanAkuntan
	filter := bson.M{"_id": objectID}
	err = config.Mongoconn.Collection("financial_reports").FindOne(context.TODO(), filter).Decode(&report)
	if err != nil {
		var response model.Response
		response.Status = "Error: Laporan keuangan tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Kirim data laporan sebagai respon
	response := map[string]interface{}{
		"status":  "success",
		"message": "Laporan keuangan ditemukan",
		"data":    report,
	}
	at.WriteJSON(w, http.StatusOK, response)
}


// Handler untuk mendapatkan semua laporan keuangan
func GetFinancialReports(w http.ResponseWriter, r *http.Request) {
	// Ambil semua data laporan keuangan dari MongoDB
	data, err := atdb.GetAllDoc[[]model.LaporanAkuntan](config.Mongoconn, "financial_reports", primitive.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Data laporan keuangan tidak ditemukan"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	if len(data) == 0 {
		var response model.Response
		response.Status = "Error: Data laporan keuangan kosong"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Format hasil sebagai slice of map dengan ID, StartDate, EndDate, Income, Expenses, Profit, CreatedAt
	var reports []map[string]interface{}
	for _, report := range data {
		reports = append(reports, map[string]interface{}{
			"id":             report.ID,
			"startDate":      report.StartDate,
			"endDate":        report.EndDate,
			"income":         report.Income,
			"expenses":       report.Expenses,
			"profit":         report.Profit,
			"createdAt":      report.CreatedAt,
		})
	}

	// Kirim data laporan keuangan sebagai respon
	at.WriteJSON(w, http.StatusOK, reports)
}


// Fungsi untuk menghapus laporan keuangan berdasarkan ID
func DeleteFinancialReport(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID laporan dari URL
	reportID := r.URL.Query().Get("id")
	if reportID == "" {
		var response model.Response
		response.Status = "Error: ID Laporan tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Konversi reportID dari string ke ObjectID MongoDB
	objectID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Laporan tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Hapus data laporan keuangan berdasarkan ID menggunakan atdb.DeleteOneDoc
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "financial_reports", bson.M{"_id": objectID})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal menghapus laporan keuangan"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Periksa apakah ada laporan yang dihapus
	if deleteResult.DeletedCount == 0 {
		var response model.Response
		response.Status = "Error: Laporan tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	// Kirim respon sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Laporan keuangan berhasil dihapus",
		"data":    deleteResult,
	}
	at.WriteJSON(w, http.StatusOK, response)
}

