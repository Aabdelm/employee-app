package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	employeedb "github.com/Aabdelm/employee-app/database"
	"github.com/Aabdelm/employee-app/regex"
	"github.com/go-chi/chi/v5"
)

type EmployeeHandler struct {
	L        *log.Logger
	DbMap    employeedb.EmployeeMapper
	Searcher employeedb.Searcher
}

func NewEmployeeHandler(l *log.Logger, DbMap employeedb.EmployeeMapper, searcher employeedb.Searcher) *EmployeeHandler {

	return &EmployeeHandler{
		L:        l,
		DbMap:    DbMap,
		Searcher: searcher,
	}
}

func (eh *EmployeeHandler) GetEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error
	idString := chi.URLParam(r, "id")

	if !regex.IsDigit(idString, eh.L) {
		eh.L.Printf("[INFO] parameter is not integer")
		http.Error(rw, "ID is not a number", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to get parameter. Error: %s", err)
		http.Error(rw, "Failed to get parameter", http.StatusBadRequest)
		return
	}

	eh.L.Println("Got id", id)

	employee, err := eh.DbMap.GetEmployeeById(id)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to Get employee %d. Error: %s", id, err)
		http.Error(rw, "An error has occured", http.StatusNotFound)
		return
	}
	if employee == nil {
		eh.L.Printf("[ERROR] Failed to Get employee %d. Employee is null", id)
		http.Error(rw, "Failed to get employee", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	enc := json.NewEncoder(rw)

	err = enc.Encode(employee)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to get JSON for employee %d. Error: %s", id, err)
		http.Error(rw, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

}

func (eh *EmployeeHandler) GetEmployeesByQuery(rw http.ResponseWriter, r *http.Request) {
	eh.L.Printf("[INFO] Starting main employee function\n")

	url := r.URL.Query()
	eh.L.Printf("[INFO] Got parameters %q", url)

	queryLength := len(url)
	if queryLength != 1 {
		eh.L.Printf("[ERROR] Query received is not 1. Length received %d", queryLength)
		http.Error(rw, "Query must be singular", http.StatusBadRequest)
		return
	}

	var query string
	var dbColumn string
	switch {
	case url.Has("email"):
		query = url.Get("email")
		dbColumn = "email"
		eh.L.Printf("[INFO] Got email %s\n", query)
	case url.Has("firstname"):
		query = url.Get("firstname")
		eh.L.Printf("[INFO] Got first name %s\n", query)
		dbColumn = "first_name"
	case url.Has("lastname"):
		query = url.Get("lastname")
		eh.L.Printf("[INFO] Got last name %s\n", query)
		dbColumn = "last_name"
	case url.Has("department"):
		query = url.Get("department")
		eh.L.Printf("[INFO] Got Department %s\n", query)
		dbColumn = "department"
	default:
		http.Error(rw, "Query not executed", http.StatusBadRequest)
		return
	}

	emps, err := eh.Searcher.Search(query, dbColumn)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to get employees. Error %s", err)
		http.Error(rw, "Failed to get employees. This is most likely due to a bad query", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)

	if emps == nil {
		//In case we get a null response
		emps = make([]*employeedb.Employee, 0)
	}

	if err := enc.Encode(emps); err != nil {
		eh.L.Printf("[ERROR] Failed to encode JSON. Error %s", err)
		http.Error(rw, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
	eh.L.Printf("[INFO] successfully encoded info")

}

func (eh *EmployeeHandler) PostEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error
	employee := employeedb.NewEmployee()

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(employee)

	eh.L.Printf("[INFO] posting new employee in department %d", employee.DepartmentId)

	err = eh.DbMap.AddNewEmployee(employee)
	if err != nil {
		eh.L.Printf("[ERROR] Failed to post employee. Error %s", err)
		http.Error(rw, "Failed to add employee", http.StatusBadRequest)
		return
	}

	enc := json.NewEncoder(rw)
	err = enc.Encode(employee)
	rw.Header().Set("Content-type", "application/json")

	if err != nil {
		http.Error(rw, "Failed to encode JSON", http.StatusInternalServerError)
		eh.L.Printf("[ERROR] Failed to add employee. Error %s", err)
		return
	}

	eh.L.Printf("[INFO] Successfully added employee %d", employee.DepartmentId)
}

func (eh *EmployeeHandler) PutEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error
	eh.L.Printf("[INFO] Function PutEmployee Called\n")
	param := chi.URLParam(r, "id")

	if !regex.IsDigit(param, eh.L) {
		eh.L.Printf("[INFO] parameter is not integer")
		http.Error(rw, "ID is not a number", http.StatusBadRequest)
		return
	}

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

	enc := json.NewEncoder(rw)

	err = enc.Encode(employee)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		eh.L.Printf("[ERROR] Failed to update employee %d. Error %s", id, err)
		return
	}

	rw.Header().Set("Content-type", "application/json")

	eh.L.Printf("[INFO] Updated %d", id)

}

func (eh *EmployeeHandler) DeleteEmployee(rw http.ResponseWriter, r *http.Request) {
	var err error

	idStr := chi.URLParam(r, "id")

	if !regex.IsDigit(idStr, eh.L) {
		eh.L.Printf("[INFO] parameter is not integer")
		http.Error(rw, "ID is not a number", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)

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

func (eh *EmployeeHandler) DeleteMultipleEmployees(rw http.ResponseWriter, r *http.Request) {
	wg := &sync.WaitGroup{}

	resume, errChan := make(chan bool), make(chan error)

	eh.L.Printf("[INFO] Starting function DeleteMultipleEmployees\n")

	emps := make([]*employeedb.Employee, 0)
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&emps); err != nil {
		eh.L.Printf("[ERROR] Failed to decode JSON. Error %d", err)
		http.Error(rw, "Failed to decode JSON. This might be due to a bad input", http.StatusBadRequest)
		return
	}

	go eh.runGoRoutines(emps, wg, resume, errChan)

	select {
	case <-errChan:
		eh.L.Printf("[ERROR] Error detected in goroutine. Returning error %s", <-errChan)
		http.Error(rw, "An error was detected while attempting to delete", http.StatusInternalServerError)
		return
	case <-resume:
		eh.L.Printf("[INFO] All goroutines have finished\n")
	}

	eh.L.Printf("[INFO] Successfully deleted all employees\n")

}

func (eh *EmployeeHandler) runGoRoutines(emps []*employeedb.Employee, wg *sync.WaitGroup, resume chan bool, errChan chan error) {
	eh.L.Printf("[INFO] Function runGoRoutines called")
	//Create a buffered channel to limit the amount of goroutines running
	sem := make(chan struct{}, 20)
	for i, emp := range emps {
		wg.Add(1)

		go func(emp *employeedb.Employee, task int, wg *sync.WaitGroup) {
			sem <- struct{}{} //Add one value to block until we're finished (to limit resources)
			eh.L.Printf("[INFO] Starting task %d", task)
			if err := eh.DbMap.DeleteEmployee(emp.Id); err != nil {
				eh.L.Printf("[ERROR] Error detected at task %d, error %s", task, err)
				eh.L.Printf("Sending error to channel\n")
				errChan <- err

			}
			<-sem //Release once we're done
			eh.L.Printf("[INFO] Task %d done", task)
			wg.Done()
		}(emp, i, wg)
	}
	eh.L.Printf("[INFO] Waiting for tasks to finish")
	go func() {
		//Wait for all tasks to finish
		wg.Wait()
		eh.L.Printf("[INFO] successfully deleted all tasks")
		//We're all good now
		resume <- true
	}()
}
