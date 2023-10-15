package handlers

import (
	"encoding/json"
	"fmt"
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

func (employeeHandler EmployeeHandler) GetEmployee(rw http.ResponseWriter, r *http.Request) {
	s := chi.URLParam(r, "id")
	i, _ := strconv.Atoi(s)
	fmt.Fprintf(rw, "Hello there, employee %d", i)
}

func (employeeHandler EmployeeHandler) AddEmployee(rw http.ResponseWriter, r *http.Request) {
	employee := employeedb.NewEmployee()
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(employee)
	employeeHandler.DbMap.AddNewEmployee(employee)
}
