package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/kamalbowselvam/Microservice/product-images/files"
	"github.com/kamalbowselvam/Microservice/product-images/handlers"
	"github.com/nicholasjackson/env"
)




var bindAddress = env.String("BIND_ADDRESS", false, ":9091","Bind Address for the server")
var loglevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info and Trace]")
var basePath = env.String("BASE_PATH",false,"./imagestore","Base path to save images")

func main(){

	env.Parse()

	l := hclog.New(
			&hclog.LoggerOptions{
				Name: "product-images",
				Level: hclog.LevelFromString(*loglevel),
			},	
	)

	// create a logger for the server from the default logger 


	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the stograge class, using local storage 
	// max size 5 MB
	stor, err := files.NewLocal(*basePath, 1024*1000*5)
	
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}


	fh := handlers.NewFiles(stor, l)
	mw := handlers.GzipHandler{}

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()


	// upload files
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadREST)
	ph.HandleFunc("/", fh.UploadMultipart)

	// get files
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))),
	)
	gh.Use(mw.GzipMiddleware)

	s := http.Server{
		Addr: *bindAddress,
		Handler: sm,
		ErrorLog: sl,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}


	// start the server
	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}