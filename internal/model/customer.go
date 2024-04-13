package model

import "time"

type Customer struct {
	Id          string
	FirstName   string
	Email       string
	CompanyName string
	Phone       string
	Address     Address
	Invoices    []Invoice
	Lead        []Lead
	Notes       []Note
	CreatedAt   time.Time
	UpdatedAtAt time.Time
}

type Address struct {
	UnitNumber   string
	StreetNumber string
	StreetName   string
	City         string
	Postcode     string
}
