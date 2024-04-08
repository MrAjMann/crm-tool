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

	log.Println("Connecting to database...")
	// Database connection code
	if err != nil {
		log.Printf("Database connection error: %v", err)
		return
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
	sideBarTmpl, err = sideBarTmpl.ParseGlob("src/templates/modals/*.html")
	if err != nil {
		log.Fatal(err)
	}

	customerRepo := repository.NewCustomerRepository(db)
	if customerRepo == nil {
		println("Creating customers table")
	}

	leadRepo := repository.NewLeadRepository(db)

	dashboardHandler := handler.NewDashboardHandler(sideBarTmpl)
	customerHandler := handler.NewCustomerHandler(customerRepo, sideBarTmpl)
	leadHandler := handler.NewLeadHandler(leadRepo, sideBarTmpl)

	// Setup routes
	// Handlers

	// Setup routes
	fs := http.FileServer(http.Dir("src"))
	css := http.FileServer(http.Dir("css"))

	http.Handle("/src/", http.StripPrefix("/src/", fs))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	// Dashboard Routes
	http.HandleFunc("/", dashboardHandler.Dashboard)

	// Customer Routes
	http.HandleFunc("/customers", customerHandler.GetAllCustomers) // Customers page
	http.HandleFunc("/customer/", customerHandler.GetCustomer)     // Handle getting a customer
	http.HandleFunc("/add-customer/", customerHandler.AddCustomer) // Handle adding a customer

	// Lead Routes
	http.HandleFunc("/leads", leadHandler.GetAllLeads) // Leads page
	http.HandleFunc("/lead/", leadHandler.GetLead)     // Handle getting a lead
	http.HandleFunc("/add-lead/", leadHandler.AddLead) // Handle adding a lead
	http.HandleFunc("src/templates/modals/create-lead-modal", func(w http.ResponseWriter, r *http.Request) {
		err := sideBarTmpl.ExecuteTemplate(w, "createLeadModal.html", nil)
		if err != nil {
			http.Error(w, "Error loading modal", http.StatusInternalServerError)
			log.Printf("Error loading modal: %v\n", err)
		}
		log.Printf("Modal loaded")

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
