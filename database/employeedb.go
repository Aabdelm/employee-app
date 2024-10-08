package employeedb

import (
	"database/sql"
	"log"
	"os"

	"github.com/Aabdelm/employee-app/security"
	"github.com/go-sql-driver/mysql"
)

type EmployeeMapper interface {
	GetEmployeeById(id int) (*Employee, error)
	UpdateEmployee(*Employee) error
	DeleteEmployee(id int) error
	AddNewEmployee(*Employee) error
}
type Searcher interface {
	Search(identifier string, category string) ([]*Employee, error)
}

type DeptMapper interface {
	GetEmployeesByDepartment(id int) ([]*Employee, error)
	AddNewDepartment(dept *EmployeeDepartment) error
	RemoveDepartment(id int) error
	UpdateDepartment(id int, department *EmployeeDepartment, newName string) error
	GetAllDepartments() ([]*EmployeeDepartment, error)
	GetAllEmployees() ([]*Employee, error)
}
type DbMap struct {
	l  *log.Logger    // for logging
	M  map[string]int // for storing key-pair vals of departments
	Db *sql.DB        //the database
}

type Sec interface {
	FetchInfo(username string) (*security.Info, error)
}

type User struct {
	UserName string `json:"username"`
	UserPass string `json:"userpass"`
}

type Employee struct {
	Id           int    `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Department   string `json:"department"`
	DepartmentId int    `json:"departmentId"`
}

type EmployeeDepartment struct {
	Id         int    `json:"departmentId"`
	Department string `json:"department"`
}

func NewDbMap(l *log.Logger) *DbMap {
	return &DbMap{l: l, M: make(map[string]int, 0)}
}

func NewEmployee() *Employee {
	return &Employee{}
}

func NewEmployeeDepartment() *EmployeeDepartment {
	return &EmployeeDepartment{}
}

func NewUser() *User {
	return &User{}
}

/*
SetupDb deals with setting up the database. It returns a nil value if no errors were found
*/
func SetupDb(DbMap *DbMap, DbName string) error {
	cfg := &mysql.Config{
		User:   os.Getenv("DBUSERNAME"),
		Passwd: os.Getenv("DBPASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: DbName,
	}

	var err error
	DbMap.Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	err = DbMap.Db.Ping()

	if err == nil {
		DbMap.l.Println("[INFO] Connected Database")
	}
	return err
}
