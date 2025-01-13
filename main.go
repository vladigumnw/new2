package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

var log = logrus.New()

func initLogger() {

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to log to file: %v\n", err)
		os.Exit(1)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.SetLevel(logrus.InfoLevel)

	log.SetFormatter(&logrus.JSONFormatter{})

	// log.Out = os.Stdout
	// log.SetLevel(logrus.InfoLevel)

	// level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	// if err == nil {
	// 	log.SetLevel(level)
	// } else if os.Getenv("LOG_LEVEL") != "" {
	// 	log.Warnf("Invalid LOG_LEVEL: %s, defaulting to InfoLevel", os.Getenv("LOG_LEVEL"))
	// }

	// log.SetFormatter(&logrus.JSONFormatter{})
}

func sendErrorEmail(errorMessage string) {
	m := gomail.NewMessage()

	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", "recipient@example.com")
	m.SetHeader("Subject", "Emergency: Application Error")

	m.SetBody("text/plain", fmt.Sprintf("An error occurred: \n\n%s", errorMessage))

	d := gomail.NewDialer("smtp.example.com", 587, "your-email@example.com", "your-email-password")

	if err := d.DialAndSend(m); err != nil {
		log.WithError(err).Error("Failed to send error email")
	}
}

func logErrorAndNotify(err error, context string) {
	log.WithFields(logrus.Fields{
		"error":   err.Error(),
		"context": context,
	}).Error("Critical error occurred")

	sendErrorEmail(fmt.Sprintf("Error: %s\nContext: %s", err.Error(), context))
}

type Response struct {
	Message string `json:"message"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(Response{Message: "API is healthy"})
}

func main() {

	initLogger()
	router := mux.NewRouter()

	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/hello", helloHandler).Methods("GET")
	router.HandleFunc("/status", statusHandler).Methods("GET")

	log.Info("Application starting...")
	log.Fatal(http.ListenAndServe(":8080", router))
	log.Println("Server is running on port 8080")

	log.Info("Connected to the database successfully")

	log.WithFields(logrus.Fields{
		"module": "main",
		"status": "initializing",
	}).Info("Initializing complete")

	log.WithFields(logrus.Fields{
		"endpoint": "/api/v1/resource",
		"method":   "GET",
		"status":   200,
	}).Info("Request handled successfully")

	router.HandleFunc("/api/health", HealthCheckHandler).Methods("GET")

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Home route accessed")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to the Home Page!")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Hello route accessed")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, User!")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Status route accessed")

	status := map[string]string{
		"status": "OK",
		"uptime": "72 Hours",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
