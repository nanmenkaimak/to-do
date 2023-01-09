package main

import (
	"crypto/tls"
	"flag"
	handlergor "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nanmenkaimak/to-do-list/internal/handlers"
	"log"
	"net/http"
	"os"
)

const portNumber = ":8080"

func main() {
	flag.Parse()

	router := mux.NewRouter()
	server := handlers.NewTaskServer()

	router.StrictSlash(true)
	// The "create task" path is protected with the BasicAuth middleware.
	//router.Handle("/task/", middleware.BasicAuth(http.HandlerFunc(server.CreateTask))).Methods("POST")
	router.HandleFunc("/task/", server.CreateTask).Methods("POST")
	router.HandleFunc("/task/", server.GetAllTasks).Methods("GET")
	router.HandleFunc("/task/", server.DeleteAllTasks).Methods("DELETE")
	router.HandleFunc("/task/{id:[0-9]+/", server.GetTask).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+/", server.DeleteTask).Methods("DELETE")
	router.HandleFunc("/task/{tag}/", server.GetTasksByTag).Methods("GET")
	router.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}/", server.GetTasksByDue).Methods("GET")

	// Set up logging and panic recovery middleware for all paths.
	router.Use(func(h http.Handler) http.Handler {
		return handlergor.LoggingHandler(os.Stdout, h)
	})
	router.Use(handlergor.RecoveryHandler(handlergor.PrintRecoveryStack(true)))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
