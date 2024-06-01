package handler

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	pdfGenUtils "github.com/MrAjMann/crm/generator"
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
	return &InvoiceHandler{repo: repo, tmpl: tmpl, customerHandler: customerHandler}
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

	if err := h.tmpl.ExecuteTemplate(w, "pdfModal", data); err != nil {
		log.Printf("Error executing modal template: %v\n", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
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

func getCustomerIdFromRequest(r *http.Request) (int, error) {
	customerIdStr := r.FormValue("customerId")
	if customerIdStr == "" {
		return 0, fmt.Errorf("customer ID is required")
	}
	return strconv.Atoi(customerIdStr)
}

func getPaymentStatusFromRequest(r *http.Request) (int, error) {
	paymentStatusStr := r.FormValue("paymentStatus")
	if paymentStatusStr == "" {
		paymentStatusStr = "0"
	}
	return strconv.Atoi(paymentStatusStr)
}

func createInvoiceFromRequest(r *http.Request, customerInfo *model.Customer, itemList []model.ItemList) (*model.Invoice, error) {
	dueDate, err := time.Parse("2006-01-02", r.FormValue("DueDate"))
	if err != nil {
		return nil, fmt.Errorf("invalid due date: %v", err)
	}

	invoice := &model.Invoice{
		CustomerId:    customerInfo.Id,
		CustomerName:  r.FormValue("customerName"),
		DueDate:       dueDate,
		CustomerEmail: r.FormValue("email"),
		CompanyName:   r.FormValue("companyName"),
		CustomerPhone: r.FormValue("phone"),
		ItemList:      itemList,
	}
	return invoice, nil
}

func (h *InvoiceHandler) saveInvoiceWithItems(invoice *model.Invoice, itemList []model.ItemList) error {
	log.Println("Attempting to save invoice with items")
	if h.repo == nil {
		log.Printf("Error: repo is nil")
		return fmt.Errorf("repo is nil")
	}

	tx, err := h.repo.BeginTransaction()
	if err != nil {
		log.Printf("Database error on beginning transaction: %v\n", err)
		return fmt.Errorf("database error on beginning transaction: %v", err)
	}

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
	invoice.InvoiceId = invoiceId

	for _, item := range itemList {
		item.InvoiceId = invoiceId // Set the invoice ID for each item as a string
		if err := h.repo.AddNewItem(tx, item); err != nil {
			log.Printf("Database error on adding new item: %v\n", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Rollback error after failing to add new item: %v\n", rbErr)
				return fmt.Errorf("database error on adding new item: %v, rollback error: %v", err, rbErr)
			}
			return fmt.Errorf("database error on adding new item: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Database error on committing transaction: %v\n", err)
		return fmt.Errorf("database error on committing transaction: %v", err)
	}

	log.Println("Successfully committed transaction")
	return nil
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

func (h *InvoiceHandler) GeneratePdf(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v\n", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	invoiceId := r.FormValue("invoiceId")
	if invoiceId == "" {
		log.Println("Invoice ID is missing")
		http.Error(w, "Invoice ID is required", http.StatusBadRequest)
		return
	}

	// Fetch the invoice details
	invoice, err := h.repo.GetInvoiceById(invoiceId)
	if err != nil {
		log.Printf("Error fetching invoice: %v\n", err)
		http.Error(w, "Error fetching invoice", http.StatusInternalServerError)
		return
	}

	// Generate the PDF
	err = pdfGenUtils.CreatePdf(invoice)
	if err != nil {
		log.Printf("Error generating PDF: %v\n", err)
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PDF generated successfully"))
}
