package finance

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"strings"
)

const bigNumberSuffix = "kMBTPE"

// CreateAnalyzer returns an Analyzer for given string. DefaultAnalyzer just prints the result.
func CreateAnalyzer(analyzer string) Analyzer {
	switch analyzer {
	case "asana":
		return &AsanaAnalyzer{}
	default:
		return &DefaultAnalyzer{}
	}
}

// AsanaAnalyzer analyses quotes and returns the result in copy paste friendly asana ticket format.
type AsanaAnalyzer struct {
}

// Analyze analyzes quotes with following criteria:
// P/E Ratio -> if its blank it is not yet profitable. Above 0 -> Green check, otherwise not ok.
// Enterprise Value < Market Cap = More Cash than Debt -> If so green check, otherwise not ok.
// Return on Equity -> 20%+ is great, green check. Above 0 ok, below 0 not ok.
// Held by insiders -> 5 to 40% green check, otherwise not ok.
func (aa *AsanaAnalyzer) Analyze(quote Quote) (string, error) {
	builder := strings.Builder{}
	p := message.NewPrinter(language.English)

	isProfitable := quote.PERatio > 0.0
	if isProfitable {
		builder.WriteString(fmt.Sprintf("PE Ratio %s -> ✅\n", FmtWithSuffix(quote.PERatio, p)))
	} else {
		builder.WriteString(fmt.Sprintf("PE Ratio %s -> ⭕\n", FmtWithSuffix(quote.PERatio, p)))
	}

	moreCashThanDebt := quote.EnterpriseValue < quote.MarketCap
	if moreCashThanDebt {
		builder.WriteString(fmt.Sprintf("Enterprise Value %s < Market Cap %s -> ✅\n", FmtWithSuffix(quote.EnterpriseValue, p), FmtWithSuffix(quote.MarketCap, p)))
	} else {
		builder.WriteString(fmt.Sprintf("Enterprise Value %s > Market Cap %s -> ⭕\n", FmtWithSuffix(quote.EnterpriseValue, p), FmtWithSuffix(quote.MarketCap, p)))
	}

	if quote.ReturnOnEquity > 20.0 {
		builder.WriteString(fmt.Sprintf("Return on Equity %s -> ✅\n", FmtWithSuffix(quote.ReturnOnEquity, p)))
	} else if quote.ReturnOnEquity > 0.0 {
		builder.WriteString(fmt.Sprintf("Return on Equity %s -> ☑️\n", FmtWithSuffix(quote.ReturnOnEquity, p)))
	} else {
		builder.WriteString(fmt.Sprintf("Return on Equity %s -> ⭕\n", FmtWithSuffix(quote.ReturnOnEquity, p)))
	}

	if quote.InsiderOwnership > 5.0 && quote.InsiderOwnership < 40.0 {
		builder.WriteString(fmt.Sprintf("Insider Ownership %s -> ✅\n", FmtWithSuffix(quote.InsiderOwnership, p)))
	} else {
		builder.WriteString(fmt.Sprintf("Insider Ownership %s -> ⭕\n", FmtWithSuffix(quote.InsiderOwnership, p)))
	}
	return builder.String(), nil
}

// DefaultAnalyzer just prints given quote.
type DefaultAnalyzer struct {
}

// Analyze just formats numbers and prints the quote.
func (da *DefaultAnalyzer) Analyze(quote Quote) (string, error) {
	builder := strings.Builder{}
	p := message.NewPrinter(language.English)
	builder.WriteString(fmt.Sprintf("PE Ratio %s\n", FmtWithSuffix(quote.PERatio, p)))
	builder.WriteString(fmt.Sprintf("Enterprise Value %s\n", FmtWithSuffix(quote.EnterpriseValue, p)))
	builder.WriteString(fmt.Sprintf("Market Cap %s\n", FmtWithSuffix(quote.MarketCap, p)))
	builder.WriteString(fmt.Sprintf("Return on Equity %s\n", FmtWithSuffix(quote.ReturnOnEquity, p)))
	builder.WriteString(fmt.Sprintf("Insider Ownership %s\n", FmtWithSuffix(quote.InsiderOwnership, p)))
	return builder.String(), nil
}

// FmtWithSuffix format big numbers to the following:
// Thousand,             1 000 -> 1 k
// Million,          1 000 000 -> 1 M
// Billion,      1 000 000 000 -> 1 B
// Trillion, 1 000 000 000 000 -> 1 T
func FmtWithSuffix(value float64, p *message.Printer) string {
	if value > -1000.0 && value < 1000.0 {
		return p.Sprintf("%.2f", value)
	}
	prefix := 1.0
	if value < 0 {
		prefix = -1.0
		value = value * prefix
	}
	exp := int(math.Log(value) / math.Log(1000))
	return p.Sprintf("%.2f %c", (value/math.Pow(1000, float64(exp)))*prefix, bigNumberSuffix[exp-1])
}
