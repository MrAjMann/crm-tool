package handler

import (
	"log"
	"net/http"
)

func (h *InvoiceHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	invoiceId := r.FormValue("invoiceId")
	// Add logic to generate PDF for the invoice with invoiceId
	log.Printf("Generating PDF for invoice ID: %s", invoiceId)

	// Redirect to /invoices after generating the PDF
	http.Redirect(w, r, "/invoices", http.StatusSeeOther)
}
