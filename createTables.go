package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// CreateTables initializes the database tables if they do not exist
func CreateTables(db *sql.DB) error {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// SQL statement to create 'leads' table if it doesn't exist
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
	if _, err = db.Exec(createLeadsTableSQL); err != nil {
		log.Fatalf("Error creating leads table: %v", err)
		return err
	}

	// SQL statement to create 'customers' table if it doesn't exist
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
	if _, err = db.Exec(createCustomersTableSQL); err != nil {
		log.Fatalf("Error creating customers table: %v", err)
		return err
	}

	// SQL statement to create 'invoices' table if it doesn't exist
	createInvoicesTableSQL := `
        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'invoices') THEN
                CREATE TABLE invoices (
                    InvoiceId SERIAL PRIMARY KEY,
                    InvoiceNumber TEXT,
                    InvoiceDate DATE,
                    DueDate DATE,
                    CustomerId TEXT,
                    CustomerName TEXT,
                    CompanyName TEXT,
                    CustomerPhone TEXT,
                    CustomerEmail TEXT,
                    PaymentStatus INTEGER
                );
            END IF;
        END
        $$;`

	if _, err = db.Exec(createInvoicesTableSQL); err != nil {
		log.Fatalf("Error creating invoices table: %v", err)
		return err
	}

	// SQL statement to create 'itemlist' table if it doesn't exist
	createItemListTableSQL := `
        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'itemlist') THEN
                CREATE TABLE itemlist (
                    InvoiceId TEXT,
                    Item      TEXT,
                    Quantity  INTEGER,
                    UnitPrice INTEGER,
                    Subtotal  INTEGER,
                    Tax       INTEGER,
                    Total     INTEGER
                );
            END IF;
        END
        $$;`

	if _, err = db.Exec(createItemListTableSQL); err != nil {
		log.Fatalf("Error creating itemlist table: %v", err)
		return err
	}
	return nil // Return nil if all tables are created successfully
}
