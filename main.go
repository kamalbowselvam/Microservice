package main

import (
	"os/signal"
	"context"
	"time"
	"github.com/kamalselvam/Microservice/handlers"
	"log"
	"os"
	"net/http"
)

func main() {
	l := log.New(os.Stdout,"product-api",log.LstdFlags)

	ph := handlers.NewProduct(l)
	
	sm := http.NewServeMux()
	sm.Handle("/",ph)

	s := &http.Server{
		Addr: ":9090",
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1* time.Second,
	}

	go func(){
		l.Println("Server started on port :9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30* time.Second)
	defer cancel() 
	s.Shutdown(tc)
}
