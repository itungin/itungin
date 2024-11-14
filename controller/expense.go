package controller

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"log"

	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateExpenseTransaction adalah handler untuk membuat transaksi pengeluaran
func CreateExpenseTransaction(w http.ResponseWriter, r *http.Request) {
	var newExpense model.ExpenseTransaction
	if err := json.NewDecoder(r.Body).Decode(&newExpense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newExpense.CreatedAt = time.Now()
	newExpense.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.ExpenseTransactionCollection.InsertOne(ctx, newExpense)
	if err != nil {
		http.Error(w, "Failed to create expense transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Expense transaction created successfully"})
}

// GetExpenses mengembalikan semua transaksi pengeluaran
func GetExpenses(w http.ResponseWriter, r *http.Request) {
	var expenses []model.ExpenseTransaction

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.ExpenseTransactionCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch expense transactions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var exp model.ExpenseTransaction
		if err := cursor.Decode(&exp); err != nil {
			http.Error(w, "Error decoding expense transaction", http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, exp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// GetExpenseByID mengambil transaksi pengeluaran berdasarkan ID dari query parameter
func GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari query parameter
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Sales transaction ID is required", http.StatusBadRequest)
		return
	}

	log.Println("ID received:", id) // Tambahkan log untuk debugging

	// Konversi ID dari string ke ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid sales transaction ID", http.StatusBadRequest)
		return
	}

	// Ambil transaksi dari MongoDB
	var transaction model.SalesTransaction
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = config.SalesTransactionCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&transaction)
	if err != nil {
		http.Error(w, "Sales transaction not found", http.StatusNotFound)
		return
	}

	// Kirim transaksi sebagai respon
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}


// UpdateExpense mengupdate transaksi pengeluaran berdasarkan ID dari query parameter
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid expense transaction ID", http.StatusBadRequest)
		return
	}

	var updatedExpense model.ExpenseTransaction
	if err := json.NewDecoder(r.Body).Decode(&updatedExpense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"expense_name":  updatedExpense.ExpenseName,
			"amount":        updatedExpense.Amount,
			"category":      updatedExpense.Category,
			"payment_method": updatedExpense.PaymentMethod,
			"expense_date":  updatedExpense.ExpenseDate,
			"notes":         updatedExpense.Notes,
			"updated_at":    time.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.ExpenseTransactionCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, "Failed to update expense transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Expense transaction updated successfully"})
}

// DeleteExpense menghapus transaksi pengeluaran berdasarkan ID dari query parameter
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid expense transaction ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.ExpenseTransactionCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete expense transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Expense transaction deleted successfully"})
}

// ExportExpensesToCSV mengekspor semua transaksi pengeluaran ke dalam file CSV
func ExportExpensesToCSV(w http.ResponseWriter, r *http.Request) {
	var expenses []model.ExpenseTransaction

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.ExpenseTransactionCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch expense transactions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var exp model.ExpenseTransaction
		if err := cursor.Decode(&exp); err != nil {
			http.Error(w, "Error decoding expense transaction", http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, exp)
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=expenses.csv")

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Menulis header CSV
	csvWriter.Write([]string{"ID", "Expense Name", "Amount", "Category", "Payment Method", "Expense Date", "Notes", "Created At", "Updated At"})

	// Menulis data transaksi ke CSV
	for _, exp := range expenses {
		csvWriter.Write([]string{
			exp.ID.Hex(),
			exp.ExpenseName,
			fmt.Sprintf("%.2f", exp.Amount),
			exp.Category,
			exp.PaymentMethod,
			exp.ExpenseDate.Format("2006-01-02"),
			exp.Notes,
			exp.CreatedAt.Format("2006-01-02 15:04:05"),
			exp.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}
