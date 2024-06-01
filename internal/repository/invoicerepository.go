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

func (r *InvoiceRepository) GetLastInvoiceId() (string, error) {
	var lastInvoiceId string
	query := `SELECT InvoiceId FROM invoices ORDER BY InvoiceId DESC LIMIT 1`

	err := r.db.QueryRow(query).Scan(&lastInvoiceId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No rows found, meaning this will be the first invoice
		}
		return "", err
	}
	return lastInvoiceId, nil
}

func (r *InvoiceRepository) generateInvoiceId(lastId string) (string, error) {
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
	lastInvoiceId, err := r.GetLastInvoiceId()
	if err != nil {
		return "", fmt.Errorf("error fetching last invoice ID: %v", err)
	}
	newInvoiceId, err := r.generateInvoiceId(lastInvoiceId)
	if err != nil {
		return "", fmt.Errorf("error generating new invoice ID: %v", err)
	}
	invoice.InvoiceId = newInvoiceId

	// PostgreSQL query with numbered placeholders
	query := `
        INSERT INTO invoices (InvoiceId, InvoiceNumber, InvoiceDate, DueDate, CustomerId, CustomerName, CompanyName, CustomerPhone, CustomerEmail, PaymentStatus)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	log.Printf("Executing query: %s\n", query)
	log.Printf("With parameters: InvoiceId=%s, InvoiceNumber=%s, InvoiceDate=%s, DueDate=%s, CustomerId=%d, CustomerName=%s, CompanyName=%s, CustomerPhone=%s, CustomerEmail=%s, PaymentStatus=%d",
		newInvoiceId, invoice.InvoiceNumber, invoice.InvoiceDate, invoice.DueDate, invoice.CustomerId, invoice.CustomerName, invoice.CompanyName, invoice.CustomerPhone, invoice.CustomerEmail, invoice.PaymentStatus)

	// Execute the query
	_, err = tx.Exec(
		query,
		invoice.InvoiceId, invoice.InvoiceNumber, invoice.InvoiceDate, invoice.DueDate, invoice.CustomerId, invoice.CustomerName, invoice.CompanyName, invoice.CustomerPhone, invoice.CustomerEmail, invoice.PaymentStatus,
	)
	if err != nil {
		return "", err
	}

	return invoice.InvoiceId, nil
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

func (r *InvoiceRepository) GetInvoiceById(invoiceId string) (*model.Invoice, error) {
	query := `
        SELECT InvoiceId, InvoiceNumber, InvoiceDate, DueDate, CustomerId, CustomerName, CompanyName, CustomerPhone, CustomerEmail, PaymentStatus
        FROM invoices
        WHERE InvoiceId = $1`

	row := r.db.QueryRow(query, invoiceId)

	var invoice model.Invoice
	err := row.Scan(&invoice.InvoiceId, &invoice.InvoiceNumber, &invoice.InvoiceDate, &invoice.DueDate, &invoice.CustomerId, &invoice.CustomerName, &invoice.CompanyName, &invoice.CustomerPhone, &invoice.CustomerEmail, &invoice.PaymentStatus)
	if err != nil {
		return nil, err
	}

	// Fetch items related to this invoice
	itemsQuery := `
        SELECT invoiceid, item, quantity, unitprice, subtotal, tax, total
        FROM item_lists
        WHERE invoiceid = $1`

	rows, err := r.db.Query(itemsQuery, invoiceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.ItemList
	for rows.Next() {
		var item model.ItemList
		err := rows.Scan(&item.InvoiceId, &item.Item, &item.Quantity, &item.UnitPrice, &item.Subtotal, &item.Tax, &item.Total)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	invoice.ItemList = items

	return &invoice, nil
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
