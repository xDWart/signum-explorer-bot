package common

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p = message.NewPrinter(language.English)

func FormatNumber(number float64, decimals byte) string {
	switch decimals {
	case 1:
		return p.Sprintf("%.1f", number)
	case 2:
		return p.Sprintf("%.2f", number)
	case 3:
		return p.Sprintf("%.3f", number)
	case 4:
		return p.Sprintf("%.4f", number)
	case 5:
		return p.Sprintf("%.5f", number)
	case 6:
		return p.Sprintf("%.6f", number)
	case 7:
		return p.Sprintf("%.7f", number)
	case 8:
		return p.Sprintf("%.8f", number)
	default:
		return p.Sprintf("%.f", number)
	}
}
