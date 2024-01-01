package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	employeedb "github.com/Aabdelm/employee-app/database"
	"github.com/Aabdelm/employee-app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// A collection for both employees and departments (for initial rendering)
type Collections struct {
	Employees   []*employeedb.Employee
	Departments []*employeedb.EmployeeDepartment
}

func newCollections(emps []*employeedb.Employee, depts []*employeedb.EmployeeDepartment) *Collections {
	return &Collections{
		Employees:   emps,
		Departments: depts,
	}
}

func main() {
	/*
		For the time being, os.stdout will be used
		A file will be used later
	*/
	l := log.New(os.Stdout, "", log.LstdFlags)

	router := chi.NewRouter()
	//Method router for usage in tandem with current logger
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"http://127.0.0.1:5500"},
			AllowedMethods: []string{"PUT", "GET", "POST", "DELETE"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		},
	))
	//Get (and serve) the static file directory
	static := http.FileServer(http.Dir("static"))

	//Handle static files
	router.Handle("/static/*", http.StripPrefix("/static/", static))

	//Create a signal channel to detect interrupts and/or shutdowns
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	//Create and setup the database
	db := employeedb.NewDbMap(l)
	if err := employeedb.SetupDb(db); err != nil {
		l.Printf("ERROR: FAILED TO CONNECT TO DATABASE. Error %s", err)
		time.Sleep(5 * time.Second)
		return
	}

	defer db.Db.Close()

	employeeHandler := handlers.NewEmployeeHandler(l, db)
	deptHandler := handlers.NewDepartmentHandler(l, db)

	router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/employees.html"))
		employees, err := db.GetAllEmployees()
		if err != nil {
			log.Fatal(err)
		}
		depts, err := db.GetAllDepartments()
		if err != nil {
			log.Fatal(err)
		}
		collection := newCollections(employees, depts)

		if err := tmpl.Execute(rw, collection); err != nil {
			log.Fatal(err)
		}

	})

	router.Get("/employees/{id}", employeeHandler.GetEmployee)
	router.Put("/employees/{id}", employeeHandler.PutEmployee)
	router.Post("/employees/", employeeHandler.PostEmployee)
	router.Delete("/employees/{id}", employeeHandler.DeleteEmployee)
	router.Delete("/employees/", employeeHandler.DeleteMultipleEmployees)
	router.Get("/employees/", deptHandler.GetAllEmployees)

	router.Post("/departments/", deptHandler.PostDepartment)
	router.Get("/departments/{id}", deptHandler.GetDepartment)
	router.Put("/departments/{id}", deptHandler.PutDepartment)
	router.Delete("/departments/{id}", deptHandler.DeleteDepartment)
	router.Get("/departments/", deptHandler.GetAllDepartments)

	s := &http.Server{
		Addr:    "localhost:80",
		Handler: router,
	}

	//run goroutine
	go func() {
		l.Println("Listening on port", s.Addr)
		s.ListenAndServe()

	}()

	//sig is blocking
	//go routine should run here until one of the two signals are recieved
	x := <-sig

	//shutdown
	l.Println("Signal Recieved:", x)
	l.Println("Attempting shutdown")

}
