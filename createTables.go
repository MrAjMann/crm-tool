package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// CreateTables initializes the database tables if they do not exist
func CreateTables(db *sql.DB) error {
	if err := loadEnvVariables(); err != nil {
		return err
	}

	if err := createStatusTable(db); err != nil {
		return err
	}

	if err := populateStatusTable(db); err != nil {
		return err
	}

	if err := createOtherTables(db); err != nil {
		return err
	}

	return nil
}

// loadEnvVariables loads the environment variables from the .env file.
func loadEnvVariables() error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file", err)
		return err
	}
	return nil
}

// createStatusTable creates the 'status' table if it doesn't exist.
func createStatusTable(db *sql.DB) error {
	sql := `
    CREATE TABLE IF NOT EXISTS status (
        StatusId SERIAL PRIMARY KEY,
        StatusValue TEXT NOT NULL UNIQUE,
        IsClosed BOOLEAN NOT NULL DEFAULT FALSE,
        ClosedStatusValue TEXT
    );`
	if _, err := db.Exec(sql); err != nil {
		log.Fatalf("Error creating status table: %v", err)
		return err
	}
	return nil
}

// populateStatusTable populates the 'status' table with predefined values.
func populateStatusTable(db *sql.DB) error {
	sql := `
INSERT INTO status (StatusValue, IsClosed, ClosedStatusValue) VALUES
('New Lead', FALSE, NULL),
('Contacted', FALSE, NULL),
('Engaged', FALSE, NULL),
('Qualified', FALSE, NULL),
('Needs Analysis', FALSE, NULL),
('Proposal Sent', FALSE, NULL),
('Negotiation', FALSE, NULL),
('Closed', TRUE, 'Still Fighting'),
('Closed', TRUE, 'Won'),
('Closed', TRUE, 'Lost')
ON CONFLICT (StatusValue) DO NOTHING;`
	if _, err := db.Exec(sql); err != nil {
		log.Fatalf("Error populating status table: %v", err)
		return err
	}
	return nil
}

// createOtherTables creates all other necessary tables.
func createOtherTables(db *sql.DB) error {
	tableCreationSQLs := []string{
		`DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'leads') THEN
                CREATE TABLE leads (
                    Id SERIAL PRIMARY KEY,
                    FirstName TEXT,
                    LastName TEXT,
                    CompanyName TEXT,
                    Email TEXT,
                    Phone TEXT,
                    StatusId INTEGER NOT NULL,
                    Title TEXT,
                    Website TEXT,
                    Industry TEXT,
                    ServiceType TEXT,
                    Source TEXT,
                    CreatedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                    UpdatedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                    FOREIGN KEY (StatusId) REFERENCES status(StatusId)
                );
            END IF;
        END
        $$;`,
		`DO $$
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
        $$;`,
		`DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'invoices') THEN
                CREATE TABLE invoices (
                    InvoiceId SERIAL PRIMARY KEY,
                    InvoiceNumber TEXT,
                    InvoiceDate TIMESTAMP WITHOUT TIME ZONE NOT NULL,
                    DueDate DATE,
                    CustomerId INTEGER NOT NULL,
                    CustomerName TEXT NOT NULL,
                    CompanyName TEXT,
                    CustomerPhone TEXT NOT NULL,
                    CustomerEmail TEXT NOT NULL,
                    PaymentStatus INTEGER NOT NULL,
                    CustomerAddress TEXT NOT NULL,
                    CreatedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                    UpdatedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                    FOREIGN KEY (CustomerId) REFERENCES customers(Id) ON DELETE CASCADE
                );
            END IF;
        END
        $$;`,
		`DO $$
        BEGIN
        IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'notes') THEN
            CREATE TABLE notes (
                NoteId SERIAL PRIMARY KEY,
                CustomerId INTEGER,
                LeadId INTEGER,
                Category TEXT NOT NULL,
                AuthorId INTEGER,
                AuthorName TEXT,
                Content TEXT,
                CreatedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                UpdatedAt TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                FOREIGN KEY (CustomerId) REFERENCES customers(Id) ON DELETE SET NULL,
                FOREIGN KEY (LeadId) REFERENCES leads(Id) ON DELETE SET NULL
            );
        END IF;
    END
    $$;`,
		`DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'item_lists') THEN
                CREATE TABLE item_lists (
                    ItemId SERIAL PRIMARY KEY,
                    InvoiceId INTEGER NOT NULL,
                    Item TEXT NOT NULL,
                    Quantity INTEGER NOT NULL,
                    UnitPrice DECIMAL NOT NULL,
                    Subtotal DECIMAL NOT NULL,
                    Tax DECIMAL NOT NULL,
                    Total DECIMAL NOT NULL,
                    FOREIGN KEY (InvoiceId) REFERENCES invoices(InvoiceId) ON DELETE CASCADE
                );
            END IF;
        END
        $$;`,
		`DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'address') THEN
                CREATE TABLE address (
                    UnitNumber TEXT,
                    StreetNumber TEXT,
                    StreetName TEXT,
                    City TEXT,
                    Postcode TEXT,
                    PRIMARY KEY (StreetNumber, StreetName, City, Postcode)
                );
            END IF;
        END
        $$;`,
		`DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'service_entry') THEN
                CREATE TABLE service_entry (
                    EntryId SERIAL PRIMARY KEY,
                    ServiceType TEXT,
                    StartDate DATE,
                    DueDate DATE,
                    EndDate DATE
                );
            END IF;
        END
        $$;`,
	}
	for _, sql := range tableCreationSQLs {
		if _, err := db.Exec(sql); err != nil {
			log.Fatalf("Error creating table: %v", err)
			return err
		}
	}
	return nil
}
