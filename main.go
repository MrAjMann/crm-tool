package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/MrAjMann/crm/internal/handler"
	"github.com/MrAjMann/crm/internal/repository"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	databaseURL := os.Getenv("DATABASE_URL")

	// Connect to the database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Templates
	// Parse templates
	sideBarTmpl, err := template.ParseGlob("src/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	sideBarTmpl, err = sideBarTmpl.ParseGlob("src/templates/navigation/*.html")
	if err != nil {
		log.Fatal(err)
	}

	customerRepo := repository.NewCustomerRepository(db)
	customerHandler := handler.NewCustomerHandler(customerRepo, sideBarTmpl)
	// Setup routes
	// Handlers

	// Setup routes
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/src/", http.StripPrefix("/src/", fs))

	http.HandleFunc("/", customerHandler.GetAllCustomers)          // Home page, possibly with a sidebar
	http.HandleFunc("/add-customer/", customerHandler.AddCustomer) // Handle adding a customer

	log.Fatal(http.ListenAndServe(":8080", nil))
}
