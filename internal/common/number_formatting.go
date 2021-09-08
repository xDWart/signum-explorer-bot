package common

import (
	"github.com/xDWart/signum-explorer-bot/api/signumapi"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p = message.NewPrinter(language.English)

func FormatNumber(number float64, decimals int) string {
	switch decimals {
	case 0:
		return p.Sprintf("%.f", number)
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
		return p.Sprintf("%v", number)
	}
}

func FormatNQT(number uint64) string {
	return p.Sprintf("%.2f", float64(number)/1e8)
}

func ConvertFeeNQT(fee signumapi.FeeType) float64 {
	return float64(fee) / 1e8
}
