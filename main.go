package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	/*
		For the time being, os.stdout will be used
		A file will be used later
	*/
	l := log.New(os.Stdout, "employee-app:", log.LstdFlags)
	s := &http.Server{
		Addr: "localhost:80",
	}
	//Create a signal channel to detect interrupts and/or shutdowns
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	//run goroutine
	go func() {
		l.Println("Listening on port", s.Addr)
		s.ListenAndServe()

	}()
	x := <-sig
	l.Println("Signal Recieved:", x)
	l.Println("Attempting shutdown")

}
