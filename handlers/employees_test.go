package handlers_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	employeedb "github.com/Aabdelm/employee-app/database"
	"github.com/Aabdelm/employee-app/handlers"
	"github.com/go-chi/chi/v5"
)

// mock database
type mockDb map[int]*employeedb.Employee

/*
Mock methods to implement the EmployeeMapper interface
*/

func (mock mockDb) GetEmployeeById(id int) (*employeedb.Employee, error) {
	_, ok := mock[id]
	if !ok {
		return nil, errors.New("Error: id not present")
	}
	return mock[id], nil
}

func (mock mockDb) UpdateEmployee(emp *employeedb.Employee) error {
	_, ok := mock[emp.Id]
	if !ok {
		return errors.New("Error: id not present")
	}

	mock[emp.Id] = emp

	return nil
}

func (mock mockDb) DeleteEmployee(id int) error {
	_, ok := mock[id]
	if !ok {
		return errors.New("Error: id not present")
	}
	delete(mock, id)
	return nil
}
func (mock mockDb) AddNewEmployee(emp *employeedb.Employee) error {
	_, ok := mock[emp.Id]
	if ok {
		return errors.New("Error: ID already present")
	}
	mock[emp.Id] = emp

	return nil
}

func TestGetEmployee(t *testing.T) {
	testDb := make(mockDb)
	logger := log.Default()

	testDb[1] = &employeedb.Employee{}

	eh := &handlers.EmployeeHandler{
		L:     logger,
		DbMap: testDb,
	}
	//Setup new recorder
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/employees", nil)

	//Setup new context for id
	//Think of it as /employees/1
	//The following is needed to satisfy chi (since it uses the context for the id)
	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	eh.GetEmployee(rr, req)

	//Should check if status code is ok
	result := rr.Result()
	if result.StatusCode != http.StatusOK {
		t.Errorf("ERROR: Expected %d, got %d", http.StatusAccepted, rr.Code)
	}

	if result.Header.Get("Content-type") != "application/json" {
		t.Errorf("Expected %s, got %s", "application/json", result.Header.Get("Content-type"))
	}

}

func Test404Get(t *testing.T) {
	testDb := make(mockDb)
	testDb[0] = nil
	eh := &handlers.EmployeeHandler{
		L:     log.Default(),
		DbMap: testDb,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/employees", nil)

	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	eh.GetEmployee(rr, req)

	result := rr.Result()

	if result.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected %d, got %d", http.StatusNotFound, result.StatusCode)
	}
}

func TestNullEmployee(t *testing.T) {
	testDb := make(mockDb)
	testDb[1] = nil
	eh := &handlers.EmployeeHandler{
		L:     log.Default(),
		DbMap: testDb,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/employees", nil)

	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	eh.GetEmployee(rr, req)

	result := rr.Result()

	if result.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected %d, got %d", http.StatusNotFound, result.StatusCode)
	}
}
