package payslips

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatRupiah(value float64) string {
    p := message.NewPrinter(language.Indonesian)
    rounded := float64(int(value*100+0.5)) / 100
    return p.Sprintf("Rp.%.2f", rounded)
}