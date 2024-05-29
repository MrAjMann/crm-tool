package handler

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/MrAjMann/crm/internal/model"
	"github.com/MrAjMann/crm/internal/repository"
)

type InvoiceHandler struct {
	repo            *repository.InvoiceRepository
	tmpl            *template.Template
	customerHandler *CustomerHandler
}

type InvoiceData struct {
	Invoices []model.Invoice
}

func NewInvoiceHandler(repo *repository.InvoiceRepository, tmpl *template.Template, customerHandler *CustomerHandler) *InvoiceHandler {
	return &InvoiceHandler{repo: repo, tmpl: tmpl}
}

func (h *InvoiceHandler) GetAllInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.repo.GetAllInvoices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := InvoiceData{
		Invoices: invoices,
	}

	err = h.tmpl.ExecuteTemplate(w, "invoices.html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func (h *InvoiceHandler) AddNewInvoice(w http.ResponseWriter, r *http.Request) {
	slog.Info("h adding invoice")
	log.Printf("AddNewInvoice called with method %s", r.Method)
	if r.Method != "POST" {
		log.Printf("error Method not allowed %v\n", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v\n", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	customerIdStr := r.FormValue("customerId")
	if customerIdStr == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	customerId, err := strconv.Atoi(customerIdStr)
	if err != nil {
		http.Error(w, "Invalid Customer ID", http.StatusBadRequest)
		return
	}

	paymentStatusStr := r.FormValue("paymentStatus")
	if paymentStatusStr == "" {
		paymentStatusStr = "0"
	}
	paymentStatusInt, err := strconv.Atoi(paymentStatusStr)
	if err != nil {
		log.Printf("Invalid payment status: %v\n", err)
		http.Error(w, "Invalid payment status", http.StatusBadRequest)
		return
	}

	if paymentStatusInt < int(model.Unpaid) || paymentStatusInt > int(model.Overdue) {
		log.Printf("Payment status out of range: received %v\n", paymentStatusInt)
		http.Error(w, "Payment status out of range", http.StatusBadRequest)
		return
	}

	// Get the list of items
	itemList, err := h.AddItemsToInvoice(r)
	if err != nil {
		log.Printf("Error adding items to invoice: %v\n", err)
		http.Error(w, "Error adding items to invoice", http.StatusBadRequest)
		return
	}

	customerInfo, err := h.customerHandler.repo.GetCustomerById(customerId)
	if err != nil {
		log.Printf("Database error on fetching address: %v\n", err)
		return
	}

	customerName := fmt.Sprintf("%s %s", customerInfo.FirstName, customerInfo.LastName)

	// Create a new invoice from form values
	invoice := model.Invoice{
		CustomerId:    customerId,
		CustomerName:  getFormValueOrDefault(r, "customerName", customerName),
		DueDate:       time.Now().AddDate(0, 0, 30),
		CustomerEmail: getFormValueOrDefault(r, "email", customerInfo.Email),
		CompanyName:   getFormValueOrDefault(r, "companyName", customerInfo.CompanyName),
		CustomerPhone: getFormValueOrDefault(r, "phone", customerInfo.Phone),
		ItemList:      itemList,
	}

	// Begin a transaction
	tx, err := h.repo.BeginTransaction()
	if err != nil {
		log.Printf("Database error on beginning transaction: %v\n", err)
		http.Error(w, "Database error on creating new invoice", http.StatusInternalServerError)
		return
	}

	// Add the new invoice to the database
	invoiceId, err := h.repo.AddNewInvoice(tx, invoice)
	if err != nil {
		tx.Rollback()
		log.Printf("Database error on creating new invoice: %v\n", err)
		http.Error(w, "Database error on creating new invoice", http.StatusInternalServerError)
		return
	}

	// Convert invoiceId to string for the InvoiceId field in ItemList
	invoiceIdStr := invoiceId

	// Add items to the invoice
	for _, item := range itemList {
		item.InvoiceId = invoiceIdStr // Set the invoice ID for each item as a string
		if err := h.repo.AddNewItem(tx, item); err != nil {
			tx.Rollback()
			log.Printf("Database error on adding new item: %v\n", err)
			http.Error(w, "Database error on adding new item", http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Database error on committing transaction: %v\n", err)
		http.Error(w, "Database error on committing transaction", http.StatusInternalServerError)
		return
	}

	customerAddress, err := h.customerHandler.CheckAddress(customerId)
	if err != nil {
		log.Printf("Database error on fetching address: %v\n", err)
		return
	}

	// Prepare data for template rendering
	data := InvoiceData{
		Invoices: []model.Invoice{
			{
				InvoiceId:       invoiceIdStr, // Use the string version of invoiceId
				InvoiceNumber:   "",
				InvoiceDate:     time.Now(),
				DueDate:         invoice.DueDate,
				CustomerId:      invoice.CustomerId,
				CustomerName:    invoice.CustomerName,
				CompanyName:     invoice.CompanyName,
				CustomerPhone:   invoice.CustomerPhone,
				CustomerEmail:   invoice.CustomerEmail,
				PaymentStatus:   model.PaymentStatus(paymentStatusInt),
				CustomerAddress: *customerAddress,
				ItemList:        itemList,
			},
		},
	}
	log.Printf("Template Data: %v\n", data)

	// Execute the template
	if err := h.tmpl.ExecuteTemplate(w, "invoice-list-element", data); err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

// Helper function to get form value or default
func getFormValueOrDefault(r *http.Request, formField, defaultValue string) string {
	value := r.FormValue(formField)
	slog.Info(value)
	if value == "" {
		return defaultValue
	}
	return value
}

// Add Items to Invoice

func (h *InvoiceHandler) AddItemsToInvoice(r *http.Request) ([]model.ItemList, error) {
	slog.Info("h adding items to invoice")
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v\n", err)
		return nil, err
	}

	var itemList []model.ItemList

	itemNames := r.Form["item"]
	quantities := r.Form["quantity"]
	unitPrices := r.Form["unitPrice"]
	subTotals := r.Form["subTotal"]
	taxes := r.Form["tax"]
	totals := r.Form["total"]

	// Create a formatted string to log all values
	var formattedString string
	for i := 0; i < len(itemNames); i++ {
		formattedString += fmt.Sprintf(
			"Item %d: {item: %s, quantity: %s, unitPrice: %s, subTotal: %s, tax: %s, total: %s}\n",
			i+1, itemNames[i], quantities[i], unitPrices[i], subTotals[i], taxes[i], totals[i],
		)
	}
	log.Printf("Form values:\n%s", formattedString)

	// Validate and process each item
	for i := 0; i < len(itemNames); i++ {
		// Check if all necessary fields are present and not empty
		if itemNames[i] == "" || quantities[i] == "" || unitPrices[i] == "" || subTotals[i] == "" || taxes[i] == "" || totals[i] == "" {
			return nil, fmt.Errorf("form values missing or empty for item %d", i+1)
		}

		// Convert and validate quantity
		quantity, err := strconv.Atoi(quantities[i])
		if err != nil {
			return nil, fmt.Errorf("invalid quantity for item %d", i+1)
		}

		// Convert and validate unit price
		unitPrice, err := strconv.Atoi(unitPrices[i])
		if err != nil {
			return nil, fmt.Errorf("invalid unit price for item %d", i+1)
		}

		// Convert and validate subtotal
		subTotal, err := strconv.Atoi(subTotals[i])
		if err != nil {
			return nil, fmt.Errorf("invalid subtotal for item %d", i+1)
		}

		// Convert and validate tax
		tax, err := strconv.Atoi(taxes[i])
		if err != nil {
			return nil, fmt.Errorf("invalid tax for item %d", i+1)
		}

		// Convert and validate total
		total, err := strconv.Atoi(totals[i])
		if err != nil {
			return nil, fmt.Errorf("invalid total for item %d", i+1)
		}

		item := model.ItemList{
			Item:      itemNames[i],
			Quantity:  int32(quantity),
			UnitPrice: int32(unitPrice),
			Subtotal:  int32(subTotal),
			Tax:       int32(tax),
			Total:     int32(total),
		}

		itemList = append(itemList, item)
	}

	return itemList, nil
}

// CALCULATION

func (h *InvoiceHandler) InvoiceCalculationHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("h calculating invoice items")
	if err := r.ParseForm(); err != nil {
		log.Printf("Could not parse form: %v", err)
		http.Error(w, "<p>Error: Could not parse form.</p>", http.StatusBadRequest)
		return
	}

	quantityStr := r.FormValue("quantity")
	unitPriceStr := r.FormValue("unitPrice")
	log.Printf("Received - Quantity: %s, Unit Price: %s", quantityStr, unitPriceStr)

	if quantityStr == "" {
		http.Error(w, "<p>Error: Quantity is required.</p>", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "<p>Error: Invalid quantity.</p>", http.StatusBadRequest)
		return
	}
	if unitPriceStr == "" {
		http.Error(w, "<p>Error: UnitPrice is required.</p>", http.StatusBadRequest)
		return
	}

	unitPrice, err := strconv.ParseFloat(unitPriceStr, 64)
	if err != nil {
		log.Printf("Error parsing unit price: %v", err)
		http.Error(w, "<p>Error: Invalid unit price.</p>", http.StatusBadRequest)
		return
	}

	subtotal := float64(quantity) * unitPrice
	tax := subtotal * 0.10 // 10% tax
	total := subtotal + tax

	// Prepare JSON response
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := fmt.Sprintf(`{
		"subtotal": "%.2f",
		"tax": "%.2f",
		"total": "%.2f"
	}`, subtotal, tax, total)
	fmt.Fprint(w, jsonResponse)
}
