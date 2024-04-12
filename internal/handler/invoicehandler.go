package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/MrAjMann/crm/internal/repository"
)

type InvoiceHandler struct {
	repo *repository.InvoiceRepository
	tmpl *template.Template
}

func NewInvoiceHandler(repo *repository.InvoiceRepository, tmpl *template.Template) *InvoiceHandler {
	return &InvoiceHandler{repo: repo, tmpl: tmpl}
}

func (h *InvoiceHandler) GetAllInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.repo.GetAllInvoices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = h.tmpl.ExecuteTemplate(w, "invoices.html", invoices)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
