package handlers_test

import (
	"fmt"

	employeedb "github.com/Aabdelm/employee-app/database"
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
