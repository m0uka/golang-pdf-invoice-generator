package main

import (
	_ "embed"
	"github.com/m0uka/golang-pdf-invoice-generator/invoicepdf"
	"log"
	"os"
)

func main() {
	buffer, err := invoicepdf.GenerateInvoice(invoicepdf.ExampleInvoice())
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("test.pdf", buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
