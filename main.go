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
		log.Fatalf("Failed to initialize database tables: %v", err)
	}

	// Initialize database tables
	err = CreateTables(db)
	if err != nil {
		log.Fatalf("Failed to create database tables: %v", err)
	}

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Templates
	// Parse templates
	tmpl, err := template.ParseGlob("src/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err = tmpl.ParseGlob("src/templates/navigation/*.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.ParseGlob("src/templates/modals/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// customerRepo := repository.NewCustomerRepository(db)
	// if customerRepo == nil {
	// 	println("Creating customers table")
	// }

	leadRepo := repository.NewLeadRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)

	dashboardHandler := handler.NewDashboardHandler(tmpl)
	customerHandler := handler.NewCustomerHandler(customerRepo, tmpl)
	leadHandler := handler.NewLeadHandler(leadRepo, tmpl)
	invoiceHandler := handler.NewInvoiceHandler(invoiceRepo, tmpl, customerHandler)

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
	http.HandleFunc("/customers", customerHandler.GetAllCustomers)              // Customers page
	http.HandleFunc("/customer/", customerHandler.GetCustomerById)              // Handle getting a customer
	http.HandleFunc("/add-customer/", customerHandler.AddCustomer)              // Handle adding a customer
	http.HandleFunc("/search-customers", customerHandler.HandleSearchCustomers) // Handle searching for a customer
	http.HandleFunc("/customer/delete/", customerHandler.DeleteCustomer)        // Handle searching for a customer

	// Lead Routes
	http.HandleFunc("/leads", leadHandler.GetAllLeads) // Leads page
	http.HandleFunc("/lead/", leadHandler.GetLead)     // Handle getting a lead
	http.HandleFunc("/add-lead/", leadHandler.AddLead) // Handle adding a lead

	//Invoice Routes
	http.HandleFunc("/invoices", invoiceHandler.GetAllInvoices)
	http.HandleFunc("/add-invoice", invoiceHandler.AddNewInvoice)
	http.HandleFunc("/calculate-invoice", invoiceHandler.InvoiceCalculationHandler)

	// PDF GEN
	http.HandleFunc("/generate-pdf", invoiceHandler.GeneratePdf)

	// Session Routes
	http.HandleFunc("/customer-session", customerHandler.HandleSessionStore)

	http.HandleFunc("/create-lead-modal", func(w http.ResponseWriter, r *http.Request) {
		modalPath := "src/templates/modals/createLeadModal.html"
		http.ServeFile(w, r, modalPath)
	})

	http.HandleFunc("/create-customer-modal", func(w http.ResponseWriter, r *http.Request) {
		modalPath := "src/templates/modals/createCustomerModal.html"
		http.ServeFile(w, r, modalPath)
	})

	http.HandleFunc("/create-invoice", func(w http.ResponseWriter, r *http.Request) {
		// Execute the template that includes the sidebar
		err := tmpl.ExecuteTemplate(w, "createInvoice.html", nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			log.Println(err)
		}
	})

	// Server startup and error handling remain unchanged
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
