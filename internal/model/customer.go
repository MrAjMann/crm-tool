package model

import "time"

type Customer struct {
	Id          string
	FirstName   string
	Email       string
	CompanyName string
	Phone       string
	Lead        []Lead
	Notes       []Note
	CreatedAt   time.Time
	UpdatedAtAt time.Time
}
