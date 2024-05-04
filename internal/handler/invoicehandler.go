package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MrAjMann/crm/internal/model"
	"github.com/MrAjMann/crm/internal/repository"
)

type InvoiceHandler struct {
	repo            *repository.InvoiceRepository
	tmpl            *template.Template
	customerHandler *CustomerHandler
}

type InvoiceData struct {
	Invoices []model.Invoice
}

func NewInvoiceHandler(repo *repository.InvoiceRepository, tmpl *template.Template, customerHandler *CustomerHandler) *InvoiceHandler {
	return &InvoiceHandler{repo: repo, tmpl: tmpl}
}

func (h *InvoiceHandler) GetAllInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.repo.GetAllInvoices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := InvoiceData{
		Invoices: invoices,
	}

	err = h.tmpl.ExecuteTemplate(w, "invoices.html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func (h *InvoiceHandler) AddNewInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("error Method not allowed %v\n", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v\n", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	customerIdStr := r.FormValue("customerId")
	println(customerIdStr)
	if customerIdStr == "" {
		// Handle error: Customer ID is required
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}
	customerId, err := strconv.Atoi(customerIdStr)
	if err != nil {
		// Handle error: Invalid Customer ID
		http.Error(w, "Invalid Customer ID", http.StatusBadRequest)
		return
	}

	paymentStatusStr := r.FormValue("paymentStatus")
	if paymentStatusStr == "" {
		paymentStatusStr = "0"
	}
	paymentStatusInt, err := strconv.Atoi(paymentStatusStr)
	if err != nil {
		log.Printf("Invalid payment status: %v\n", err)
		http.Error(w, "Invalid payment status", http.StatusBadRequest)
		return
	}

	if paymentStatusInt < int(model.Unpaid) || paymentStatusInt > int(model.Overdue) {
		log.Printf("Payment status out of range: received %v\n", paymentStatusInt)
		http.Error(w, "Payment status out of range", http.StatusBadRequest)
		return
	}

	log.Printf("Received form data: %v", r.Form)
	// Create a new invoice from form values
	invoice := model.Invoice{
		CustomerId:    customerId,
		CustomerName:  r.FormValue("customerName"),
		DueDate:       time.Now().AddDate(0, 0, 30),
		CustomerEmail: r.FormValue("email"),
		CompanyName:   r.FormValue("companyName"),
		CustomerPhone: r.FormValue("phone"),
	}

	// Add the new invoice to the database
	invoiceId, err := h.repo.AddNewInvoice(invoice)
	if err != nil {
		log.Printf("Database error on creating new invoice: %v\n", err)
		http.Error(w, "Database error on creating new invoice", http.StatusInternalServerError)
		return
	}


	customerAddress, err := h.customerHandler.CheckAddress(customerId)
	if err != nil {
		log.Printf("Database error on fetching address: %v\n", err)
		http.Error(w, "Database error on fetching address", http.StatusInternalServerError)
		return
	}
	log.Printf("Deleted customer: %+v", customerAddress)
	// Prepare data for template rendering
	data := InvoiceData{
		Invoices: []model.Invoice{
			{
				InvoiceId:       invoiceId,
				InvoiceNumber:   "",
				InvoiceDate:     time.Now(),
				DueDate:         invoice.DueDate,
				CustomerId:      invoice.CustomerId,
				CustomerName:    invoice.CustomerName,
				CompanyName:     invoice.CompanyName,
				CustomerPhone:   invoice.CustomerPhone,
				CustomerEmail:   invoice.CustomerEmail,
				PaymentStatus:   model.PaymentStatus(paymentStatusInt),
				CustomerAddress: *customerAddress,
				ItemList:        []model.ItemList{},
			},
		},
	}
	log.Printf("Error executing template: %v\n", data)
	// Execute the template
	if err := h.tmpl.ExecuteTemplate(w, "invoice-list-element", data); err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func (h *InvoiceHandler) InvoiceCalculationHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Could not parse form: %v", err)
		http.Error(w, "<p>Error: Could not parse form.</p>", http.StatusBadRequest)
		return
	}

	quantityStr := r.FormValue("quantity")
	unitPriceStr := r.FormValue("unitPrice")
	log.Printf("Received - Quantity: %s, Unit Price: %s", quantityStr, unitPriceStr)

	if quantityStr == "" {
		http.Error(w, "<p>Error: Quantity is required.</p>", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "<p>Error: Invalid quantity.</p>", http.StatusBadRequest)
		return
	}
	if unitPriceStr == "" {
		http.Error(w, "<p>Error: UnitPrice is required.</p>", http.StatusBadRequest)
		return
	}

	unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
	if err != nil {
		log.Printf("Error parsing unit price: %v", err)
		http.Error(w, "<p>Error: Invalid unit price.</p>", http.StatusBadRequest)
		return
	}

	subtotal := float64(quantity) * unitPrice

	tax := subtotal * 0.10 // 10% tax

	total := subtotal + tax
	// log.Printf("Quantity: $%.2f, UnitPrice:$%.2f,Subtotal: $%.2f, Tax: $%.2f, Total: $%.2f", float64(quantity), float64(unitPrice), subtotal, tax, total)

	// Prepare HTML response
	w.Header().Set("Content-Type", "application/json")
	// return as a json object
	jsonResponse := fmt.Sprintf(`
        {
            "subtotal": "%.2f",
            "tax": "%.2f",
            "total": "%.2f"
        }`, subtotal, tax, total)
	fmt.Fprint(w, jsonResponse)
}
