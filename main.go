package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/url-shortener/internal/app"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func setupEnv() {
	err := godotenv.Load()
	if err != nil {
		logrus.Debug("Error loading .env files. ", err)
	}
}

func setupLogging() {
	logrus.SetFormatter(&logrus.TextFormatter{})
}

func main() {

	setupEnv()
	setupLogging()

	application, err := app.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize App. Error: %s", err.Error())
	}

	n := negroni.New()

	n.UseHandler(application.Router)

	var rootHandler http.Handler
	// set up CORS support
	allowedOrigin, useCors := os.LookupEnv("ORIGIN_ALLOWED")
	if useCors {
		// useful for local development by setting ORIGIN_ALLOWED="http://localhost:$APEX_PORT"
		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Content-Length", "Authorization"})
		originsOk := handlers.AllowedOrigins([]string{allowedOrigin})
		credentialsOk := handlers.AllowCredentials()
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST"})
		rootHandler = handlers.CORS(headersOk, originsOk, methodsOk, credentialsOk)(n)
	} else {
		rootHandler = n
	}

	port := os.Getenv("PORT")
	logrus.Infof("Accepting connections on :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), rootHandler))
}
