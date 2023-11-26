package handler

import (
	"html/template"
	"log"
	"net/http"
)

type DashboardHandler struct {
	tmpl *template.Template
}

func NewDashboardHandler(tmpl *template.Template) *DashboardHandler {
	return &DashboardHandler{tmpl: tmpl}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}
