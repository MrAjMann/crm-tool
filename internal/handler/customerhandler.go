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

	// Assuming you want to redirect or display a success message
	err = tmpl.ExecuteTemplate(w, "customer-list-element", model.Customer{Id: customerIdInt, FirstName: customer.FirstName, LastName: customer.LastName, Email: customer.Email, CompanyName: customer.CompanyName, Phone: customer.Phone})
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
}

// Search for a customer
func (h *CustomerHandler) HandleSearchCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("search")
	log.Printf("Received search query: '%s'", query) // This logs the received query

	if query == "" {
		http.Error(w, "Query parameter 'search' is missing or empty", http.StatusBadRequest)
		return
	}
	customers, err := h.repo.SearchCustomers(query) // Assuming Query method accepts a search string
	if err != nil {
		http.Error(w, "Database error on fetching customers", http.StatusInternalServerError)
		log.Printf("Database error on fetching customers: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	var htmlOutput string

	// Improved HTML output with added CSS for styling and JavaScript for click functionality
	htmlOutput += `
<style>
    .customer-list { list-style-type: none; padding-left: 0; }
    .customer-item { background: #f9f9f9; border: 1px solid #ddd; margin-top: 8px; padding: 8px 16px; border-radius: 4px; cursor: pointer; }
    .customer-item:hover { background-color: #f0f0f0; }
    .customer-info { margin: 0; }
    .customer-info span { font-weight: bold; }
</style>
<ul class='customer-list'>`
	for _, customer := range customers {
		htmlOutput += fmt.Sprintf(`
<li class='customer-item'>
    <p class='customer-info'><span>Id:</span> %d </p>
    <p class='customer-info'><span>Name:</span> %s %s</p>
    <p class='customer-info'><span>Email:</span> %s</p>
    <p class='customer-info'><span>Phone:</span> %s</p>
    <p class='customer-info'><span>Company:</span> %s</p>
</li>`, customer.Id, customer.FirstName, customer.LastName, customer.Email, customer.Phone, customer.CompanyName)
	}
	htmlOutput += "</ul>"
	// Temp store, append to invoice data
	htmlOutput += `
	<script>
		document.querySelectorAll('.customer-item').forEach(item => {
			item.addEventListener('click', function() {
				// Change the background color of only the clicked item
				document.querySelectorAll('.customer-item').forEach(i => {
					i.style.background = ''; // Reset background for all items
				});
				this.style.background = "#ff7f00";
	
				// Fetching the customer details correctly
				var id = this.querySelector('.customer-info:nth-child(1)').innerText;
				var name = this.querySelector('.customer-info:nth-child(2)').innerText;
				var email = this.querySelector('.customer-info:nth-child(3)').innerText;
	
				// Sending data to server to store in session
				fetch('./customer-session', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/x-www-form-urlencoded',
					},
					body: 'id=' + encodeURIComponent(id) + '&name=' + encodeURIComponent(name) + '&email=' + encodeURIComponent(email)
					
				})
				.then(response => response.text())
				.then(data => console.log(data))
				.catch(error => console.error('Error:', error));
			});
		});
	</script>`

	w.Write([]byte(htmlOutput))

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

	session, _ := store.Get(r, "customer-session")
	// Set some session values.
	session.Values["id"] = r.FormValue("id")
	session.Values["name"] = r.FormValue("name")
	session.Values["email"] = r.FormValue("email")

	// Save session
	if err := session.Save(r, w); err != nil {
		fmt.Printf("Error: %s", err)
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session updated successfully"))
}
