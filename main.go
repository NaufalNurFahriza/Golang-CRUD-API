package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"golang_MySQL/controllers"
	"golang_MySQL/middlewares"
	"golang_MySQL/models"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Database connection
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	db.AutoMigrate(&models.User{})

	// Initialize router
	router := mux.NewRouter()

	// Auth routes (not requiring token)
	router.HandleFunc("/api/register", controllers.Register(db)).Methods("POST")
	router.HandleFunc("/api/login", controllers.Login(db)).Methods("POST")

	// User CRUD routes (requiring token)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middlewares.JWTMiddleware)

	apiRouter.HandleFunc("/users", controllers.CreateUser(db)).Methods("POST")
	apiRouter.HandleFunc("/users", controllers.GetAllUsers(db)).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", controllers.GetUser(db)).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", controllers.UpdateUser(db)).Methods("PUT")
	apiRouter.HandleFunc("/users/{id}", controllers.DeleteUser(db)).Methods("DELETE")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
