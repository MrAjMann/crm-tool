package handler

import (
	"html/template"
	"log"
	"net/http"

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
		log.Printf("Error fetching customers: %v\n", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", customers)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}

}

func (h *CustomerHandler) AddCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
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

	tmpl, err := template.ParseFiles("src/templates/index.html")
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
