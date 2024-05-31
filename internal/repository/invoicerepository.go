package repository

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"github.com/MrAjMann/crm/internal/model"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}
func (r *InvoiceRepository) BeginTransaction() (*sql.Tx, error) {
	slog.Info("r Beginning transaction")
	return r.db.Begin()
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
	slog.Info("r Generating Invoice ID")
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

func (r *InvoiceRepository) AddNewInvoice(tx *sql.Tx, invoice model.Invoice) (string, error) {
	slog.Info("r Adding the Invoice")

	// Generate invoice ID
	invoiceId, err := generateInvoiceId(invoice.InvoiceId)
	if err != nil {
		return "", err
	}

	// PostgreSQL query with numbered placeholders
	query := `
        INSERT INTO invoices (InvoiceId, InvoiceNumber, InvoiceDate, DueDate, CustomerId, CustomerName, CompanyName, CustomerPhone, CustomerEmail, PaymentStatus)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	log.Printf("Executing query: %s\n", query)
	log.Printf("With parameters: InvoiceId=%s, InvoiceNumber=%s, InvoiceDate=%s, DueDate=%s, CustomerId=%d, CustomerName=%s, CompanyName=%s, CustomerPhone=%s, CustomerEmail=%s, PaymentStatus=%d",
		invoiceId, invoice.InvoiceNumber, invoice.InvoiceDate, invoice.DueDate, invoice.CustomerId, invoice.CustomerName, invoice.CompanyName, invoice.CustomerPhone, invoice.CustomerEmail, invoice.PaymentStatus)

	// Execute the query
	_, err = tx.Exec(
		query,
		invoiceId, invoice.InvoiceNumber, invoice.InvoiceDate, invoice.DueDate, invoice.CustomerId, invoice.CustomerName, invoice.CompanyName, invoice.CustomerPhone, invoice.CustomerEmail, invoice.PaymentStatus,
	)
	if err != nil {
		return "", err
	}

	return invoiceId, nil
}

func (r *InvoiceRepository) AddNewItem(tx *sql.Tx, item model.ItemList) error {
	slog.Info("r Adding an item")

	query := `
        INSERT INTO item_lists (invoiceid, item, quantity, unitprice, subtotal, tax, total)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	log.Printf("Executing query: %s\n", query)
	log.Printf("With parameters: invoiceid=%s, item=%s, quantity=%d, unit_price=%f, subtotal=%f, tax=%f, total=%f",
		item.InvoiceId, item.Item, item.Quantity, item.UnitPrice, item.Subtotal, item.Tax, item.Total)

	_, err := tx.Exec(
		query,
		item.InvoiceId, item.Item, item.Quantity, item.UnitPrice, item.Subtotal, item.Tax, item.Total,
	)
	if err != nil {
		log.Printf("Database error on adding new item: %v\n", err)
	}
	return err
}

func GenerateInvoiceNumber(lastInvoiceNumber string) (string, error) {
	slog.Info("r Generating Invoice Number")
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
