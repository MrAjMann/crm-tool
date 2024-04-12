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
	rows, err := repo.db.Query("SELECT InvoiceDate, DueDate, CustomerId, CustomerName, CompanyName FROM invoices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []model.Invoice
	for rows.Next() {
		var i model.Invoice
		if err := rows.Scan(&i.InvoiceDate, &i.DueDate, &i.CustomerId, &i.CustomerName, &i.CompanyName); err != nil {
			return nil, err
		}
		invoices = append(invoices, i)
	}
	return invoices, nil

}
