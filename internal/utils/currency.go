package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatCurrency(amount float64) string {
	p := message.NewPrinter(language.Indonesian) // Use Indonesian locale for formatting
	return p.Sprintf("%.0f", amount)             // Format without decimals
}
