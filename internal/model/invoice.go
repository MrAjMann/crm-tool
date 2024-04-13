package model

import (
	"time"
)

type Invoice struct {
	InvoiceId       string
	InvoiceNumber   string
	InvoiceDate     time.Time
	DueDate         time.Time
	CustomerId      string
	CustomerName    string
	CompanyName     string
	CustomerPhone   string
	CustomerEmail   string
	PaymentStatus   PaymentStatus
	CustomerAddress Address
	ItemList        []ItemList
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ItemList struct {
	InvoiceId string
	Item      string
	Quantity  int32
	UnitPrice int32
	Subtotal  int32
	Tax       int32
	Total     int32
}
type PaymentStatus int

// This code will be used for the business logic, however i want to keep track roughly of my intention before I get to that point

// import (
//     "fmt"
//     "model"  - Just import the model folder
// )

const (
	Paid    PaymentStatus = iota // 0
	Pending                      // 1
	Overdue                      // 2
)

// func main() {
//     fmt.Println(model.Paid)    // Output: 0
//     fmt.Println(model.Pending) // Output: 1
//     fmt.Println(model.Overdue) // Output: 2

//     // Use the enum in a function
//     printStatus(model.Pending)
// }

// // Function that uses the PaymentStatus type
// func printStatus(status model.PaymentStatus) {
//     switch status {
//     case model.Paid:
//         fmt.Println("Payment is Paid")
//     case model.Pending:
//         fmt.Println("Payment is Pending")
//     case model.Overdue:
//         fmt.Println("Payment is Overdue")
//     default:
//         fmt.Println("Unknown Payment Status")
//     }
// }
