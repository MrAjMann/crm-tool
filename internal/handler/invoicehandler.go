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

type InvoiceData struct {
	Invoices []model.Invoice
}
type InvoiceHandler struct {
	repo            *repository.InvoiceRepository
	customerHandler *CustomerHandler
	tmpl            *template.Template
}

func NewInvoiceHandler(repo *repository.InvoiceRepository, tmpl *template.Template, customerHandler *CustomerHandler) *InvoiceHandler {
	return &InvoiceHandler{
		repo:            repo,
		customerHandler: customerHandler,
		tmpl:            tmpl,
	}
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
	log.Printf("AddNewInvoice called with method %s", r.Method)

	if r.Method != "POST" {
		log.Printf("Error: Method not allowed %v\n", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v\n", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// log.Printf("Form values: %v\n", r.Form)

	customerId, err := getCustomerIdFromRequest(r)
	if err != nil {
		log.Printf("Error getting customer ID: %v\n", err)
		http.Error(w, "Invalid Customer ID", http.StatusBadRequest)
		return
	}
	log.Printf("getCustomerIdFromRequest %d", customerId)

	paymentStatusInt, err := getPaymentStatusFromRequest(r)
	if err != nil {
		log.Printf("Error getting payment status: %v\n", err)
		http.Error(w, "Invalid payment status", http.StatusBadRequest)
		return
	}
	log.Printf("getPaymentStatusFromRequest %d", paymentStatusInt)
	itemList, err := h.AddItemsToInvoice(r)
	if err != nil {
		log.Printf("Error adding items to invoice: %v\n", err)
		http.Error(w, "Error adding items to invoice", http.StatusBadRequest)
		return
	}

	if h.customerHandler == nil {
		log.Printf("Error: customerHandler is nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if h.customerHandler.repo == nil {
		log.Printf("Error: customerHandler.repo is nil")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	customerInfo, err := h.customerHandler.repo.GetCustomerById(customerId)
	if err != nil {
		log.Printf("Database error on fetching customer info: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if customerInfo == nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	invoice, err := createInvoiceFromRequest(r, customerInfo, itemList)
	if err != nil {
		log.Printf("Error creating invoice: %v\n", err)
		http.Error(w, "Error creating invoice", http.StatusInternalServerError)
		return
	}

	if err := h.saveInvoiceWithItems(invoice, itemList); err != nil {
		log.Printf("Error saving invoice with items: %v\n", err)
		http.Error(w, "Error saving invoice", http.StatusInternalServerError)
		return
	}

	customerAddress, err := h.customerHandler.CheckAddress(customerId)
	if err != nil {
		log.Printf("Database error on fetching address: %v\n", err)
		return
	}

	data := prepareInvoiceData(invoice, itemList, customerAddress, paymentStatusInt)



    if err := h.tmpl.ExecuteTemplate(w, "modal", data); err != nil {
        log.Printf("Error executing modal template: %v\n", err)
        http.Error(w, "Error executing template", http.StatusInternalServerError)
    }
}

func getCustomerIdFromRequest(r *http.Request) (int, error) {
	customerIdStr := r.FormValue("customerId")
	if customerIdStr == "" {
		return 0, fmt.Errorf("customer ID is required")
	}

	customerId, err := strconv.Atoi(customerIdStr)
	if err != nil {
		return 0, fmt.Errorf("invalid Customer ID")
	}
	return customerId, nil
}

func getPaymentStatusFromRequest(r *http.Request) (int, error) {
	paymentStatusStr := r.FormValue("paymentStatus")
	if paymentStatusStr == "" {
		paymentStatusStr = "0"
	}
	paymentStatusInt, err := strconv.Atoi(paymentStatusStr)
	if err != nil {
		return 0, fmt.Errorf("invalid payment status")
	}
	if paymentStatusInt < int(model.Unpaid) || paymentStatusInt > int(model.Overdue) {
		return 0, fmt.Errorf("payment status out of range")
	}
	return paymentStatusInt, nil
}

func createInvoiceFromRequest(r *http.Request, customerInfo *model.Customer, itemList []model.ItemList) (*model.Invoice, error) {
	customerName := fmt.Sprintf("%s %s", customerInfo.FirstName, customerInfo.LastName)

	customerNameFormValue := r.FormValue("customerName")
	if customerNameFormValue != "" {
		customerName = customerNameFormValue
	}

	customerEmail := customerInfo.Email
	if r.FormValue("email") != "" {
		customerEmail = r.FormValue("email")
	}

	companyName := customerInfo.CompanyName
	if r.FormValue("companyName") != "" {
		companyName = r.FormValue("companyName")
	}

	customerPhone := customerInfo.Phone
	if r.FormValue("phone") != "" {
		customerPhone = r.FormValue("phone")
	}

	invoice := &model.Invoice{
		CustomerId:    customerInfo.Id,
		CustomerName:  customerName,
		DueDate:       time.Now().AddDate(0, 0, 30),
		CustomerEmail: customerEmail,
		CompanyName:   companyName,
		CustomerPhone: customerPhone,
		ItemList:      itemList,
	}

	return invoice, nil
}

func (h *InvoiceHandler) saveInvoiceWithItems(invoice *model.Invoice, itemList []model.ItemList) error {
	if h.repo == nil {
		log.Println("Error: repo is nil")
		return fmt.Errorf("repo is nil")
	}

	// Begin a new transaction
	log.Println("Starting a new transaction")
	tx, err := h.repo.BeginTransaction()
	if err != nil {
		log.Printf("Database error on beginning transaction: %v\n", err)
		return fmt.Errorf("database error on beginning transaction: %v", err)
	}

	// Attempt to add the new invoice
	log.Println("Attempting to add a new invoice")
	invoiceId, err := h.repo.AddNewInvoice(tx, *invoice)
	if err != nil {
		log.Printf("Database error on creating new invoice: %v\n", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Rollback error after failing to create new invoice: %v\n", rbErr)
			return fmt.Errorf("database error on creating new invoice: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("database error on creating new invoice: %v", err)
	}

	log.Printf("Successfully added new invoice with ID: %s\n", invoiceId)

	// Attempt to add items to the invoice
	for i, item := range itemList {
		item.InvoiceId = invoiceId
		log.Printf("Attempting to add item %d: %+v\n", i+1, item)
		if err := h.repo.AddNewItem(tx, item); err != nil {
			log.Printf("Database error on adding new item: %v\n", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Rollback error after failing to add new item: %v\n", rbErr)
				return fmt.Errorf("database error on adding new item: %v, rollback error: %v", err, rbErr)
			}
			return fmt.Errorf("database error on adding new item: %v", err)
		}
		log.Printf("Successfully added item %d\n", i+1)
	}

	// Attempt to commit the transaction
	log.Println("Attempting to commit the transaction")
	if err := tx.Commit(); err != nil {
		log.Printf("Database error on committing transaction: %v\n", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Rollback error after failing to commit transaction: %v\n", rbErr)
			return fmt.Errorf("database error on committing transaction: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("database error on committing transaction: %v", err)
	}

	log.Println("Transaction committed successfully")
	return nil
}

func prepareInvoiceData(invoice *model.Invoice, itemList []model.ItemList, customerAddress *model.Address, paymentStatusInt int) InvoiceData {
	return InvoiceData{
		Invoices: []model.Invoice{
			{
				InvoiceId:       invoice.InvoiceId,
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
}

// Add Items to Invoice

func (h *InvoiceHandler) AddItemsToInvoice(r *http.Request) ([]model.ItemList, error) {
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v\n", err)
		return nil, err
	}

	itemNames := r.Form["item"]
	quantities := r.Form["quantity"]
	unitPrices := r.Form["unitPrice"]
	subTotals := r.Form["subtotal"]
	taxes := r.Form["tax"]
	totals := r.Form["total"]

	// Log the lengths of the form arrays
	log.Printf("itemNames length: %d, quantities length: %d, unitPrices length: %d, subTotals length: %d, taxes length: %d, totals length: %d",
		len(itemNames), len(quantities), len(unitPrices), len(subTotals), len(taxes), len(totals))

	// Check if any of the slices are empty
	if len(itemNames) == 0 || len(quantities) == 0 || len(unitPrices) == 0 || len(subTotals) == 0 || len(taxes) == 0 || len(totals) == 0 {
		return nil, fmt.Errorf("form values are missing")
	}

	// Create a formatted string to log all values
	var formattedString string
	for i := 0; i < len(itemNames); i++ {
		formattedString += fmt.Sprintf(
			"Item %d: {item: %s, quantity: %s, unitPrice: %s, subTotal: %s, tax: %s, total: %s}\n",
			i+1, itemNames[i], quantities[i], unitPrices[i], subTotals[i], taxes[i], totals[i],
		)
	}
	log.Printf("Form values:\n%s", formattedString)

	var itemList []model.ItemList

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
		unitPrice, err := strconv.ParseFloat(unitPrices[i], 32)
		if err != nil {
			return nil, fmt.Errorf("invalid unit price for item %d", i+1)
		}

		// Convert and validate subtotal
		subTotal, err := strconv.ParseFloat(subTotals[i], 32)
		if err != nil {
			return nil, fmt.Errorf("invalid subtotal for item %d", i+1)
		}

		// Convert and validate tax
		tax, err := strconv.ParseFloat(taxes[i], 32)
		if err != nil {
			return nil, fmt.Errorf("invalid tax for item %d", i+1)
		}

		// Convert and validate total
		total, err := strconv.ParseFloat(totals[i], 32)
		if err != nil {
			return nil, fmt.Errorf("invalid total for item %d", i+1)
		}

		// Create an ItemList struct
		item := model.ItemList{
			Item:      itemNames[i],
			Quantity:  int32(quantity),
			UnitPrice: float32(unitPrice),
			Subtotal:  float32(subTotal),
			Tax:       float32(tax),
			Total:     float32(total),
		}

		// Add the item to the list
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

	unitPrice, err := strconv.ParseFloat(unitPriceStr, 32)
	if err != nil {
		log.Printf("Error parsing unit price: %v", err)
		http.Error(w, "<p>Error: Invalid unit price.</p>", http.StatusBadRequest)
		return
	}

	subtotal := float32(quantity) * float32(unitPrice)
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
