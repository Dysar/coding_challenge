package main

import (
	"challenge/internal/config"
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
	logrus.Infof("Go version: %s (%s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	conf, confErr := config.New("conf.json")
	if confErr != nil {
		logrus.Panic("config parser", confErr)
	}

	router := server.NewRouter(conf)

	srv := &http.Server{
		ReadHeaderTimeout: conf.Server.ReadHeaderTimeoutSeconds * time.Second,
		ReadTimeout:       conf.Server.ReadTimeoutSeconds * time.Second,
		WriteTimeout:      conf.Server.WriteTimeoutSeconds * time.Second,
		Handler:           router,
		Addr:              fmt.Sprintf(":%d", conf.Port),
	}

	logrus.Infof("Server is listening on %d", conf.Port)
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
