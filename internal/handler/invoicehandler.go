package handler

import (
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
