package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Rennerp/microservices_tutorial/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Create the handlers
	productHandler := handlers.NewProducts(l)

	// Create A new ServerMux and register handlers
	serveMux := http.NewServeMux()
	serveMux.Handle("/", productHandler)

	// Server Configuration
	server := &http.Server{
		Addr:         ":9090",           // Configure the bind address
		Handler:      serveMux,          // set the default handler
		ErrorLog:     l,                 // Set the default logger for the server
		IdleTimeout:  120 * time.Second, // Max time for connections
		ReadTimeout:  1 * time.Second,   // Max time for read request from client
		WriteTimeout: 1 * time.Second,   // Max time for write request to the client
	}

	// Start Server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Trap sigterm or interrupt and gracefully shutdown the server
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt)
	signal.Notify(sigChannel, os.Kill)

	sig := <-sigChannel
	l.Println("Recieved termiante, graceful shutdown", sig)

	timeOutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeOutContext)
}
