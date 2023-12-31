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
		DbMap.l.Printf("[ERROR] Query failed. In function GetEmployeeById. error:%s\n", err)
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

func (DbMap *DbMap) UpdateEmployee(employee *Employee) error {
	var err error
	DbMap.l.Println("Function UpdateEmployee called")

	// if err = DbMap.disableForeignKeyChecks(); err != nil {
	// 	return err
	// }

	stmt, err := DbMap.Db.Prepare(`
	UPDATE employee_department
	INNER JOIN employee ON employee_department.department_id = employee.department_id
	SET first_name = ?, last_name = ?, email = ?, employee.department_id = ?, employee.date_modified = ?
	WHERE id = ? AND employee_department.department_id = ?
	`)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement in function UpdateEmployee. Error %s\n", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(employee.FirstName, employee.LastName, employee.Email,
		employee.DepartmentId, time.Now(), employee.Id, employee.DepartmentId)

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute statement in function UpdateEmployee. Error %s\n", err)
		return err
	}

	if err = DbMap.enableForeignKeyChecks(); err != nil {
		DbMap.l.Printf("[ERROR] Failed to enable foreign key checks. Error %d \n", err)
		return err
	}

	DbMap.l.Printf("[INFO] Successfully updated employee %d\n", employee.Id)

	return nil
}

// delete employee removes the employee off the database
func (DbMap *DbMap) DeleteEmployee(id int) error {
	DbMap.l.Println("[INFO] Function Delete Employee called")
	statement, err := DbMap.Db.Prepare(`DELETE FROM employee
	WHERE id = ?`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute query for id %d in Function DeleteEmployee. Error: %s\n", id, err)
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
	DbMap.l.Println("[INFO] Function AddNewEmployee called")

	statement, err := DbMap.Db.Prepare(`INSERT INTO employee
	 (first_name, last_name, email, department_id, date_added, date_modified) 
	 VALUES (?,?,?,?,?,?)`)
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to prepare statement on funcion AddNewEmployee. Error %s\n", err)
		return err

	}
	res, err := statement.Exec(employee.FirstName, employee.LastName, employee.Email,
		employee.DepartmentId, time.Now(), time.Now())

	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute statement. Error %s\n", err)
		return err
	}
	//We'll need a new id for the json request
	DbMap.l.Printf("Assigning new id in function AddNewEmployee\n")
	id, err := res.LastInsertId()
	if err != nil {
		DbMap.l.Printf("Failed to parse last inserted Id. Error %s", err)
		return err
	}

	employee.Id = int(id)

	DbMap.l.Printf("[INFO] Added employee to database \n")
	return nil
}
