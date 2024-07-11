package invoicepdf

import "time"

type Invoice struct {
	Id     string `json:"id" yaml:"id"`
	Number string `json:"number"`
	Title  string `json:"title" yaml:"title"`

	LogoUrl string         `json:"logo" yaml:"logo"`
	From    InvoiceCompany `json:"from" yaml:"from"`
	To      InvoiceCompany `json:"to" yaml:"to"`
	Date    string         `json:"date" yaml:"date"`
	Due     string         `json:"due" yaml:"due"`

	Items      []string  `json:"items" yaml:"items"`
	Quantities []int     `json:"quantities" yaml:"quantities"`
	Rates      []float64 `json:"rates" yaml:"rates"`
	AmountOnly bool      `json:"amount_only"`

	Tax      float64 `json:"tax" yaml:"tax"`
	Discount float64 `json:"discount" yaml:"discount"`
	Currency string  `json:"currency" yaml:"currency"`

	Note       string `json:"note" yaml:"note"`
	HeaderNote string `json:"header_note"`
}

type InvoiceCompany struct {
	Hide         bool   `json:"hide"`
	Name         string `json:"name"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	Country      string `json:"country"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	TaxID        string `json:"tax_id"`
	Email        string `json:"email"`
	Website      string `json:"website"`
}

func ExampleInvoice() *Invoice {
	return &Invoice{
		Id:         "pn-payout-4we8a9ew6",
		Number:     "1r56war4ea-1",
		Title:      "INVOICE",
		Rates:      []float64{2500},
		Quantities: []int{1},
		Items:      []string{"Digital Services licensed to PayNow for resale"},
		LogoUrl:    "https://cdn.paynow.gg/logo/full/logotype-color.png",
		From: InvoiceCompany{
			Name:         "Fake Company",
			AddressLine1: "Somewhere over the rainbow",
			AddressLine2: "",
			City:         "Hoboken",
			Country:      "United States",
			State:        "NJ",
			PostalCode:   "12312",
			TaxID:        "",
		},
		To: InvoiceCompany{
			Name:         "PayNow Services, Inc.",
			AddressLine1: "123 Fake Street",
			AddressLine2: "",
			City:         "New York",
			Country:      "United States",
			State:        "NY",
			PostalCode:   "10001",
			TaxID:        "",
			Email:        "support@paynow.gg",
			Website:      "www.paynow.gg",
		},
		Date:       time.Now().Format("Jan 02, 2006 MST"),
		Tax:        0,
		Discount:   0,
		AmountOnly: true,
		Currency:   "USD",
		Note:       "NO VAT - REVERSE CHARGE IF APPLICABLE",
		HeaderNote: "* You can change your billing details in payout provider settings",
	}
}
