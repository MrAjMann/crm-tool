package model

import (
	"time"
)

type Status struct {
	StatusId          int    // Unique identifier for the status
	StatusValue       string // Descriptive name of the status
	IsClosed          bool   // Indicates if the status is a closed state
	ClosedStatusValue string // Specific closed status description, if applicable
}

type Lead struct {
	LeadId      int
	FirstName   string
	LastName    string
	Email       string
	CompanyName string
	Phone       string
	Status      Status
	Title       string
	Website     string
	Industry    string
	ServiceType string
	Source      string
	Notes       []Note
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
