package repository

import (
	"database/sql"

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
	err := repo.db.QueryRow("INSERT INTO invoices ( InvoiceNumber, InvoiceDate, DueDate, CustomerId, CustomerName, CompanyName, CustomerPhone, CustomerEmail, PaymentStatus ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING InvoiceId").Scan(&invoiceId)
	if err != nil {
		return "", err
	}
	return invoiceId, nil
}
