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

// Fungsi untuk menambahkan transaksi pengeluaran baru
func CreateExpenseTransaction(respw http.ResponseWriter, req *http.Request) {
	var expense model.ExpenseTransaction
	if err := json.NewDecoder(req.Body).Decode(&expense); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Inisialisasi data transaksi pengeluaran baru
	expense.ID = primitive.NewObjectID()
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()

	_, err := atdb.InsertOneDoc(config.Mongoconn, "expense_transaction", expense)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Transaksi pengeluaran berhasil ditambahkan",
		"data":    expense,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Fungsi untuk mendapatkan daftar semua transaksi pengeluaran
func GetExpenses(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.ExpenseTransaction](config.Mongoconn, "expense_transaction", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data pengeluaran tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Format hasil sebagai slice of map untuk tiap transaksi pengeluaran
	var expenses []map[string]interface{}
	for _, expense := range data {
		expenses = append(expenses, map[string]interface{}{
			"id":             expense.ID,
			"expense_name":   expense.ExpenseName,
			"amount":         expense.Amount,
			"category":       expense.Category,
			"payment_method": expense.PaymentMethod,
			"expense_date":   expense.ExpenseDate,
			"notes":          expense.Notes,
			"created_at":     expense.CreatedAt,
			"updated_at":     expense.UpdatedAt,
		})
	}

	at.WriteJSON(respw, http.StatusOK, expenses)
}

// Fungsi untuk mendapatkan detail transaksi pengeluaran berdasarkan ID
func GetExpenseByID(respw http.ResponseWriter, req *http.Request) {
	expenseID := req.URL.Query().Get("id")
	if expenseID == "" {
		var respn model.Response
		respn.Status = "Error: ID Pengeluaran tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Pengeluaran tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var expense model.ExpenseTransaction
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.ExpenseTransaction](config.Mongoconn, "expense_transaction", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Pengeluaran tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pengeluaran ditemukan",
		"data":    expense,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Fungsi untuk mengupdate transaksi pengeluaran berdasarkan ID
func UpdateExpense(respw http.ResponseWriter, req *http.Request) {
	expenseID := req.URL.Query().Get("id")
	if expenseID == "" {
		var respn model.Response
		respn.Status = "Error: ID Pengeluaran tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Pengeluaran tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var updatedExpense model.ExpenseTransaction
	if err := json.NewDecoder(req.Body).Decode(&updatedExpense); err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	updateData := bson.M{
		"expense_name":   updatedExpense.ExpenseName,
		"amount":         updatedExpense.Amount,
		"category":       updatedExpense.Category,
		"payment_method": updatedExpense.PaymentMethod,
		"expense_date":   updatedExpense.ExpenseDate,
		"notes":          updatedExpense.Notes,
		"updated_at":     time.Now(),
	}

	filter := bson.M{"_id": objectID}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "expense_transaction", filter, bson.M{"$set": updateData})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate pengeluaran"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pengeluaran berhasil diupdate",
		"data":    updateData,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Fungsi untuk menghapus transaksi pengeluaran berdasarkan ID
func DeleteExpense(respw http.ResponseWriter, req *http.Request) {
	expenseID := req.URL.Query().Get("id")
	if expenseID == "" {
		var respn model.Response
		respn.Status = "Error: ID Pengeluaran tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(expenseID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Pengeluaran tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "expense_transaction", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus pengeluaran"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	if deleteResult.DeletedCount == 0 {
		var respn model.Response
		respn.Status = "Error: Pengeluaran tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pengeluaran berhasil dihapus",
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Fungsi untuk mengekspor data pengeluaran ke CSV
// func ExportExpensesToCSV(respw http.ResponseWriter, req *http.Request) {
// 	var expenses []model.ExpenseTransaction
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	cursor, err := config.ExpenseTransactionCollection.Find(ctx, bson.M{})
// 	if err != nil {
// 		var respn model.Response
// 		respn.Status = "Error: Gagal mengambil data pengeluaran"
// 		at.WriteJSON(respw, http.StatusInternalServerError, respn)
// 		return
// 	}
// 	defer cursor.Close(ctx)

// 	for cursor.Next(ctx) {
// 		var expense model.ExpenseTransaction
// 		if err := cursor.Decode(&expense); err != nil {
// 			var respn model.Response
// 			respn.Status = "Error: Gagal mendekode data pengeluaran"
// 			at.WriteJSON(respw, http.StatusInternalServerError, respn)
// 			return
// 		}
// 		expenses = append(expenses, expense)
// 	}

// 	respw.Header().Set("Content-Disposition", "attachment; filename=expenses.csv")
// 	respw.Header().Set("Content-Type", "text/csv")

// 	csvWriter := csv.NewWriter(respw)
// 	defer csvWriter.Flush()

// 	headers := []string{"ID", "Expense Name", "Amount", "Category", "Payment Method", "Expense Date", "Notes", "Created At", "Updated At"}
// 	if err := csvWriter.Write(headers); err != nil {
// 		var respn model.Response
// 		respn.Status = "Error: Gagal menulis header CSV"
// 		at.WriteJSON(respw, http.StatusInternalServerError, respn)
// 		return
// 	}

// 	for _, expense := range expenses {
// 		record := []string{
// 			expense.ID.Hex(),
// 			expense.ExpenseName,
// 			fmt.Sprintf("%.2f", expense.Amount),
// 			expense.Category,
// 			expense.PaymentMethod,
// 			expense.ExpenseDate.String(),
// 			expense.Notes,
// 			expense.CreatedAt.String(),
// 			expense.UpdatedAt.String(),
// 		}
// 		if err := csvWriter.Write(record); err != nil {
// 			var respn model.Response
// 			respn.Status = "Error: Gagal menulis data CSV"
// 			at.WriteJSON(respw, http.StatusInternalServerError, respn)
// 			return
// 		}
// 	}
// }
