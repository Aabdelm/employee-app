package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	employeedb "github.com/Aabdelm/employee-app/database"
	"github.com/go-chi/chi/v5"
)

type EmployeeHandler struct {
	L     *log.Logger
	DbMap employeedb.EmployeeMapper
}

func NewEmployeeHandler(l *log.Logger, DbMap employeedb.EmployeeMapper) EmployeeHandler {

	return EmployeeHandler{
		L:     l,
		DbMap: DbMap,
	}
}

func (eh EmployeeHandler) GetEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error
	idString := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to get parameter. Error: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	eh.L.Println("Got id", id)

	employee, err := eh.DbMap.GetEmployeeById(id)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to Get employee %d. Error: %s", id, err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if employee == nil {
		eh.L.Printf("[ERROR] Failed to Get employee %d. Employee is null", id)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	enc := json.NewEncoder(rw)

	err = enc.Encode(employee)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to get JSON for employee %d. Error: %s", id, err)
		return
	}

}

func (eh EmployeeHandler) PostEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error
	employee := employeedb.NewEmployee()

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(employee)

	eh.L.Printf("[INFO] posting new employee %d", employee.DepartmentId)

	err = eh.DbMap.AddNewEmployee(employee)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (eh EmployeeHandler) PutEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error
	param := chi.URLParam(r, "id")

	id, err := strconv.Atoi(param)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		eh.L.Printf("[ERROR] Failed to convert to integer on PutEmployee. Error %s", err)
		return
	}

	employee := employeedb.NewEmployee()
	employee.Id = id

	dec := json.NewDecoder(r.Body)
	if err = dec.Decode(employee); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		eh.L.Printf("[ERROR] Failed to decode json for employee %d. Error: %s", id, err)
		return
	}

	if err = eh.DbMap.UpdateEmployee(employee); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		eh.L.Printf("[ERROR] Failed to update employee %d. Error: %s", id, err)
		return
	}

	eh.L.Printf("[INFO] updated %d", id)

	enc := json.NewEncoder(rw)

	err = enc.Encode(employee)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		eh.L.Printf("[ERROR] Failed to update employee %d. Error %s", id, err)
		return
	}

	eh.L.Printf("[INFO] Updated %d", id)

}

func (eh EmployeeHandler) DeleteEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error

	idS := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idS)

	if err != nil {
		eh.L.Printf("[ERROR] Failed to convert id into integer. Error %s", err)
		http.Error(rw, "Failed to convert id into integer", http.StatusBadRequest)
		return
	}

	err = eh.DbMap.DeleteEmployee(id)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to delete employee %d. Error %s", id, err)
		http.Error(rw, "Failed to delete employee", http.StatusInternalServerError)
		return
	}

	eh.L.Printf("[INFO] Deleted employee %d", id)

}
