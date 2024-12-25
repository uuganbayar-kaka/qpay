package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"qpay/qpay"

	"github.com/gorilla/handlers"
)

// GlVersion :
var GlVersion string

// GlBuildDate :
var GlBuildDate string

// GlRunMode :
var GlRunMode string

// func main() {
// 	router := mux.NewRouter()

// 	qpay.LinkHandlersV1(router)

// 	fmt.Println("Server is running on port 4000...")
// 	log.Fatal(http.ListenAndServe(":4000", nil))
// }

func main() {
	router := qpay.LinkHandlersV1() // Set up the application's configuration
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, QPay Middleware!"))
	})

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// Create an HTTP server
	httpServer := &http.Server{
		Addr:    "0.0.0.0:4000",
		Handler: loggedRouter,
	}

	// Start the server
	log.Println("Server is running on port 4000")
	listenErr := httpServer.ListenAndServe() // Listen and serve incoming requests
	if listenErr != nil {
		fmt.Printf("REST SERVER LISTEN ERROR [%s]\n", listenErr.Error()) // Log the error
		os.Exit(2)                                                       // Exit the program with error status
	}
}
