package main

import (
	_ "embed"
	"github.com/maaslalani/invoice/invoicepdf"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "invoice",
	Short: "Invoice generates invoices from the command line.",
	Long:  `Invoice generates invoices from the command line.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return invoicepdf.GenerateInvoice(invoicepdf.DefaultInvoice())
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
