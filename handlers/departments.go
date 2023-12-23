package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	employeedb "github.com/Aabdelm/employee-app/database"
	"github.com/go-chi/chi/v5"
)

type DepartmentHandler struct {
	L     *log.Logger
	DbMap employeedb.DeptMapper
}

func NewDepartmentHandler(l *log.Logger, DbMap employeedb.DeptMapper) DepartmentHandler {
	return DepartmentHandler{
		L:     l,
		DbMap: DbMap,
	}
}

func (dh DepartmentHandler) PostDepartment(rw http.ResponseWriter, r *http.Request) {
	var err error

	dept := employeedb.NewEmployeeDepartment()

	j := json.NewDecoder(r.Body)

	if err := j.Decode(dept); err != nil {
		dh.L.Printf("[ERROR] Failed to parse json. Error: %s\n", err)
		http.Error(rw, "Failed to decode json", http.StatusBadRequest)
		return
	}

	dh.L.Printf("[INFO] Decoded json: ID %d\n Department: %s\n", dept.Id, dept.Department)

	err = dh.DbMap.AddNewDepartment(dept)

	if err != nil {
		dh.L.Println("[ERROR] Failed to insert department")
		http.Error(rw, "Failed to insert department", http.StatusBadRequest)
		return
	}
	dh.L.Printf("[INFO] Added new department %s", dept.Department)

}

func (dh DepartmentHandler) PutDepartment(rw http.ResponseWriter, r *http.Request) {
	var err error

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		dh.L.Printf("[ERROR] Failed to convert id to integer. Error %s", err)
		http.Error(rw, "Failed to convert id to integer", http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	newDept := employeedb.NewEmployeeDepartment()

	err = dec.Decode(newDept)

	if err != nil {
		dh.L.Printf("[ERROR] Failed to decode JSON. Error %s", err)
		http.Error(rw, "Failed to convert JSON. Data might be malformed", http.StatusBadRequest)
		return
	}

	err = dh.DbMap.UpdateDepartment(id, newDept, newDept.Department)
	if err != nil {
		dh.L.Printf("[ERROR] Failed to update department %d. Error %s", id, err)
		http.Error(rw, "Failed to update department", http.StatusInternalServerError)
		return
	}

	dh.L.Printf("Updated department %d", id)

}

func (dh DepartmentHandler) DeleteDepartment(rw http.ResponseWriter, r *http.Request) {
	var err error

	idS := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idS)

	if err != nil {
		dh.L.Printf("[ERROR] Failed to parse int from URL. Error %s", err)
		http.Error(rw, `Failed to parse int from URL. 
		This might be due to the parameter being a string`, http.StatusBadRequest)
		return
	}

	err = dh.DbMap.RemoveDepartment(id)
	if err != nil {
		dh.L.Printf("[ERROR] Failed to delete id %d", id)
		http.Error(rw, "Failed to delete department", http.StatusInternalServerError)
		return
	}

	dh.L.Printf("[INFO] Deleted id %d", id)
}

func (dh DepartmentHandler) GetDepartment(rw http.ResponseWriter, r *http.Request) {
	var err error
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		dh.L.Printf("[ERROR] Failed to convert URL parameter to integer. Error %s", err)
		http.Error(rw, `Failed to convert URL parameter to integer. 
		Your URL might contain a string rather than an integer`, http.StatusBadRequest)
		return
	}

	emps, err := dh.DbMap.GetEmployeesByDepartment(id)
	if err != nil {
		dh.L.Printf("[ERROR] Failed to get employee department. Error %s", err)
		http.Error(rw, "Failed to retrieve employee department", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(rw)
	err = enc.Encode(emps)

	if err != nil {
		dh.L.Printf("[ERROR] Failed to encode JSON for employees. Error %s", err)
		http.Error(rw, "Failed to encode JSON for employees", http.StatusInternalServerError)
		return
	}
	dh.L.Printf("[INFO] Successfully retrieved Employeed for dept %d", id)

}
