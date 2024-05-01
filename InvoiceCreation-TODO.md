
# Invoice Creation TODO List

## High Priority Tasks

### XSS Attack Prevention
- **Validate User Inputs**: Ensure all inputs are validated for type, length, format, and range before processing.
- **Encode Outputs**: Use HTML encoding on any outputs that could be manipulated through user input to prevent script injection.
- **Sanitize Inputs on Arrival**: Utilize libraries or built-in functions to sanitize inputs on server arrival.

### Address Handling for Customers
- **Check for Existing Address**: Verify if an existing customer has an address.

  ```go
  customer, err := repo.GetCustomerByID(customerID)
  if customer.Address == nil {
      // Prompt for address
  }
  ```
  
- **Address Update or Addition**: Update or add a new address if the provided invoice data includes address changes.
  ```go
  if addressChanged {
      err := repo.UpdateCustomerAddress(customerID, newAddress)
  }
  ```

### Managing Invoice Items
- Ensure invoice items parsed from the request are stored in the database.
  ```go
  itemsData := parseItems(r)
  for _, item := range itemsData {
      repo.AddItemToInvoice(invoiceID, item)
  }
  ```

## Medium Priority Tasks

### Payment Status Initialization
- Initialize and manage the payment status for each invoice.
  ```go
  status := parsePaymentStatus(r.FormValue("status"))
  updateInvoiceStatus(invoiceID, status)
  ```

### Content Security Policy (CSP)
- Implement CSP headers to control sources of scripts and other resources.
  ```go
  w.Header().Set("Content-Security-Policy", "default-src 'self';")
  ```

### Linking with Customer and Address Models
- Correctly link invoices with their customer records and include address details.
  ```go
  invoice.Customer = getCustomerWithAddress(customerID)
  ```

## Low Priority Tasks

### Consistency and Integrity
- **Database Transactions**: Use transactions for multi-step operations to ensure data integrity.
  ```go
  tx, err := db.Begin()
  if err != nil {
      return err
  }
  defer tx.Rollback()
  // Perform operations
  tx.Commit()
  ```
- **Error Handling**: Comprehensive handling of potential data errors.
  ```go
  if err != nil {
      log.Error("Failed operation:", err)
      handleErrorResponse(w, "Operation failed", http.StatusInternalServerError)
  }
  ```

### Example Code for Address Check and Update
- Demonstrate how to check and update a customer's address during invoice creation.
  ```go
  if customer.Address == nil {
      updateAddress(customerID, newAddress)
  }
  ```
