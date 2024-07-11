package invoicepdf

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/pkg/errors"
	"github.com/signintech/gopdf"
)

//go:embed "Inter/Inter Variable/Inter.ttf"
var interFont []byte

//go:embed "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf"
var interBoldFont []byte

func GenerateInvoice(invoice *Invoice) (*bytes.Buffer, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	pdf.SetMargins(40, 40, 40, 40)
	pdf.AddPage()
	err := pdf.AddTTFFontData("Inter", interFont)
	if err != nil {
		return nil, err
	}

	err = pdf.AddTTFFontData("Inter-Bold", interBoldFont)
	if err != nil {
		return nil, err
	}

	// This is nasty, but so is the rest of the library, soooo
	from := fmt.Sprintf("%s\n%s\n%s\n%s, %s %s\n%s\n%s\n \n%s\n%s", invoice.From.Name, invoice.From.AddressLine1, invoice.From.AddressLine2, invoice.From.City, invoice.From.State, invoice.From.PostalCode, invoice.From.Country, invoice.From.TaxID, invoice.From.Email, invoice.From.Website)
	to := fmt.Sprintf("%s\n%s\n%s\n%s, %s %s\n%s\n%s\n \n%s\n%s", invoice.To.Name, invoice.To.AddressLine1, invoice.To.AddressLine2, invoice.To.City, invoice.To.State, invoice.To.PostalCode, invoice.From.Country, invoice.To.TaxID, invoice.To.Email, invoice.To.Website)

	if writeLogo(&pdf, invoice.LogoUrl, from, invoice.HeaderNote) != nil {
		return nil, errors.Wrap(err, "failed to call writeLogo")
	}

	writeTitle(&pdf, invoice.Title, invoice.Id, invoice.Date, invoice.Number)
	writeBillTo(&pdf, to)
	writeHeaderRow(&pdf, invoice.AmountOnly)
	subtotal := 0.0
	for i := range invoice.Items {
		q := 1
		if len(invoice.Quantities) > i {
			q = invoice.Quantities[i]
		}

		r := 0.0
		if len(invoice.Rates) > i {
			r = invoice.Rates[i]
		}

		writeRow(&pdf, invoice.Items[i], q, r, invoice.Currency, invoice.AmountOnly)
		subtotal += float64(q) * r
	}
	if invoice.Note != "" {
		writeNotes(&pdf, invoice.Note)
	}
	writeTotals(&pdf, subtotal, invoice.Tax, invoice.Discount, invoice.Currency)
	if invoice.Due != "" {
		writeDueDate(&pdf, invoice.Due)
	}
	writeFooter(&pdf, invoice.Id)

	var buf bytes.Buffer
	_, err = pdf.WriteTo(&buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call pdf.WriteTo")
	}

	return &buf, nil
}
