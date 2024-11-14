package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateEmployee handles creating a new employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var newEmployee model.Employee
	if err := json.NewDecoder(r.Body).Decode(&newEmployee); err != nil {
		var response model.Response
		response.Status = "Error: Bad Request"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	newEmployee.ID = primitive.NewObjectID()
	newEmployee.CreatedAt = time.Now()
	newEmployee.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.EmployeeCollection.InsertOne(ctx, newEmployee)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal membuat employee"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee created successfully",
		"data":    newEmployee,
	}
	at.WriteJSON(w, http.StatusCreated, response)
}

// GetEmployees returns all employees
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	var employees []model.Employee

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.EmployeeCollection.Find(ctx, bson.M{})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengambil data employees"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var emp model.Employee
		if err := cursor.Decode(&emp); err != nil {
			var response model.Response
			response.Status = "Error: Gagal mendekode employee"
			at.WriteJSON(w, http.StatusInternalServerError, response)
			return
		}
		employees = append(employees, emp)
	}

	at.WriteJSON(w, http.StatusOK, employees)
}

// GetEmployeeByID retrieves an employee by ID
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		var response model.Response
		response.Status = "Error: ID tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	var employee model.Employee
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = config.EmployeeCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&employee)
	if err != nil {
		var response model.Response
		response.Status = "Error: Employee tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee ditemukan",
		"data":    employee,
	}
	at.WriteJSON(w, http.StatusOK, response)
}

// UpdateEmployee updates an employee by ID
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		var response model.Response
		response.Status = "Error: ID tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	var updatedEmployee model.Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
		var response model.Response
		response.Status = "Error: Gagal membaca data JSON"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":         updatedEmployee.Name,
			"email":        updatedEmployee.Email,
			"phone_number": updatedEmployee.PhoneNumber,
			"position":     updatedEmployee.Position,
			"updated_at":   time.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.EmployeeCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal mengupdate employee"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee updated successfully",
		"data":    update,
	}
	at.WriteJSON(w, http.StatusOK, response)
}

// DeleteEmployee deletes an employee by ID
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		var response model.Response
		response.Status = "Error: ID tidak ditemukan"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		var response model.Response
		response.Status = "Error: ID tidak valid"
		at.WriteJSON(w, http.StatusBadRequest, response)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deleteResult, err := config.EmployeeCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		var response model.Response
		response.Status = "Error: Gagal menghapus employee"
		response.Response = err.Error()
		at.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	if deleteResult.DeletedCount == 0 {
		var response model.Response
		response.Status = "Error: Employee tidak ditemukan"
		at.WriteJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee deleted successfully",
	}
	at.WriteJSON(w, http.StatusOK, response)
}
