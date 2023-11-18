package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	employeedb "github.com/Aabdelm/employee-app/database"
	"github.com/Aabdelm/employee-app/handlers"
	"github.com/go-chi/chi/v5"
)

type MockDeptStruct struct {
	eDept     *employeedb.EmployeeDepartment
	employees []*employeedb.Employee
}
type MockDeptDb map[int]MockDeptStruct

// Methods to implement DeptMapper
func (md MockDeptDb) GetEmployeesByDepartment(id int) ([]*employeedb.Employee, error) {
	dept, exists := md[id]
	if !exists {
		return nil, fmt.Errorf("Failed to fetch employee %d", id)
	}
	return dept.employees, nil
}

func (md MockDeptDb) AddNewDepartment(dept *employeedb.EmployeeDepartment) error {
	_, exists := md[dept.Id]

	if exists {
		return fmt.Errorf("Error: department %d already exists", dept.Id)
	}

	md[dept.Id] = MockDeptStruct{
		eDept:     dept,
		employees: make([]*employeedb.Employee, 0),
	}

	return nil
}

func (md MockDeptDb) RemoveDepartment(deptId int) error {
	_, exists := md[deptId]
	if !exists {
		return fmt.Errorf("Error: department %d does not exist", deptId)
	}

	delete(md, deptId)
	return nil
}

func (md MockDeptDb) UpdateDepartment(id int, dept *employeedb.EmployeeDepartment, newName string) error {
	_, exists := md[id]
	if !exists {
		return fmt.Errorf("Error: deptartment %d does not exist.", id)
	}
	md[id].eDept.Department = newName
	return nil
}

// Unit tests
func TestAddNewDepartment(t *testing.T) {
	md := make(MockDeptDb, 0)
	testDept := &employeedb.EmployeeDepartment{
		Id:         1,
		Department: "Engineering",
	}
	j, _ := json.Marshal(testDept)
	bytes := bytes.NewReader(j)

	dh := handlers.DepartmentHandler{
		L:     log.Default(),
		DbMap: md,
	}
	req := httptest.NewRequest("POST", "/departments/", bytes)
	rr := httptest.NewRecorder()

	dh.PostDepartment(rr, req)

	body := rr.Result()
	defer body.Body.Close()

	if body.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusAccepted, body.StatusCode)

	}

}

// unimplemented methods
func TestUpdateDepartment(t *testing.T) {
	md := make(MockDeptDb, 0)
	testOldDept := &employeedb.EmployeeDepartment{
		Id:         1,
		Department: "Engineering",
	}
	md[1] = MockDeptStruct{
		eDept:     testOldDept,
		employees: make([]*employeedb.Employee, 0),
	}

	testNewDept := &employeedb.EmployeeDepartment{
		Id:         1,
		Department: "Finance",
	}

	j, _ := json.Marshal(testNewDept)
	bytes := bytes.NewReader(j)

	dh := handlers.DepartmentHandler{
		L:     log.Default(),
		DbMap: md,
	}
	req := httptest.NewRequest("PUT", "/departments/", bytes)
	rr := httptest.NewRecorder()

	dh.PostDepartment(rr, req)

	body := rr.Result()
	defer body.Body.Close()

	if body.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusAccepted, body.StatusCode)

	}
}

func TestGetEmployeesByDepartment(t *testing.T) {
	req := httptest.NewRequest("GET", "/employees/{id}", nil)
	rr := httptest.NewRecorder()

	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	md := make(MockDeptDb, 0)
	testOldDept := MockDeptStruct{
		eDept: &employeedb.EmployeeDepartment{
			Department: "Engineering",
			Id:         1,
		},
		employees: []*employeedb.Employee{
			{
				Id:           1,
				FirstName:    "First",
				LastName:     "Last",
				Email:        "email@yes.xyz",
				Department:   "Engineering",
				DepartmentId: 1,
			},
			{
				Id:           2,
				FirstName:    "Also",
				LastName:     "Another",
				Email:        "email@maybe.xyz",
				Department:   "Engineering",
				DepartmentId: 1,
			},
		},
	}
	md[1] = testOldDept

	dh := &handlers.DepartmentHandler{
		L:     log.Default(),
		DbMap: md,
	}

	dh.GetDepartment(rr, req)

	result := rr.Result()

	if result.StatusCode != http.StatusAccepted {
		t.Fatalf("Expected %d, got %d", http.StatusAccepted, result.StatusCode)
	}

}

func TestDeleteDepartments(t *testing.T) {
	req := httptest.NewRequest("GET", "/employees/{id}", nil)
	rr := httptest.NewRecorder()

	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx))

	md := make(MockDeptDb, 0)
	testOldDept := MockDeptStruct{
		eDept: &employeedb.EmployeeDepartment{
			Department: "Engineering",
			Id:         1,
		},
		employees: []*employeedb.Employee{
			{
				Id:           1,
				FirstName:    "First",
				LastName:     "Last",
				Email:        "email@yes.xyz",
				Department:   "Engineering",
				DepartmentId: 1,
			},
			{
				Id:           2,
				FirstName:    "Also",
				LastName:     "Another",
				Email:        "email@maybe.xyz",
				Department:   "Engineering",
				DepartmentId: 1,
			},
		},
	}
	md[1] = testOldDept

	dh := &handlers.DepartmentHandler{
		L:     log.Default(),
		DbMap: md,
	}

	dh.PutDepartment(rr, req)

	result := rr.Result()

	if result.StatusCode != http.StatusOK {
		t.Fatalf("Error: Expected %d, got %d", http.StatusOK, result.StatusCode)
	}

}
