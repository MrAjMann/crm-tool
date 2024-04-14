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

func (repo *InvoiceRepository) AddNewInvoice(invoice model.Invoice) (string, error) {
	var invoiceId string
	var lastInvoiceNumber string

	err := repo.db.QueryRow("SELECT InvoiceNumber FROM invoices ORDER by InvoiceNumber DESC LIMIT 1").Scan(&lastInvoiceNumber)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error fetching last invoice number: %v", err)
	}

	newInvoiceNumber, err := GenerateInvoiceNumber(lastInvoiceNumber)
	if err != nil {
		return "", err
	}

	invoiceDate := time.Now()

	

	// The query must include actual parameters from the 'invoice' object
	err = repo.db.QueryRow(
		"INSERT INTO invoices ( InvoiceNumber, InvoiceDate, DueDate, CustomerId, CustomerName, CompanyName, CustomerPhone, CustomerEmail, PaymentStatus) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING InvoiceId",
		newInvoiceNumber,      // $1
		invoiceDate,           // $2
		invoice.DueDate,       // $3
		invoice.CustomerId,    // $4 Assume there's a CustomerId field in your model
		invoice.CustomerName,  // $5
		invoice.CompanyName,   // $6
		invoice.CustomerPhone, // $7
		invoice.CustomerEmail, // $8
		invoice.PaymentStatus, // $9
	).Scan(&invoiceId)

	if err != nil {
		return "", fmt.Errorf("error returning InvoiceId: %v", err)
	}
	return invoiceId, nil
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
