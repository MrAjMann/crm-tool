package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/MrAjMann/crm/internal/model"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (repo *InvoiceRepository) GetAllInvoices() ([]model.Invoice, error) {
	rows, err := repo.db.Query("SELECT InvoiceId, InvoiceNumber, InvoiceDate, DueDate, CustomerId, CustomerName, CompanyName, CustomerPhone, CustomerEmail, PaymentStatus FROM invoices")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []model.Invoice
	for rows.Next() {
		var i model.Invoice
		if err := rows.Scan(&i.InvoiceId, &i.InvoiceNumber, &i.InvoiceDate, &i.DueDate, &i.CustomerId, &i.CustomerName, &i.CompanyName, &i.CustomerPhone, &i.CustomerEmail, &i.PaymentStatus); err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}
	return invoices, nil

}

func generateInvoiceId(lastId string) (string, error) {
	// Assuming lastId is in the format "INV0001"
	if lastId == "" {
		return "INV0001", nil
	}
	prefix := lastId[:3]                    // "INV"
	number, err := strconv.Atoi(lastId[3:]) // "0001" -> 1
	if err != nil {
		return "", err
	}
	newId := fmt.Sprintf("%s%04d", prefix, number+1)
	return newId, nil
}

func (repo *InvoiceRepository) AddNewInvoice(invoice model.Invoice) (string, error) {
	var lastInvoiceId string

	// Fetch the last InvoiceId to generate the next one
	err := repo.db.QueryRow("SELECT InvoiceId FROM invoices ORDER BY InvoiceId DESC LIMIT 1").Scan(&lastInvoiceId)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error fetching last invoice ID: %v", err)
	}

	// Generate new InvoiceId based on the last InvoiceId
	newInvoiceId, err := generateInvoiceId(lastInvoiceId)
	if err != nil {
		return "", err
	}

	// Prepare the new invoice with the generated InvoiceId and current time
	invoice.InvoiceId = newInvoiceId
	invoice.InvoiceDate = time.Now() // Set the invoice date to now

	// Perform the insert operation and return the newly created InvoiceId
	err = repo.db.QueryRow(
		"INSERT INTO invoices (InvoiceId, InvoiceNumber, InvoiceDate, DueDate, CustomerId, CustomerName, CompanyName, CustomerPhone, CustomerEmail, PaymentStatus) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING InvoiceId",
		newInvoiceId,          // $1
		invoice.InvoiceNumber, // $2 (Optionally generate this as well, if needed)
		invoice.InvoiceDate,   // $3
		invoice.DueDate,       // $4
		invoice.CustomerId,    // $5
		invoice.CustomerName,  // $6
		invoice.CompanyName,   // $7
		invoice.CustomerPhone, // $8
		invoice.CustomerEmail, // $9
		invoice.PaymentStatus, // $10
	).Scan(&newInvoiceId)
	if err != nil {
		return "", fmt.Errorf("error inserting new invoice: %v", err)
	}

	return newInvoiceId, nil
}

func GenerateInvoiceNumber(lastInvoiceNumber string) (string, error) {
	if lastInvoiceNumber == "" {
		return "INV0001", nil
	}

	numberPart := lastInvoiceNumber[3:]
	number, err := strconv.Atoi(numberPart)
	if err != nil {
		return "", fmt.Errorf("error converting invoice number to integer: %v", err)
	}

	if number >= 9999 {
		return "", fmt.Errorf("maximum invoice number reached")
	}

	newInvoiceNumber := fmt.Sprintf("INV%04d", number+1)
	return newInvoiceNumber, nil
}
