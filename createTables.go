package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func CreateTables() {

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

	// Create leads table if it doesn't exist
	createLeadsTableSQL := `
		   DO $$
		   BEGIN
			   IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'leads') THEN
				   CREATE TABLE leads (
					   Id SERIAL PRIMARY KEY,
					   FirstName TEXT,
					   LastName TEXT,
					   Email TEXT,
					   CompanyName TEXT,
					   Phone TEXT,
					   Title TEXT,
					   Website TEXT,
					   Industry TEXT,
					   Source TEXT
				   );
			   END IF;
		   END
		   $$;`

	_, err = db.Exec(createLeadsTableSQL)
	if err != nil {
		log.Fatalf("Error creating invoices table: %v", err)
	}
	// Create customers table if it doesn't exist
	createCustomersTableSQL := `
		   DO $$
		   BEGIN
			   IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'customers') THEN
				   CREATE TABLE customers (
					   Id SERIAL PRIMARY KEY,
					   FirstName TEXT,
					   LastName TEXT,
					   Email TEXT,
					   CompanyName TEXT,
					   Phone TEXT,
					   Title TEXT,
					   Website TEXT,
					   Industry TEXT,
					   Source TEXT
				   );
			   END IF;
		   END
		   $$;`

	_, err = db.Exec(createCustomersTableSQL)
	if err != nil {
		log.Fatalf("Error creating invoices table: %v", err)
	}

	// Create customers table if it doesn't exist
	createInvoicesTableSQL := `
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'invoices') THEN
				CREATE TABLE invoices (
					InvoiceId SERIAL PRIMARY KEY,
					InvoiceNumber   TEXT 
					InvoiceDate     DATE
					DueDate         DATE
					CustomerId      TEXT
					CustomerName    TEXT
					CompanyName     TEXT
					CustomerPhone   TEXT
					CustomerEmail   TEXT
				);
			END IF;
		END
		$$;`

	_, err = db.Exec(createInvoicesTableSQL)
	if err != nil {
		log.Fatalf("Error creating invoices table: %v", err)
	}

}
