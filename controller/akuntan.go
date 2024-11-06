package controller

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


// Fungsi untuk menambahkan produk baru
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct model.Product // Model Produk
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set waktu pembuatan produk
	newProduct.CreatedAt = time.Now()

	// Insert produk ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.ProductCollection.InsertOne(ctx, newProduct)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Kirim respon sukses
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product created successfully"})
}

// Fungsi untuk mendapatkan daftar produk
func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []model.Product

	// Ambil data dari MongoDB
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

	// Kirim data produk sebagai respon
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Fungsi untuk mendapatkan detail produk berdasarkan ID
func GetProductByID(w http.ResponseWriter, r *http.Request) {
    // Ambil parameter ID dari URL menggunakan mux.Vars
    vars := mux.Vars(r)
    id := vars["id"]

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    // Ambil produk dari MongoDB
    var product model.Product
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err = config.ProductCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
    if err != nil {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    // Kirim produk sebagai respon
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(product)
}


// Fungsi untuk mengupdate produk berdasarkan ID
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL menggunakan gorilla/mux
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Decode data produk yang akan diupdate
	var updatedProduct model.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update produk di MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":        updatedProduct.Name,
			"description": updatedProduct.Description,
			"price":       updatedProduct.Price,
			"category":    updatedProduct.Category,
			"stock":       updatedProduct.Stock,
			"updatedAt":   time.Now(),
		},
	}

	_, err = config.ProductCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	// Kirim respon sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully"})
}



// Fungsi untuk menghapus produk berdasarkan ID
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL menggunakan gorilla/mux
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Hapus produk dari MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.ProductCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	// Kirim respon sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted successfully"})
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
	// Ambil parameter ID dari URL menggunakan gorilla/mux
	vars := mux.Vars(r)
	id := vars["id"]
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var updatedCustomer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&updatedCustomer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":      updatedCustomer.Name,
			"email":     updatedCustomer.Email,
			"phone":     updatedCustomer.Phone,
			"address":   updatedCustomer.Address,
			"updatedAt": time.Now(),
		},
	}

	_, err = config.CustomerCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Customer updated successfully"})
}

// DeleteCustomer handles deleting a customer by ID
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ID dari URL menggunakan gorilla/mux
	vars := mux.Vars(r)
	id := vars["id"]
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.CustomerCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

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
	// Ambil parameter ID dari URL menggunakan gorilla/mux
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.ReportCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete financial report", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Financial report deleted successfully"})
}