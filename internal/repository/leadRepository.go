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

// Example function to fetch all leads
func (repo *LeadRepository) GetAllLeads() ([]model.Lead, error) {
	rows, err := repo.db.Query("SELECT Id, FirstName,Lastname, Email, CompanyName, Phone, Title, Website, Industry, Source  FROM leads")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []model.Lead
	for rows.Next() {
		var lead model.Lead
		if err := rows.Scan(
			&lead.LeadId,
			&lead.FirstName,
			&lead.LastName,
			&lead.Email,
			&lead.CompanyName,
			&lead.Phone,
			&lead.Website,
			&lead.Title,
			&lead.Industry,
			&lead.Source); err != nil {
			return nil, err
		}
		leads = append(leads, lead)
	}

	return leads, nil
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

func (repo *LeadRepository) GetLeadById(id string) (model.Lead, error) {
	println(id)
	var lead model.Lead

	query := `SELECT Id, FirstName, LastName, Email, Phone, CompanyName, Website, Title, Industry, Source 
						FROM leads
						WHERE Id = $1`

	err := repo.db.QueryRow(query, id).Scan(
		&lead.LeadId,
		&lead.FirstName,
		&lead.LastName,
		&lead.Email,
		&lead.Phone,
		&lead.CompanyName,
		&lead.Website,
		&lead.Title,
		&lead.Industry,
		&lead.Source,
	)
	if err != nil {
		return lead, err
	}
	return lead, nil
}
