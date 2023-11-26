package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/MrAjMann/crm/internal/model"
	"github.com/MrAjMann/crm/internal/repository"
)

type LeadHandler struct {
	repo *repository.LeadRepository
	tmpl *template.Template
}

func NewLeadHandler(repo *repository.LeadRepository, tmpl *template.Template) *LeadHandler {
	return &LeadHandler{repo: repo, tmpl: tmpl}
}

func (h *LeadHandler) AddLead(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	println(r.Method)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		println("Error parsing form")
		return
	}

	lead := model.Lead{
		FirstName:   r.FormValue("firstName"),
		LastName:    r.FormValue("lastName"),
		Email:       r.FormValue("email"),
		CompanyName: r.FormValue("companyName"),
		Phone:       r.FormValue("phone"),
		Title:       r.FormValue("title"),
		Website:     r.FormValue("website"),
		Industry:    r.FormValue("industry"),
		Source:      r.FormValue("source"),
	}

	leadId, err := h.repo.AddLead(lead)
	println(leadId)
	if err != nil {
		log.Printf("Database error on inserting new lead: %v\n", err)
		http.Error(w, "Database error on inserting new lead", http.StatusInternalServerError)
		return
	}

	// Redirect to the newly created lead's page
	http.Redirect(w, r, fmt.Sprintf("/lead/%s", leadId), http.StatusSeeOther)

}
