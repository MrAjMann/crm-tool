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
	repo *repository.InvoiceRepository
	tmpl *template.Template
}

type InvoiceData struct {
	Invoices []model.Invoice
}

func NewInvoiceHandler(repo *repository.InvoiceRepository, tmpl *template.Template) *InvoiceHandler {
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

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		log.Printf("Error parsing form: %v\n", err)
		return
	}

	paymentStatusStr := r.FormValue("paymentStatus")
	paymentStatusInt, err := strconv.Atoi(paymentStatusStr)
	if err != nil {
		http.Error(w, "Invalid payment status", http.StatusBadRequest)
		log.Printf("Invalid payment status: %v\n", err)
		return
	}

	// Validate payment status
	if paymentStatusInt < int(model.Paid) || paymentStatusInt > int(model.Overdue) {
		http.Error(w, "Payment status out of range", http.StatusBadRequest)
		log.Printf("Payment status out of range: received %v\n", paymentStatusInt)
		return
	}

	invoice := model.Invoice{
		CustomerName:  r.FormValue("customerName"),
		DueDate:       time.Now().AddDate(0, 0, 30),
		CustomerEmail: r.FormValue("email"),
		CompanyName:   r.FormValue("companyName"),
		CustomerPhone: r.FormValue("phone"),
	}

	invoiceId, err := h.repo.AddNewInvoice(invoice)
	if err != nil {
		http.Error(w, "Database error on creating new invoice", http.StatusInternalServerError)
		log.Printf("Database error on creating new invoice: %v\n", err)
		return
	}

	// Corrected to use the structure InvoiceData with a slice of Invoices
	data := InvoiceData{
		Invoices: []model.Invoice{
			{
				InvoiceId:     invoiceId, // Now including InvoiceId
				InvoiceNumber: invoice.InvoiceNumber,
				InvoiceDate:   invoice.InvoiceDate,
				DueDate:       invoice.DueDate,
				CustomerEmail: invoice.CustomerEmail,
				CompanyName:   invoice.CompanyName,
				CustomerPhone: invoice.CustomerPhone,
				PaymentStatus: invoice.PaymentStatus,
			},
		},
	}

	err = h.tmpl.ExecuteTemplate(w, "invoice-list-element", data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Printf("Error executing template: %v\n", err)
	}
}

func (h *InvoiceHandler) InvoiceCalculationHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Printf("Could not parse form: %v", err)
		http.Error(w, `{"error": "Could not parse form"}`, http.StatusBadRequest)
		return
	}

	quantityStr := r.FormValue("quantity")
	unitPriceStr := r.FormValue("unitPrice")
	log.Printf("Received - Quantity: %s, Unit Price: %s", quantityStr, unitPriceStr)

	if quantityStr == "" {
		http.Error(w, "Quantity is required", http.StatusBadRequest)
		return
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	unitPrice, err := strconv.ParseFloat(r.FormValue("unitPrice"), 64)
	if err != nil {
		log.Printf("Error parsing unit price: %v", err)
		http.Error(w, `{"error": "Invalid unit price"}`, http.StatusBadRequest)
		return
	}

	subtotal := float64(quantity) * unitPrice
	tax := subtotal * 0.10 // 10% tax
	total := subtotal + tax

	// Prepare HTML response
	htmlResponse := fmt.Sprintf(`
        <div>
            <p>Subtotal: $%.2f</p>
            <p>Tax: $%.2f</p>
            <p>Total: $%.2f</p>
        </div>
    `, subtotal, tax, total)

	// Write HTML response to the client
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlResponse)

}
