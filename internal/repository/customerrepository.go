package repository

import (
	"database/sql"
	"fmt"
	"strings"

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
	rows, err := repo.db.Query("SELECT Id, FirstName, LastName, Email ,Phone, CompanyName, Title, Website, Industry FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.CompanyName, &c.Title, &c.Website, &c.Industry); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, nil
}

// AddCustomer inserts a new customer into the database
func (repo *CustomerRepository) AddCustomer(customer model.Customer) (string, error) {
	var customerId string
	err := repo.db.QueryRow("INSERT INTO customers (FirstName, LastName, Email, Phone, CompanyName, Title, Website, Industry) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING Id",
		customer.FirstName, customer.LastName, customer.Email, customer.Phone, customer.CompanyName, customer.Title, customer.Website, customer.Industry).Scan(&customerId)
	if err != nil {
		return "", err
	}
	return customerId, nil
}

func (repo *CustomerRepository) GetCustomerById(id string) (model.Customer, error) {
	println(id)
	var customer model.Customer

	query := `SELECT Id, FirstName, LastName, Email, Phone, CompanyName, Title, Website, Industry 
						FROM customers
						WHERE Id = $1`

	err := repo.db.QueryRow(query, id).Scan(
		&customer.Id,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.Phone,
		&customer.CompanyName,
		&customer.Title,
		&customer.Website,
		&customer.Industry,
	)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (repo *CustomerRepository) SearchCustomers(query string) ([]model.Customer, error) {
	var customers []model.Customer

	// Adjust the SQL query to better handle searches for both first and last names together
	sqlQuery := `SELECT Id, FirstName, LastName, Email, Phone, CompanyName
                 FROM customers
                 WHERE CONCAT(FirstName, ' ', LastName) ILIKE $1 OR FirstName ILIKE $1 OR LastName ILIKE $1 OR Email ILIKE $1 OR Phone ILIKE $1 OR CompanyName ILIKE $1`
	// This allows for a more flexible search that considers both individual and full names.

	searchQuery := "%" + strings.TrimSpace(query) + "%"

	rows, err := repo.db.Query(sqlQuery, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("error querying customers with search query %s: %v", query, err)
	}
	defer rows.Close()

	for rows.Next() {
		var customer model.Customer
		if err := rows.Scan(&customer.Id, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone, &customer.CompanyName); err != nil {
			return nil, fmt.Errorf("error scanning customer: %v", err)
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating customer rows: %v", err)
	}

	return customers, nil
}

func (repo *CustomerRepository) DeleteCustomerById(id string) (model.Customer, error) {
	println(id)
	var customer model.Customer

	query := `SELECT Id, FirstName, LastName, Email, Phone, CompanyName
						FROM customers
						WHERE Id = $1`

	err := repo.db.QueryRow(query, id).Scan(
		&customer.Id,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.Phone,
		&customer.CompanyName,
	)
	if err != nil {
		return customer, err
	}
	// If the customer exists, proceed with deletion.
	deleteQuery := `DELETE FROM customers WHERE Id = $1`
	_, err = repo.db.Exec(deleteQuery, id)
	if err != nil {
		return model.Customer{}, err // Return an error if the delete operation fails
	}

	// Return the details of the deleted customer for confirmation/logging.
	return customer, nil

}
