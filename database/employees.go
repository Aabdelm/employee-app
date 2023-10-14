package employeedb

import "time"

func (DbMap *DbMap) GetEmployeeById(id int) (*Employee, error) {
	row, err := DbMap.Db.Query(`
	SELECT id, first_name, last_name, email, department_id
	FROM employee_department
	INNER JOIN employee ON id = employee.department_id
	WHERE department_id = ?`, id)

	employee := NewEmployee()

	if err != nil {
		row.Close()
		DbMap.l.Printf("[ERROR] Query failed. error:%s\n", err)
		return nil, err
	}

	if err := row.Scan(employee.Id, employee.FirstName, employee.LastName, employee.Email, employee.Department); err != nil {
		DbMap.l.Printf("[ERROR] Qailed to scan row. error:%s\n", err)
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
	oldId := Employee.DepartmentId
	_, err := DbMap.Db.Exec(`UPDATE employee 
	SET department_id = ?`, departmentId)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to update id %d. Error: %s\n", departmentId, err)
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
	oldName := Employee.FirstName
	_, err := DbMap.Db.Exec(`UPDATE employee 
	SET first_name = ?`, newFirstName)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to update first name %s. Error: %s\n", oldName, err)
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
	oldName := Employee.LastName
	_, err := DbMap.Db.Exec(`UPDATE employee 
	SET last_name = ?, date_modified = ?`, newLastName, time.Now())
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
	oldEmail := Employee.Email
	_, err := DbMap.Db.Exec(`UPDATE employee 
	SET email = ?, date_modified = ?`, newEmail, time.Now())
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to update email %s. Error: %s\n", oldEmail, err)
		return err
	}

	Employee.Email = newEmail
	DbMap.l.Printf("[INFO] Updated Employee %d email %s to %s\n", Employee.Id, oldEmail, Employee.Email)
	return nil
}

//

//delete employee removes the employee off the database
func (DbMap *DbMap) DeleteEmployee(id int) error {
	_, err := DbMap.Db.Exec(`DELETE FROM employee
	WHERE id = ?`, id)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute query for id %d. Error: %s\n", id, err)
		return err
	}

	DbMap.l.Printf("[INFO] Successfully deleted %d\n", id)
	return nil
}
