package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MrAjMann/crm/internal/model"
	"github.com/MrAjMann/crm/internal/repository"
	"github.com/gorilla/sessions"
)

type CustomerHandler struct {
	repo *repository.CustomerRepository
	tmpl *template.Template
}

func NewCustomerHandler(repo *repository.CustomerRepository, tmpl *template.Template) *CustomerHandler {
	return &CustomerHandler{
		repo: repo,
		tmpl: tmpl,
	}
}

func httpError(w http.ResponseWriter, logMessage string, err error, statusCode int) {
	log.Printf("%s: %v", logMessage, err)
	http.Error(w, http.StatusText(statusCode), statusCode)
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
		LastName:    r.FormValue("lastName"),
		Email:       r.FormValue("email"),
		CompanyName: r.FormValue("companyName"),
		Phone:       r.FormValue("phone"),
		Title:       r.FormValue("title"),
		Website:     r.FormValue("website"),
		Industry:    r.FormValue("industry"),
	}

	customerId, err := h.repo.AddCustomer(customer)
	if err != nil {
		http.Error(w, "Database error on inserting new customer", http.StatusInternalServerError)
		log.Printf("Database error on inserting new customer: %v\n", err)
		return
	}
	customerIdInt, err := strconv.Atoi(customerId)
	if err != nil {
		log.Fatalf("Error converting customerId to int: %v", err)
	}

	tmpl, err := template.ParseFiles("src/templates/customers.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Error loading template: %v\n", err)
		return
	}

	// Redirect or display a success message
	err = tmpl.ExecuteTemplate(w, "customer-list-element", model.Customer{Id: customerIdInt, FirstName: customer.FirstName, LastName: customer.LastName, Email: customer.Email, CompanyName: customer.CompanyName, Phone: customer.Phone})
	if err != nil {
		httpError(w, "Error executing template", err, http.StatusInternalServerError)

	}
}

// Get a Customer
func (h *CustomerHandler) GetCustomerById(w http.ResponseWriter, r *http.Request) {
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
	customerId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}
	// Get the customer by id from the repository
	customer, err := h.repo.GetCustomerById(customerId)
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
}

// Search for a customer
func (h *CustomerHandler) HandleSearchCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("search")
	log.Printf("Received search query: '%s'", query)

	if query == "" {
		http.Error(w, "Query parameter 'search' is missing or empty", http.StatusBadRequest)
		return
	}

	customers, err := h.repo.SearchCustomers(query)
	if err != nil {
		http.Error(w, "Database error on fetching customers", http.StatusInternalServerError)
		log.Printf("Database error on fetching customers: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	var htmlOutput string

	htmlOutput += `
    <style>
        .customer-list { list-style-type: none; margin: 0; padding: 0; }
        .customer-item { background-color: #f9f9f9; border-left: 5px solid #007bff; margin-bottom: 8px; padding: 12px; border-radius: 4px; cursor: pointer; transition: background-color 0.3s; }
        .customer-item:hover { background-color: #f0f0f0; }
        .customer-info { margin: 0; color: #333; }
        .customer-info span { font-weight: bold; }
    </style>
    <ul class='customer-list'>`

	for _, customer := range customers {
		htmlOutput += fmt.Sprintf(`
			<li class='customer-item' data-customer-id="%d" onclick="selectCustomer(event)">
				<div class='customer-info'><span>CustomerId:</span> %d</div>
				<div class='customer-info'><span>Name:</span> %s %s</div>
				<div class='customer-info'><span>Email:</span> %s</div>
				<div class='customer-info'><span>Phone:</span> %s</div>
				<div class='customer-info'><span>Company:</span> %s</div>
			</li>`, customer.Id, customer.Id, customer.FirstName, customer.LastName, customer.Email, customer.Phone, customer.CompanyName)
	}
	htmlOutput += "</ul>"

	// Additional JavaScript for interactive functionality
	htmlOutput += `
    <script>
        document.querySelectorAll('.customer-item').forEach(item => {
            item.addEventListener('click', function() {
                
                document.querySelectorAll('.customer-item').forEach(i => {
                    i.style.borderLeft = '5px solid #007bff';
                });
                this.style.borderLeft = '5px solid #ff7f00';  // Highlight color change on click
                
                const customerId = this.querySelector('.customer-info:nth-child(1)').innerText.split(':')[1].trim();
					
                console.log('Selected Customer ID:', customerId);  
            });
        });
    </script>`

	if _, err = w.Write([]byte(htmlOutput)); err != nil {
		log.Printf("Failed to write data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

// Update a Customer

// Delete a Customer

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parsing the customer ID from the URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}
	idStr := parts[3] // Assuming URL is formatted as /customer/delete/{id}

	if idStr == "" {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}
	// Get the customer by id from the repository
	deletedCustomer, err := h.repo.DeleteCustomerById(idStr)
	if err != nil {
		http.Error(w, "Database error on fetching customer", http.StatusInternalServerError)
		log.Printf("Database error on fetching customer: %v\n", err)
		return
	}
	log.Printf("Deleted customer: %+v", deletedCustomer)
	// Assuming tmpl is a template instance parsed at application initialization
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Customer deleted successfully: %d - %s %s", deletedCustomer.Id, deletedCustomer.FirstName, deletedCustomer.LastName)
}

var store = sessions.NewCookieStore([]byte("askjdn23undm-dc2-3njdknwr"))

func (h *CustomerHandler) HandleSessionStore(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "customer-session")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Session retrieval failed", http.StatusInternalServerError)
		return
	}

	// Set some session values.
	session.Values["id"] = r.FormValue("id")
	session.Values["name"] = r.FormValue("name")
	session.Values["email"] = r.FormValue("email")

	// Save session
	if err := session.Save(r, w); err != nil {
		log.Printf("Failed to save session: %v", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte("Session updated successfully")); err != nil {
		log.Printf("Failed to write data: %v", err)
		return // Only log here, as the client already got the "200 OK" response code
	}
}

func (h *CustomerHandler) CheckAddress(customerId int) (*model.Address, error) {
	if customerId < 1 {
		return nil, fmt.Errorf("invalid customer ID provided")
	}

	address, err := h.repo.GetAddressByCustomerId(customerId)
	if err != nil {
		log.Printf("Error finding Address: %+v", address)
		return nil, err
	}

	// Assuming GetAddressByCustomerId handles nil address correctly
	log.Printf("Retrieved customer address: %+v", address)
	return address, nil
}
