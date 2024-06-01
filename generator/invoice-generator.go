package pdfGenUtils

import (
	"fmt"
	"log"

	"time"

	"github.com/MrAjMann/crm/internal/model"
	"github.com/go-pdf/fpdf"
)

func CreatePdf(invoice *model.Invoice) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	print("Creating PDF...")
	pdf.SetMargins(10, 10, 10)
	log.Print(invoice)
	header(pdf)
	footer(pdf)
	clientInfo(pdf, invoice)

	// Output the PDF
	err := pdf.OutputFileAndClose("invoice_" + invoice.InvoiceId + ".pdf")
	if err != nil {
		log.Fatalf("Error creating PDF: %v", err) // More descriptive error handling
	}
	return nil
}

func header(pdf *fpdf.Fpdf) {
	pdf.SetHeaderFunc(func() {
		// Add logo to the header
		pdf.CellFormat(7, 7, "A&R TECH", "", 0, "R", false, 0, "")

		// Add "INVOICE" at the right corner
		pdf.SetFont("Arial", "B", 32)
		pdf.SetXY(170, 10) // Adjust X and Y positions as needed
		pdf.CellFormat(30, 10, "INVOICE", "", 0, "R", false, 0, "")

		// Add "GENERATED at May 12, 2024" below "EUR STATEMENT"
		pdf.SetXY(170, 15) // Adjust Y position to place the text below "INVOICE"
		pdf.SetFont("Arial", "I", 6)
		text := fmt.Sprintf("Generated at %v %v, %v", time.Now().Month(), time.Now().Day(), time.Now().Year())
		pdf.CellFormat(30, 10, text, "", 0, "R", false, 0, "")

		// Add Company Name
		pdf.SetXY(170, 19) // Adjust Y position to place the text below "INVOICE"
		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(30, 10, "Your Company Ltd", "", 1, "R", false, 0, "")

	})
}

func footer(pdf *fpdf.Fpdf) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-35)

		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(20, 20, "", "", 0, "L", false, 0, "") // empty cell to simulate the image inserted
		pdf.CellFormat(30, 6.5, "Report lost/stolen card", "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", "I", 4)
		pdf.CellFormat(129, 6.5, "Revolut Bank UAB is a credit institution...", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6.5, "", "", 1, "L", false, 0, "") //
		pdf.SetY(pdf.GetY() - 4)

		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(20, 6.5, "", "", 0, "L", false, 0, "") // empty cell to simulate the image inserted
		pdf.CellFormat(30, 6.5, "+44 20 3322 2222", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "I", 4)
		pdf.CellFormat(129, 6.5, "Vilnius, Lithuania. Licensed and regulated by...", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6.5, "", "", 1, "L", false, 0, "") //
		pdf.SetY(pdf.GetY() - 4)

		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(20, 6.5, "", "", 0, "L", false, 0, "") // empty cell to simulate the image inserted
		pdf.CellFormat(30, 6.5, "Get help in the app", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "I", 4)
		pdf.CellFormat(129, 6.5, "More information available at www.iidraudimas.lt", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6.5, "", "", 1, "L", false, 0, "") //
		pdf.SetY(pdf.GetY() - 4)

		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(20, 6.5, "", "", 0, "L", false, 0, "") // empty cell to simulate the image inserted
		pdf.CellFormat(30, 6.5, "Scan the QR code", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "I", 4)
		pdf.CellFormat(129, 6.5, "Contact us via in-app chat.", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6.5, "", "", 1, "L", false, 0, "") //

		pdf.SetFont("Arial", "B", 8)
		pdf.SetY(pdf.GetY() + 2)

		tr := pdf.UnicodeTranslatorFromDescriptor("")
		pdf.CellFormat(95, 10, fmt.Sprintf("%v Your Company Ltd", tr("©")), "", 0, "L", false, 0, "")
		pdf.CellFormat(95, 10, fmt.Sprintf("Page %d of NP", pdf.PageNo()), "", 0, "R", false, 0, "")
	})
}

func clientInfo(pdf *fpdf.Fpdf, invoice *model.Invoice) {
	// Set font
	pdf.SetFont("Arial", "B", 20)

	// name of the client
	//CellFormat(w, h, txtStr, borderStr, ln, alignStr, fill bool, link int, linkStr string)
	pdf.SetXY(10, 50)
	pdf.CellFormat(189, 20, invoice.CustomerName, "", 1, "L", false, 0, "")

	pdf.SetY(pdf.GetY() - 4)

	address := fmt.Sprintf("%s %s %s %s %s", invoice.CustomerAddress.UnitNumber, invoice.CustomerAddress.StreetNumber, invoice.CustomerAddress.StreetName, invoice.CustomerAddress.City, invoice.CustomerAddress.Postcode)
	// address of the client
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(100, 9, address, "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 9, "", "", 0, "L", false, 0, "")

	// Add a small vertical space between IBAN and BIC
	pdf.SetY(pdf.GetY() - 4)

	pdf.CellFormat(100, 9, invoice.CustomerEmail, "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 9, "", "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 9, "BIC", "", 0, "L", false, 0, "")
}
