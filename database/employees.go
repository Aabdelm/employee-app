package employeedb

import (
	"time"
)

func (DbMap *DbMap) GetEmployeeById(id int) (*Employee, error) {
	var err error
	statement, err := DbMap.Db.Prepare(`
	SELECT id, first_name, last_name, email, department,
	employee_department.department_id
	FROM employee_department
	INNER JOIN employee ON employee_department.department_id = employee.department_id
	WHERE id= ?`)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement on function GetEmployeeById. error: %s", err)
		return nil, err
	}

	defer statement.Close()
	row := statement.QueryRow(id)

	if err != nil {
		DbMap.l.Printf("[ERROR] Query failed. error:%s\n", err)
		return nil, err
	}

	employee := NewEmployee()
	err = row.Scan(&employee.Id, &employee.FirstName, &employee.LastName, &employee.Email, &employee.Department, &employee.DepartmentId)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to scan row. error:%s\n", err)
		return nil, err
	}

	return employee, nil

}

//The update methods here are really the same with little exceptions

/*
UpdateEmployeeDepartment updates the employee department
It utilizes an employee pointer for later json parsing
*/
func (DbMap *DbMap) UpdateEmployeeDepartment(Employee *Employee, departmentId int) error {
	var err error

	oldId := Employee.DepartmentId
	statement, err := DbMap.Db.Prepare(`UPDATE employee 
	SET department_id = ?`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement on function UpdateEmployeeDepartment %d. Error: %s\n", departmentId, err)
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(departmentId)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute statement on function UpdateEmployeeDepartment. Error: %s\n", err)
		return err
	}

	Employee.Id = departmentId
	DbMap.l.Printf("[INFO] Updated Employee %d from department %d to %d\n", Employee.Id, oldId, Employee.DepartmentId)
	return nil
}

/*
UpdateEmployeeDepartment updates the employee's first name
It utilizes an employee pointer for later json parsing
*/
func (DbMap *DbMap) UpdateEmployeeFirstName(Employee *Employee, newFirstName string, departmentId int) error {
	var err error

	oldName := Employee.FirstName

	statement, err := DbMap.Db.Prepare(`UPDATE employee 
	SET first_name = ?`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement on function UpdateEmployeeFirstName. Error: %s\n", err)
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(newFirstName)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute statement on function UpdateEmployeeFirstName. Error %s\n", err)
		return err
	}

	Employee.FirstName = newFirstName
	DbMap.l.Printf("[INFO] Updated Employee %d first name %s to %s\n", Employee.Id, oldName, Employee.LastName)
	return nil
}

/*
UpdateEmployeeDepartment updates the employee's last name
It utilizes an employee pointer for later json parsing
*/
func (DbMap *DbMap) UpdateEmployeeLastName(Employee *Employee, newLastName string, departmentId int) error {
	var err error

	oldName := Employee.LastName
	statement, err := DbMap.Db.Prepare(`UPDATE employee 
	SET last_name = ?, date_modified = ?`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement on function UpdateEmployeeLastName. Error: %s\n", err)
	}

	defer statement.Close()

	_, err = statement.Exec(newLastName, time.Now())

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to update last name %s. Error: %s\n", oldName, err)
		return err
	}

	Employee.LastName = newLastName
	DbMap.l.Printf("[INFO] Updated Employee %d last name %s to %s\n", Employee.Id, oldName, Employee.LastName)
	return nil
}

/*
UpdateEmployeeDepartment updates the employee's email
It utilizes an employee pointer for later json parsing
*/
func (DbMap *DbMap) UpdateEmployeeEmail(Employee *Employee, newEmail string, departmentId int) error {
	var err error

	oldEmail := Employee.Email
	statement, err := DbMap.Db.Prepare(`UPDATE employee 
	SET email = ?, date_modified = ?`)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement on function UpdateEmployeeEmail. Error: %s\n", err)
	}

	_, err = statement.Exec(newEmail, time.Now())
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to update email %s. Error: %s\n", oldEmail, err)
		return err
	}

	Employee.Email = newEmail
	DbMap.l.Printf("[INFO] Updated Employee %d email %s to %s\n", Employee.Id, oldEmail, Employee.Email)
	return nil
}

//

// delete employee removes the employee off the database
func (DbMap *DbMap) DeleteEmployee(id int) error {
	statement, err := DbMap.Db.Prepare(`DELETE FROM employee
	WHERE id = ?`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute query for id %d. Error: %s\n", id, err)
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute statement on function DeleteEmployee. Error %s\n", err)
		return err
	}

	DbMap.l.Printf("[INFO] Successfully deleted %d\n", id)
	return nil
}

func (DbMap *DbMap) AddNewEmployee(employee *Employee) error {
	var err error

	statement, err := DbMap.Db.Prepare(`INSERT INTO employee
	 (id,first_name, last_name, email, department_id, date_added, date_modified) 
	 VALUES (?,?,?,?,?,?,?)`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement on funcion AddNewEmployee. Error %s\n", err)
		return err

	}
	_, err = statement.Exec(employee.Id, employee.FirstName, employee.LastName, employee.Email,
		employee.DepartmentId, time.Now(), time.Now())

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute statement. Error %s\n", err)
		return err
	}

	DbMap.l.Printf("[INFO] Added employee to database \n")

	return nil
}
