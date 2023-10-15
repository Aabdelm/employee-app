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
	l     *log.Logger
	DbMap *employeedb.DbMap
}

func NewEmployeHandler(l *log.Logger, DbMap *employeedb.DbMap) EmployeeHandler {
	return EmployeeHandler{
		l:     l,
		DbMap: DbMap,
	}
}

func (eh EmployeeHandler) GetEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error
	idString := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		eh.l.Printf("[ERROR] Failed to get parameter. Error: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	eh.l.Println("Got id", id)

	employee, err := eh.DbMap.GetEmployeeById(id)
	if err != nil {
		eh.l.Printf("[ERROR] Failed to Get employee %d. Error: %s", id, err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if employee == nil {
		eh.l.Printf("[ERROR] Failed to Get employee %d. Employee is null", id)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	enc := json.NewEncoder(rw)

	err = enc.Encode(employee)
	if err != nil {
		eh.l.Printf("[ERROR] Failed to get JSON for employee %d. Error: %s", id, err)
		return
	}

}

func (eh EmployeeHandler) PostEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error
	employee := employeedb.NewEmployee()

	routeString := chi.URLParam(r, "id")

	id, err := strconv.Atoi(routeString)
	if err != nil {
		eh.l.Printf("[ERROR] failed to parse integer from string. error: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	employee.Id = id

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(employee)

	eh.l.Printf("[INFO] got dept id %d", employee.DepartmentId)

	err = eh.DbMap.AddNewEmployee(employee)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

}
