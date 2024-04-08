package handler

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/MrAjMann/crm/internal/model"
	"github.com/MrAjMann/crm/internal/repository"
)

type CustomerHandler struct {
	repo *repository.CustomerRepository
	tmpl *template.Template
}

func NewCustomerHandler(repo *repository.CustomerRepository, tmpl *template.Template) *CustomerHandler {
	return &CustomerHandler{repo: repo, tmpl: tmpl}
}

func (h *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.repo.GetAllCustomers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "customers.html", customers)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}

}

// Add a Customer

func (h *CustomerHandler) AddCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("Method not allowed: %v\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		log.Printf("Error parsing form: %v\n", err)
		return
	}

	customer := model.Customer{
		FirstName:   r.FormValue("firstName"),
		Email:       r.FormValue("email"),
		CompanyName: r.FormValue("companyName"),
		Phone:       r.FormValue("phone"),
	}

	customerId, err := h.repo.AddCustomer(customer)
	if err != nil {
		http.Error(w, "Database error on inserting new customer", http.StatusInternalServerError)
		log.Printf("Database error on inserting new customer: %v\n", err)
		return
	}

	tmpl, err := template.ParseFiles("src/templates/customers.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Error loading template: %v\n", err)
		return
	}

	// Assuming you want to redirect or display a success message
	err = tmpl.ExecuteTemplate(w, "customer-list-element", model.Customer{Id: customerId, FirstName: customer.FirstName, Email: customer.Email, CompanyName: customer.CompanyName, Phone: customer.Phone})
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Printf("Error executing template: %v\n", err)
	}
}

// Get a Customer
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/customer/")
	idStr = strings.TrimSuffix(idStr, "/") // Optional, based on URL structure

	if idStr == "" {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	// Get the customer by id from the repository
	customer, err := h.repo.GetCustomerById(idStr)
	if err != nil {
		http.Error(w, "Database error on fetching customer", http.StatusInternalServerError)
		log.Printf("Database error on fetching customer: %v\n", err)
		return
	}

	// Assuming tmpl is a template instance parsed at application initialization
	err = h.tmpl.ExecuteTemplate(w, "customer.html", customer)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Printf("Error executing template: %v\n", err)
	}

	// err = tmpl.ExecuteTemplate(w, "customer.html", customer)
	// log.Println(customer)
	// if err != nil {
	// 	http.Error(w, "Error executing template", http.StatusInternalServerError)
	// 	log.Printf("Error executing template: %v\n", err)
	// }

	// Call the repository function to get the customer
	// Execute the template with the customer data
}

// Update a Customer

// Delete a Customer
