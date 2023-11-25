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

func (repo *LeadRepository) GetAllLeads() ([]model.Lead, error) {
	rows, err := repo.db.Query("SELECT Id, FirstName, Email, CompanyName, Phone FROM leads")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []model.Lead
	for rows.Next() {
		var l model.Lead
		if err := rows.Scan(&l.LeadId, &l.FirstName, &l.Email, &l.CompanyName, &l.Phone); err != nil {
			return nil, err
		}
		leads = append(leads, l)
	}

	return leads, nil
}


