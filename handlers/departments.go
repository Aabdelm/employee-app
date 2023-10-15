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
	l     *log.Logger
	DbMap *employeedb.DbMap
}

func NewDepartmentHandler(l *log.Logger, DbMap *employeedb.DbMap) DepartmentHandler {
	return DepartmentHandler{
		l:     l,
		DbMap: DbMap,
	}
}

func (dh DepartmentHandler) PostDepartment(rw http.ResponseWriter, r *http.Request) {
	var err error

	dept := employeedb.NewEmployeeDepartment()

	s := chi.URLParam(r, "id")

	paramId, err := strconv.Atoi(s)
	if err != nil {
		dh.l.Printf("[ERROR] Failed to convert string. Error: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	dept.Id = paramId

	j := json.NewDecoder(r.Body)

	if err := j.Decode(dept); err != nil {
		dh.l.Printf("[ERROR] Failed to parse json. Error: %s\n", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	dh.l.Printf("[INFO] Decoded json: ID %d\n Department: %s\n", dept.Id, dept.Department)

	err = dh.DbMap.AddNewDepartment(dept)

	if err != nil {
		dh.l.Println("[ERROR] Failed to insert department")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	dh.l.Printf("[INFO] Added new department %s", dept.Department)

}
