package employeedb

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DbMap struct {
	l  *log.Logger    // for logging
	M  map[int]string // for storing key-pair vals of departments
	Db *sql.DB        //the database
}

type Employee struct {
	Id         int    `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Department *EmployeeDepartment
}

type EmployeeDepartment struct {
	Id         int    `json:"employeeId"`
	Department string `json:"employeeDepartment"`
}

func NewDbMap(l *log.Logger) *DbMap {
	return &DbMap{l: l, M: make(map[int]string)}
}

func NewEmployee() *Employee {
	return &Employee{}
}

/*
SetupDb deals with setting up the database. It returns a nil value if no errors were found
*/
func SetupDb(DbMap *DbMap) error {
	cfg := &mysql.Config{
		User:   os.Getenv("DBUSERNAME"),
		Passwd: os.Getenv("DBPASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "employee_management",
	}

	var err error = nil
	DbMap.Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	err = DbMap.Db.Ping()

	DbMap.l.Println("Connected Database")
	return err
}

/*
GetQuery gets the row(s) for employee(s) using the arguments provided
it returns a nil error if no error was found
*/
func (DbMap *DbMap) GetEmployeesByDepartment(department string) (employees []*Employee, err error) {
	//Query rows
	rows, err := DbMap.Db.Query(`
	SELECT id, first_name, last_name, email, department_id
	FROM employee_department
	WHERE department = ?
	INNER JOIN employee ON id = employee.department_id`, department)

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
		if err := rows.Scan(employee.Id, employee.FirstName, employee.LastName, employee.Email, employee.Department.Id); err != nil {
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
*/
func (DbMap *DbMap) AddNewDepartment(department string) {
	id, err := DbMap.Db.Exec(`INSERT INTO employee_department (department, date_added, date_modified)
	VALUES (?,?,?)`, department, time.Now(), time.Now())
	if err != nil {
		DbMap.l.Fatalf("ERROR: Failed to execute query. Error:%s", err)
	}
	DbMap.l.Printf("Successfully added id %d\n", id)

}

/*
 */
func (DbMap *DbMap) RemoveDepartment(department string) {
	id, err := DbMap.Db.Exec(`DELETE FROM employee_department
	WHERE department = ?)`, department)
	if err != nil {
		DbMap.l.Fatalf("FATAL: failed to delete, error:%s", err)
	}
	i, err := id.LastInsertId()
	if err != nil {
		DbMap.l.Fatalf("FATAL: failed to retrieve id %d, error:%s", i, err)
	}

	DbMap.l.Printf("Database: Successfully removed %d", i)
}
