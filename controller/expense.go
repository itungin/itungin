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
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Fungsi untuk menambahkan transaksi pengeluaran baru
func CreateExpenseTransaction(w http.ResponseWriter, r *http.Request) {
	var expense model.ExpenseTransaction
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		var response model.Response
		response.Status = "Error: Bad Request"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	// Inisialisasi data transaksi pengeluaran baru
	expense.ID = primitive.NewObjectID()
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()

	_, err := config.ExpenseTransactionCollection.InsertOne(context.Background(), expense)
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
		"message": "Transaksi pengeluaran berhasil ditambahkan",
		"data":    expense,
	}
	at.WriteJSON(w, http.StatusCreated, response)
}

// Fungsi untuk mendapatkan daftar semua transaksi pengeluaran
func GetExpenses(w http.ResponseWriter, r *http.Request) {
	data, err := config.ExpenseTransactionCollection.Find(context.Background(), bson.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Data pengeluaran tidak ditemukan"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}
	defer data.Close(context.Background())

	var expenses []map[string]interface{}
	for data.Next(context.Background()) {
		var expense model.ExpenseTransaction
		if err := data.Decode(&expense); err != nil {
			var response model.Response
			response.Status = "Error: Gagal mendekode data pengeluaran"
			at.WriteJSON(w, http.StatusInternalServerError, response)
			return
		}
		expenses = append(expenses, map[string]interface{}{
			"id":            expense.ID,
			"expense_name":  expense.ExpenseName,
			"amount":        expense.Amount,
			"category":      expense.Category,
			"payment_method": expense.PaymentMethod,
			"expense_date":  expense.ExpenseDate,
			"notes":         expense.Notes,
			"created_at":    expense.CreatedAt,
			"updated_at":    expense.UpdatedAt,
		})
	}

	// Kirim data transaksi pengeluaran sebagai respon
	at.WriteJSON(w, http.StatusOK, expenses)
}

// Fungsi untuk mendapatkan detail transaksi pengeluaran berdasarkan ID
func GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	expenseID := r.URL.Query().Get("id")
	if expenseID == "" {
		var response model.Response
		response.Status = "Error: ID Pengeluaran tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Pengeluaran tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	var expense model.ExpenseTransaction
	err = config.ExpenseTransactionCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&expense)
	if err != nil {
		var response model.Response
		response.Status = "Error: Pengeluaran tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pengeluaran ditemukan",
		"data":    expense,
	}
	at.WriteJSON(w, http.StatusOK, response)
}

// Fungsi untuk mengupdate transaksi pengeluaran berdasarkan ID
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	expenseID := r.URL.Query().Get("id")
	if expenseID == "" {
		var response model.Response
		response.Status = "Error: ID Pengeluaran tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Pengeluaran tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	var updatedExpense model.ExpenseTransaction
	if err := json.NewDecoder(r.Body).Decode(&updatedExpense); err != nil {
		var response model.Response
		response.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	updateData := bson.M{
		"expense_name":  updatedExpense.ExpenseName,
		"amount":        updatedExpense.Amount,
		"category":      updatedExpense.Category,
		"payment_method": updatedExpense.PaymentMethod,
		"expense_date":  updatedExpense.ExpenseDate,
		"notes":         updatedExpense.Notes,
		"updated_at":    time.Now(),
	}

	_, err = config.ExpenseTransactionCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updateData})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengupdate pengeluaran"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pengeluaran berhasil diupdate",
		"data":    updateData,
	}
	at.WriteJSON(w, http.StatusOK, response)
}

// Fungsi untuk menghapus transaksi pengeluaran berdasarkan ID
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	expenseID := r.URL.Query().Get("id")
	if expenseID == "" {
		var response model.Response
		response.Status = "Error: ID Pengeluaran tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID Pengeluaran tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	deleteResult, err := config.ExpenseTransactionCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal menghapus pengeluaran"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	if deleteResult.DeletedCount == 0 {
		var response model.Response
		response.Status = "Error: Pengeluaran tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pengeluaran berhasil dihapus",
	}
	at.WriteJSON(w, http.StatusOK, response)
}

// Fungsi untuk mengekspor data pengeluaran ke CSV
func ExportExpensesToCSV(w http.ResponseWriter, r *http.Request) {
	var expenses []model.ExpenseTransaction
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.ExpenseTransactionCollection.Find(ctx, bson.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengambil data pengeluaran"
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var expense model.ExpenseTransaction
		if err := cursor.Decode(&expense); err != nil {
			var response model.Response
			response.Status = "Error: Gagal mendekode data pengeluaran"
			at.WriteJSON(w, http.StatusInternalServerError, response)
			return
		}
		expenses = append(expenses, expense)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=expenses.csv")
	w.Header().Set("Content-Type", "text/csv")

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	headers := []string{"ID", "Expense Name", "Amount", "Category", "Payment Method", "Expense Date", "Notes", "Created At", "Updated At"}
	if err := csvWriter.Write(headers); err != nil {
		var response model.Response
		response.Status = "Error: Gagal menulis header CSV"
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	for _, expense := range expenses {
		row := []string{
			expense.ID.Hex(),
			expense.ExpenseName,
			fmt.Sprintf("%.2f", expense.Amount),
			expense.Category,
			expense.PaymentMethod,
			expense.ExpenseDate.Format(time.RFC3339),
			expense.Notes,
			expense.CreatedAt.Format(time.RFC3339),
			expense.UpdatedAt.Format(time.RFC3339),
		}
		if err := csvWriter.Write(row); err != nil {
			var response model.Response
			response.Status = "Error: Gagal menulis data CSV"
			at.WriteJSON(w, http.StatusInternalServerError, response)
			return
		}
	}
}
