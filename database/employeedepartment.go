package employeedb

import (
	"time"
)

/*
GetQuery gets the row(s) for employee(s) using the arguments provided
it returns a nil error if no error was found
*/
func (DbMap *DbMap) GetEmployeesByDepartment(deptId int) (employees []*Employee, err error) {
	//Query rows
	statement, err := DbMap.Db.Prepare(`
	SELECT id, first_name, last_name, email, employee_department.department_id
	FROM employee_department
	INNER JOIN employee ON id = employee.department_id
	WHERE employee_department.department_id = ?`)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute statement. error: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(deptId)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to Query for function GetEmployeesByDepartment. error: %s", err)
	}

	if err != nil {
		DbMap.l.Printf("error: ")
		return nil, err
	}

	//Make an employee slice (in case of an empty slice, we don't want to return nil)
	employees = make([]*Employee, 0)

	//iterate and append
	for rows.Next() {
		employee := NewEmployee()
		if err := rows.Scan(&employee.Id, &employee.FirstName, &employee.LastName, &employee.Email, &employee.DepartmentId); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return employees, nil
}

/*
AddNewDepartment adds a new department
The method returns an error that will help us later with error codes
*/
func (DbMap *DbMap) AddNewDepartment(dept *EmployeeDepartment) error {

	statement, err := DbMap.Db.Prepare(`INSERT INTO employee_department 
	(department_id, department, date_added, date_modified)
	VALUES (?,?,?,?)`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare query. Error %s", err)
		return err
	}
	defer statement.Close()

	row, err := statement.Exec(dept.Id, dept.Department, time.Now(), time.Now())

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute query for AddNewDepartment. Error:%s\n", err)
		return err
	}

	lastId, err := row.LastInsertId()
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to retrieve last inserted id. Error:%s\n", err)
		return err
	}

	DbMap.l.Printf("[INFO] Successfully added id %d\n", lastId)

	dept.Id = int(lastId)

	DbMap.M[dept.Department] = int(lastId)

	return nil

}

/*
RemoveDepartment removes a department
The method returns an error for HTTP error codes
*/
func (DbMap *DbMap) RemoveDepartment(deptId int) error {
	var err error

	statement, err := DbMap.Db.Prepare(`DELETE FROM employee_department
	WHERE department = ?)`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement for function RemoveDepartment. error:%s\n", err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(deptId)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute delete statement for function RemoveDepartment. error:%s\n", err)
		return err
	}

	return nil
}

func (DbMap *DbMap) UpdateDepartment(id int, department *EmployeeDepartment, newName string) error {
	var err error
	oldDepartment := department.Department

	statement, err := DbMap.Db.Prepare(`UPDATE employee_department
	SET department = ?, date_modified = ?
	WHERE department_id = ?`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement for function UpdateDepartment. error:%s\n", err)
	}

	defer statement.Close()

	_, err = statement.Exec(newName, time.Now(), id)

	if err != nil {
		DbMap.l.Printf("[ERROR] failed to delete, error:%s\n", err)
	}

	//Update the map
	delete(DbMap.M, oldDepartment)

	DbMap.M[newName] = id

	DbMap.l.Printf("[INFO] Successfully updated %d\n", id)

	return nil
}

func (DbMap *DbMap) GetAllEmployees() (employees []*Employee, err error) {
	employees = make([]*Employee, 0)

	stmt, err := DbMap.Db.Prepare(`SELECT 
	id, first_name, last_name, email, 
	department, employee_department.department_id 
	FROM employee_department
	INNER JOIN employee ON 
	employee_department.department_id = employee.department_id`)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement. Error %s", err)
		return employees, err
	}

	rows, err := stmt.Query()
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to query statement. Error %s", err)
		return nil, err
	}

	for rows.Next() {
		emp := NewEmployee()
		if err := rows.Scan(&emp.Id, &emp.FirstName, &emp.LastName, &emp.Email, &emp.Department, &emp.DepartmentId); err != nil {
			DbMap.l.Printf("[ERROR] Failed to scan row. Error %s", err)
			return employees, err
		}
		employees = append(employees, emp)
	}
	return employees, err
}

func (DbMap *DbMap) GetAllDepartments() ([]*EmployeeDepartment, error) {
	stmt, err := DbMap.Db.Prepare(`SELECT department_id, department FROM employee_department`)
	depts := make([]*EmployeeDepartment, 0)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement in function getAllDepartments, error %s", err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to query statement in function getAllDepartments, error %s", err)
		return nil, err
	}

	for rows.Next() {
		dept := NewEmployeeDepartment()
		rows.Scan(&dept.Id, &dept.Department)
		depts = append(depts, dept)

	}

	return depts, nil
}
