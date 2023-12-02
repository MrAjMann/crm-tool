package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

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



// Get all leads
func (h *LeadHandler) GetAllLeads(w http.ResponseWriter, r *http.Request) {
	leads, err := h.repo.GetAllLeads()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "leads.html", leads)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
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

// Get a Lead
func (h *LeadHandler) GetLead(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	base, idStr := path.Split(r.URL.Path)
	if base != "/lead/" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Get the lead id from request
	lead, err := h.repo.GetLeadById(idStr)
	if err != nil {
		http.Error(w, "Database error on fetching lead", http.StatusInternalServerError)
		log.Printf("Database error on fetching lead: %v\n", err)
		return
	}

	tmpl, err := template.ParseFiles("src/templates/lead.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Error loading template: %v\n", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "lead.html", lead)

	if err != nil {
		http.Error(w, "Error executing lead template", http.StatusInternalServerError)
		log.Printf("Error executing lead template: %v\n", err)
	}

}
