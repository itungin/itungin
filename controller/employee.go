package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/model"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateEmployee handles creating a new employee
func CreateEmployee(respw http.ResponseWriter, req *http.Request) {
	var newEmployee model.Employee
	if err := json.NewDecoder(req.Body).Decode(&newEmployee); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	newEmployee.ID = primitive.NewObjectID()
	newEmployee.CreatedAt = time.Now()
	newEmployee.UpdatedAt = time.Now()

	// Insert into MongoDB collection
	_, err := atdb.InsertOneDoc(config.Mongoconn, "employee", newEmployee)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee created successfully",
		"data":    newEmployee,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// GetAllEmployees returns all employees
func GetAllEmployees(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Employee](config.Mongoconn, "employee", bson.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data employees tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	if len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data employees kosong"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	// Format hasil sebagai slice of map dengan ID dan name untuk setiap employee
	var employees []map[string]interface{}
	for _, employee := range data {
		employees = append(employees, map[string]interface{}{
			"id":       employee.ID,
			"name":     employee.Name,
			"position": employee.Position,
		})
	}

	at.WriteJSON(respw, http.StatusOK, employees)
}

// GetEmployeeByID retrieves an employee by ID
func GetEmployeeByID(respw http.ResponseWriter, req *http.Request) {
	employeeID := req.URL.Query().Get("id")
	if employeeID == "" {
		var respn model.Response
		respn.Status = "Error: ID Employee tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(employeeID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Employee tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var employee model.Employee
	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.Employee](config.Mongoconn, "employee", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Employee tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee ditemukan",
		"data":    employee,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// UpdateEmployee updates an employee by ID
func UpdateEmployee(respw http.ResponseWriter, req *http.Request) {
	employeeID := req.URL.Query().Get("id")
	if employeeID == "" {
		var respn model.Response
		respn.Status = "Error: ID Employee tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(employeeID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Employee tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var updatedEmployee model.Employee
	if err := json.NewDecoder(req.Body).Decode(&updatedEmployee); err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal membaca data JSON"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
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

	// Perform the update
	filter := bson.M{"_id": objectID}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "employee", filter, update)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate employee"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee updated successfully",
		"data":    update,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// DeleteEmployee deletes an employee by ID
func DeleteEmployee(respw http.ResponseWriter, req *http.Request) {
	employeeID := req.URL.Query().Get("id")
	if employeeID == "" {
		var respn model.Response
		respn.Status = "Error: ID Employee tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(employeeID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Employee tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteResult, err := atdb.DeleteOneDoc(config.Mongoconn, "employee", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus employee"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	if deleteResult.DeletedCount == 0 {
		var respn model.Response
		respn.Status = "Error: Employee tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Employee deleted successfully",
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
