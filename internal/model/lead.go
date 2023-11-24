package model

type Lead struct {
	leadId      string
	firstName   string
	lastName    string
	email       string
	companyName string
	phone       string
	title       string
	website     string
	industry    string
	source      string
	notes       []Note
}


