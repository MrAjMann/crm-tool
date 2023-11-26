package repository

import (
	"database/sql"

	"github.com/MrAjMann/crm/internal/model"
)

type LeadRepository struct {
	db *sql.DB
}

func NewLeadRepository(db *sql.DB) *LeadRepository {
	return &LeadRepository{db: db}
}

// Addlead inserts a new lead into the database
func (repo *LeadRepository) AddLead(lead model.Lead) (string, error) {
	var leadId string
	err := repo.db.QueryRow("INSERT INTO leads (FirstName, LastName, Email, CompanyName, Phone, Title, Website, Industry, Source ) VALUES ($1, $2, $3, $4, $5,$6,$7, $8, $9) RETURNING Id",
		lead.FirstName, lead.LastName, lead.Email, lead.CompanyName, lead.Phone, lead.Title, lead.Website, lead.Industry, lead.Source).Scan(&leadId)
	if err != nil {
		return "", err
	}
	return leadId, nil
}
