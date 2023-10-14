package employeedb

import (
	"time"
)

/*
GetQuery gets the row(s) for employee(s) using the arguments provided
it returns a nil error if no error was found
*/
func (DbMap *DbMap) GetEmployeesByDepartment(department string) (employees []*Employee, err error) {
	//Query rows
	rows, err := DbMap.Db.Query(`
	SELECT id, first_name, last_name, email, department_id
	FROM employee_department
	INNER JOIN employee ON id = employee.department_id
	WHERE department = ?`, department)

	if err != nil {
		rows.Close()
		DbMap.l.Printf("error: ")
		return nil, err
	}

	//Make an employee slice (in case of an empty slice, we don't want to return nil)
	employees = make([]*Employee, 0)

	//iterate and append
	for rows.Next() {
		employee := NewEmployee()
		if err := rows.Scan(employee.Id, employee.FirstName, employee.LastName, employee.Email, employee.DepartmentId); err != nil {
			rows.Close()
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
func (DbMap *DbMap) AddNewDepartment(department string) error {
	id, err := DbMap.Db.Exec(`INSERT INTO employee_department (department, date_added, date_modified)
	VALUES (?,?,?)`, department, time.Now(), time.Now())
	if err != nil {
		DbMap.l.Printf("[ERROR] Failed to execute query. Error:%s\n", err)
	}
	DbMap.l.Printf("[INFO] Successfully added id %d\n", id)

	return nil

}

/*
RemoveDepartment removes a department
The method returns an error for HTTP error codes
*/
func (DbMap *DbMap) RemoveDepartment(department string) error {
	_, err := DbMap.Db.Exec(`DELETE FROM employee_department
	WHERE department = ?)`, department)
	if err != nil {
		DbMap.l.Printf("[ERROR] failed to delete, error:%s\n", err)
	}
	return nil
}

func (DbMap *DbMap) UpdateDepartment(id int, newName string) error {
	_, err := DbMap.Db.Exec(`UPDATE employee_department
	SET department = ?, date_modified = ?
	WHERE department_id = ?`, newName, time.Now(), id)
	if err != nil {
		DbMap.l.Printf("[ERROR] failed to delete, error:%s\n", err)
	}

	DbMap.l.Printf("[INFO] Successfully updated %d\n", id)

	return nil
}
