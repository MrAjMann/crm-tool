package model

import "time"

type Customer struct {
	Id                 int
	FirstName          string
	LastName           string
	CompanyName        string
	Email              string
	Phone              string
	Title              string
	Website            string
	Industry           string
	InitialServiceType string
	CurrentServiceType string
	Address            *Address
	Invoices           []Invoice
	LeadId             int
	Notes              []Note
	ServiceHistory     []ServiceEntry
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
type Address struct {
	CustomerId      int
	UnitNumber   string
	StreetNumber string
	StreetName   string
	City         string
	State        string
	Postcode     string
}

type ServiceEntry struct {
	ServiceType string
	StartDate   time.Time
	EndDate     time.Time // Can be nil if currently active
}
