package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	protos "github.com/kamalbowselvam/Microservice/currency/protos/currency"
	gohandlers "github.com/gorilla/handlers"
	"github.com/go-openapi/runtime/middleware"
	"google.golang.org/grpc" 
	"github.com/gorilla/mux"
	"github.com/kamalbowselvam/Microservice/product-api/data"
	"github.com/kamalbowselvam/Microservice/product-api/handlers"
)

func main() {

	
	l := log.New(os.Stdout, "products-api ", log.LstdFlags)
	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:9092")
	
	if err != nil {
		panic(err)

	}

	defer conn.Close()

	cc := protos.NewCurrencyClient(conn)

	// create the handlers
	ph := handlers.NewProducts(l, v, cc)

	

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.ListAll)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products", ph.Update)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))



	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))


	s := &http.Server{
		Addr:         ":9090",
		Handler:      ch(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Server started on port :9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(tc)
	
}
