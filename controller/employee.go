package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateEmployee handles creating a new employee
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var newEmployee model.Employee
	if err := json.NewDecoder(r.Body).Decode(&newEmployee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mengatur ID otomatis dan timestamps
	newEmployee.ID = primitive.NewObjectID()
	newEmployee.CreatedAt = time.Now()
	newEmployee.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.EmployeeCollection.InsertOne(ctx, newEmployee)
	if err != nil {
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee created successfully"})
}


// GetEmployees returns all employees
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	var employees []model.Employee

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.EmployeeCollection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var emp model.Employee
		if err := cursor.Decode(&emp); err != nil {
			http.Error(w, "Error decoding employee", http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

// GetEmployeeByID retrieves an employee by ID
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var employee model.Employee
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = config.EmployeeCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&employee)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee updates an employee by ID
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var updatedEmployee model.Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Failed to update employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee updated successfully"})
}

// DeleteEmployee deletes an employee by ID
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = config.EmployeeCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Failed to delete employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee deleted successfully"})
}


