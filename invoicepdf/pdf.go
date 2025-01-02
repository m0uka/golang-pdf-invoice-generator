package invoicepdf

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"image"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

const (
	quantityColumnOffset = 360
	rateColumnOffset     = 405
	amountColumnOffset   = 480
)

const (
	subtotalLabel = "Subtotal"
	discountLabel = "Discount"
	feeLabel      = "Fees"
	taxLabel      = "Tax"
	totalLabel    = "Total"
)

func writeLogo(pdf *gopdf.GoPdf, logo string, from string, headerNote string) error {
	if logo != "" {

		resp, err := http.Get(logo)
		if err != nil {
			return errors.Wrap(err, "failed to get logo from URL")
		}
		defer resp.Body.Close()

		var buf bytes.Buffer
		tee := io.TeeReader(resp.Body, &buf)

		img, _, err := image.Decode(tee)
		if err != nil {
			return errors.Wrap(err, "failed to decode image from URL")
		}

		imgConfig, _, err := image.DecodeConfig(&buf)
		if err != nil {
			return errors.Wrap(err, "failed to decode image as config from URL")
		}

		scaledWidth := 100.0
		scaledHeight := float64(imgConfig.Height) * scaledWidth / float64(imgConfig.Width)
		_ = pdf.ImageFrom(img, pdf.GetX(), pdf.GetY(), &gopdf.Rect{W: scaledWidth, H: scaledHeight})
		pdf.Br(scaledHeight + 24)
	}

	formattedFrom := strings.ReplaceAll(from, `\n`, "\n")
	fromLines := strings.Split(formattedFrom, "\n")

	pdf.SetTextColor(75, 75, 75)
	_ = pdf.SetFont("Inter", "", 9)
	_ = pdf.Cell(nil, "INVOICE FROM")
	pdf.Br(18)

	pdf.SetTextColor(55, 55, 55)

	for i := 0; i < len(fromLines); i++ {
		line := fromLines[i]
		if line == "" {
			continue
		}

		if i == 0 {
			_ = pdf.SetFont("Inter", "", 12)
			_ = pdf.Cell(nil, line)
			pdf.Br(18)
		} else {
			_ = pdf.SetFont("Inter", "", 10)
			_ = pdf.Cell(nil, line)
			pdf.Br(15)
		}
	}

	if headerNote != "" {
		pdf.Br(4)
		pdf.SetTextColor(100, 100, 100)
		_ = pdf.SetFont("Inter", "", 8)
		_ = pdf.Cell(nil, headerNote)
		pdf.Br(4)
		pdf.SetTextColor(55, 55, 55)
	}

	pdf.Br(21)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX(), pdf.GetY(), 260, pdf.GetY())
	pdf.Br(36)

	return nil
}

func writeTitle(pdf *gopdf.GoPdf, title, id, date string, invNumber string) {
	_ = pdf.SetFont("Inter-Bold", "", 24)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.Cell(nil, title)
	pdf.Br(36)
	_ = pdf.SetFont("Inter", "", 12)
	pdf.SetTextColor(100, 100, 100)
	_ = pdf.Cell(nil, id)
	pdf.SetTextColor(150, 150, 150)
	_ = pdf.Cell(nil, "  ·  ")
	if invNumber != "" {
		pdf.SetTextColor(100, 100, 100)
		_ = pdf.Cell(nil, invNumber)

		pdf.SetTextColor(150, 150, 150)
		_ = pdf.Cell(nil, "  ·  ")
	}
	pdf.SetTextColor(100, 100, 100)
	_ = pdf.Cell(nil, date)
	pdf.Br(48)
}

func writeDueDate(pdf *gopdf.GoPdf, due string) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, "Due Date")
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(11)
	pdf.SetX(amountColumnOffset - 15)
	_ = pdf.Cell(nil, due)
	pdf.Br(12)
}

func writeBillTo(pdf *gopdf.GoPdf, to string) {
	pdf.SetTextColor(75, 75, 75)
	_ = pdf.SetFont("Inter", "", 9)
	_ = pdf.Cell(nil, "INVOICE TO")
	pdf.Br(18)
	pdf.SetTextColor(75, 75, 75)

	formattedTo := strings.ReplaceAll(to, `\n`, "\n")
	toLines := strings.Split(formattedTo, "\n")

	for i := 0; i < len(toLines); i++ {
		line := toLines[i]
		if line == "" {
			continue
		}

		if i == 0 {
			_ = pdf.SetFont("Inter", "", 15)
			_ = pdf.Cell(nil, line)
			pdf.Br(20)
		} else {
			_ = pdf.SetFont("Inter", "", 10)
			_ = pdf.Cell(nil, line)
			pdf.Br(15)
		}
	}
	pdf.Br(46)
}

func writeHeaderRow(pdf *gopdf.GoPdf, amountOnly bool) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "ITEM")
	if !amountOnly {
		pdf.SetX(quantityColumnOffset)
		_ = pdf.Cell(nil, "QTY")
		pdf.SetX(rateColumnOffset)
		_ = pdf.Cell(nil, "RATE")
	}
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, "AMOUNT")
	pdf.Br(24)
}

func writeNotes(pdf *gopdf.GoPdf, notes string) {
	pdf.SetY(600)

	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, "NOTES")
	pdf.Br(18)
	_ = pdf.SetFont("Inter", "", 12)
	pdf.SetTextColor(0, 0, 0)

	formattedNotes := strings.ReplaceAll(notes, `\n`, "\n")
	notesLines := strings.Split(formattedNotes, "\n")

	for i := 0; i < len(notesLines); i++ {
		_ = pdf.Cell(nil, notesLines[i])
		pdf.Br(15)
	}

	pdf.Br(48)
}
func writeFooter(pdf *gopdf.GoPdf, id string) {
	pdf.SetY(800)

	_ = pdf.SetFont("Inter", "", 10)
	pdf.SetTextColor(55, 55, 55)
	_ = pdf.Cell(nil, id)
	pdf.SetStrokeColor(225, 225, 225)
	pdf.Line(pdf.GetX()+10, pdf.GetY()+6, 550, pdf.GetY()+6)
	pdf.Br(48)
}

func writeRow(pdf *gopdf.GoPdf, item string, quantity int, rate float64, currency string, amountOnly bool) {
	_ = pdf.SetFont("Inter", "", 11)
	pdf.SetTextColor(0, 0, 0)

	total := float64(quantity) * rate
	amount := strconv.FormatFloat(total, 'f', 2, 64)

	_ = pdf.Cell(nil, item)
	if !amountOnly {
		pdf.SetX(quantityColumnOffset)
		_ = pdf.Cell(nil, strconv.Itoa(quantity))
		pdf.SetX(rateColumnOffset)
		_ = pdf.Cell(nil, currencySymbols[currency]+strconv.FormatFloat(rate, 'f', 2, 64))
	}
	pdf.SetX(amountColumnOffset)
	_ = pdf.Cell(nil, currencySymbols[currency]+amount)
	pdf.Br(24)
}

func writeTotals(pdf *gopdf.GoPdf, subtotal float64, tax float64, discount float64, fees float64, currency string) {
	pdf.SetY(600)

	writeTotal(pdf, subtotalLabel, subtotal, currency)
	writeTotal(pdf, taxLabel, tax, currency)
	if discount > 0 {
		writeTotal(pdf, discountLabel, discount, currency)
	}
	if fees > 0 {
		writeTotal(pdf, feeLabel, fees, currency)
	}
	writeTotal(pdf, totalLabel, subtotal+tax-discount-fees, currency)
}

func writeTotal(pdf *gopdf.GoPdf, label string, total float64, currency string) {
	_ = pdf.SetFont("Inter", "", 9)
	pdf.SetTextColor(75, 75, 75)
	pdf.SetX(rateColumnOffset)
	_ = pdf.Cell(nil, label)
	pdf.SetTextColor(0, 0, 0)
	_ = pdf.SetFontSize(12)
	pdf.SetX(amountColumnOffset - 15)
	if label == totalLabel {
		_ = pdf.SetFont("Inter-Bold", "", 11.5)
	}
	_ = pdf.Cell(nil, currencySymbols[currency]+strconv.FormatFloat(total, 'f', 2, 64))
	pdf.Br(24)
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
