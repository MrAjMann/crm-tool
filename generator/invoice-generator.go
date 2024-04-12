package main

import (
	"log"

	"github.com/go-pdf/fpdf"
)

func main() {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetCompression(false) // Consider enabling for production

	// // Define a template - 210 x 297 mm
	template := pdf.CreateTemplate(func(tpl *fpdf.Tpl) {
		tpl.Image("./logo.png", 6, 10, 15, 0, false, "", 0, "")
		tpl.SetFont("Arial", "B", 16)
		tpl.Text(165, 280, "A&R TECH")
		tpl.SetFont("Arial", "B", 12)
		tpl.Text(142, 287, "PC SUPPORT ON THE GO")

	})

	template2 := pdf.CreateTemplate(func(tpl *fpdf.Tpl) {
		tpl.UseTemplate(template)
		subtemplate := tpl.CreateTemplate(func(tpl2 *fpdf.Tpl) {
			tpl2.SetFont("Arial", "B", 16)
			tpl2.Text(40, 200, "Subtemplate says hello")
			tpl2.SetDrawColor(0, 200, 100)
			tpl2.SetLineWidth(2.5)
			tpl2.Line(102, 92, 112, 102)
		})
		tpl.UseTemplate(subtemplate)
	})

	// Use the template
	pdf.AddPage()
	pdf.UseTemplate(template2)

	// Output the PDF
	err := pdf.OutputFileAndClose("hello1.pdf")
	if err != nil {
		log.Fatalf("Error creating PDF: %v", err) // More descriptive error handling
	}
}
