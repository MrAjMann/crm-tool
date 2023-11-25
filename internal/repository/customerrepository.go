package repository

import (
	"database/sql"

	"github.com/MrAjMann/crm/internal/model"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Example function to fetch all customers
func (repo *CustomerRepository) GetAllCustomers() ([]model.Customer, error) {
	rows, err := repo.db.Query("SELECT Id, FirstName, Email, CompanyName, Phone FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(&c.Id, &c.FirstName, &c.Email, &c.CompanyName, &c.Phone); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, nil
}

// AddCustomer inserts a new customer into the database
func (repo *CustomerRepository) AddCustomer(customer model.Customer) (string, error) {
	var customerId string
	err := repo.db.QueryRow("INSERT INTO customers (FirstName, Email, CompanyName, Phone) VALUES ($1, $2, $3, $4) RETURNING Id",
		customer.FirstName, customer.Email, customer.CompanyName, customer.Phone).Scan(&customerId)
	if err != nil {
		return "", err
	}
	return customerId, nil
}

func (repo *CustomerRepository) GetCustomerById(id string) (model.Customer, error) {
	println(id)
	var customer model.Customer

	query := `SELECT Id, FirstName,  Email, Phone, CompanyName
						FROM customers
						WHERE Id = $1`

	err := repo.db.QueryRow(query, id).Scan(
		&customer.Id,
		&customer.FirstName,
		&customer.Email,
		&customer.Phone,
		&customer.CompanyName,
	)
	if err != nil {
		return customer, err
	}
	return customer, nil
}
