package main

import (
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
)

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

	router.Get("/employees/{id}", employeeHandler.GetEmployee)
	router.Put("/employees/{id}", employeeHandler.PutEmployee)
	router.Post("/employees/", employeeHandler.PostEmployee)
	router.Delete("/employees/{id}", employeeHandler.DeleteEmployee)

	router.Post("/departments/", deptHandler.PostDepartment)
	router.Get("/departments/{id}", deptHandler.GetDepartment)
	router.Put("/departments/{id}", deptHandler.PutDepartment)
	router.Delete("/departments/{id}", deptHandler.DeleteDepartment)

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
