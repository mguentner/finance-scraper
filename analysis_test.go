package finance_test

import (
	"github.com/hill-daniel/finance-scraper"
	"github.com/hill-daniel/finance-scraper/test"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
	"testing"
)

func Test_big_number_suffix_format(t *testing.T) {
	testCases := []struct {
		name     string
		input    float64
		expected string
	}{
		{name: "should format under - a Trillion", input: -5608991029248, expected: "-5.61 T"},
		{name: "should format under - a Billion", input: -28991029248, expected: "-28.99 B"},
		{name: "should format under - a Million", input: -47123456, expected: "-47.12 M"},
		{name: "should format under -10,000 with decimal", input: -12345.78, expected: "-12.35 k"},
		{name: "should format under -10,000", input: -10000.0, expected: "-10.00 k"},
		{name: "should format -1,000", input: -1000.0, expected: "-1.00 k"},
		{name: "should format over -1,000", input: -999.0, expected: "-999.00"},
		{name: "should format under 0", input: -34.24, expected: "-34.24"},
		{name: "should format under 1000", input: 999.0, expected: "999.00"},
		{name: "should format over 1,000", input: 1000.0, expected: "1.00 k"},
		{name: "should format over 10,000", input: 10000.0, expected: "10.00 k"},
		{name: "should format over 10,000 with decimal", input: 12345.78, expected: "12.35 k"},
		{name: "should format over 100,000", input: 110592, expected: "110.59 k"},
		{name: "should format over a Million", input: 47123456, expected: "47.12 M"},
		{name: "should format over a Billion", input: 28991029248, expected: "28.99 B"},
		{name: "should format over a Trillion", input: 5608991029248, expected: "5.61 T"},
	}

	p := message.NewPrinter(language.English)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := finance.FmtWithSuffix(tc.input, p)

			test.Equals(t, tc.expected, actual)
		})
	}
}

func TestAsanaAnalyzer_Analyze(t *testing.T) {
	testCases := []struct {
		name     string
		input    finance.Quote
		expected string
	}{
		{name: "should be green check if PERatio is above 0", input: finance.Quote{
			PERatio: 10,
		}, expected: "PE Ratio 10.00 -> ✅"},
		{name: "should not ok if PERatio is 0", input: finance.Quote{
			PERatio: 0,
		}, expected: "PE Ratio 0.00 -> ⭕"},
		{name: "should green check if enterprise value is smaller than market cap", input: finance.Quote{
			MarketCap:       100,
			EnterpriseValue: 99,
		}, expected: "Enterprise Value 99.00 < Market Cap 100.00 -> ✅"},
		{name: "should not be ok if enterprise value is bigger than market cap", input: finance.Quote{
			MarketCap:       99,
			EnterpriseValue: 100,
		}, expected: "Enterprise Value 100.00 > Market Cap 99.00 -> ⭕"},
		{name: "should green check if return on equity is over 20", input: finance.Quote{
			ReturnOnEquity: 20.1,
		}, expected: "Return on Equity 20.10 -> ✅"},
		{name: "should ok if return on equity is below 20 but above 0", input: finance.Quote{
			ReturnOnEquity: 5.67,
		}, expected: "Return on Equity 5.67 -> ☑️"},
		{name: "should not be ok if return on equity is below 0", input: finance.Quote{
			ReturnOnEquity: -12.84,
		}, expected: "Return on Equity -12.84 -> ⭕"},
		{name: "should be green check if insider ownership is above 5 and below 40", input: finance.Quote{
			InsiderOwnership: 33.45,
		}, expected: "Insider Ownership 33.45 -> ✅"},
		{name: "should not be ok if insider ownership is below 5", input: finance.Quote{
			InsiderOwnership: 4.99,
		}, expected: "Insider Ownership 4.99 -> ⭕"},
		{name: "should not be ok if insider ownership is above 40", input: finance.Quote{
			InsiderOwnership: 40.99,
		}, expected: "Insider Ownership 40.99 -> ⭕"},
		{name: "should be all green if all criteria are met", input: finance.Quote{
			PERatio:          10,
			EnterpriseValue:  99,
			MarketCap:        100,
			ReturnOnEquity:   25,
			InsiderOwnership: 33,
		}, expected: "PE Ratio 10.00 -> ✅\nEnterprise Value 99.00 < Market Cap 100.00 -> ✅\nReturn on Equity 25.00 -> ✅\nInsider Ownership 33.00 -> ✅"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyzer := finance.AsanaAnalyzer{}

			result, err := analyzer.Analyze(tc.input)

			test.Ok(t, err)
			test.Assert(t, strings.Contains(result, tc.expected), "expected %s to contain %s", result, tc.expected)
		})
	}
}
