package main

import (
	"challenge/internal/server"
	"fmt"
	"github.com/braintree/manners"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

// @contact.email

func main() {
	fmt.Printf("Go version: %s (%s/%s)\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	router := server.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	srv := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       600 * time.Second,
		WriteTimeout:      600 * time.Second,
		Handler:           router,
		Addr:              fmt.Sprintf(":%s", port),
	}

	fmt.Printf("Server is listening on %s\n", port)
	if err := manners.NewWithServer(srv).ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go func(ch <-chan os.Signal) {
		<-ch
		manners.Close()
	}(ch)
}
