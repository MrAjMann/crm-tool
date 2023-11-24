package model

import (
	"time"
)

type Status string
type ClosedStatus string

const (
	NewLead       Status = "New Lead"
	Contacted     Status = "Contacted"
	Engaged       Status = "Engaged"
	Qualified     Status = "Qualified"
	NeedsAnalysis Status = "Needs Analysis"
	ProposalSent  Status = "Proposal Sent"
	Negotiation   Status = "Negotiation"
	Closed        Status = "Closed"
)

const (
	Won  ClosedStatus = "Won"
	Lost ClosedStatus = "Lost"
)

type Lead struct {
	LeadId      string
	FirstName   string
	LastName    string
	Email       string
	CompanyName string
	Phone       string
	Status      Status
	Title       string
	Website     string
	Industry    string
	Source      string
	Notes       []Note
	CreatedAt   time.Time
	UpdatedAtAt time.Time
}
