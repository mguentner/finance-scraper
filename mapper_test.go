package finance_test

import (
	"testing"

	"github.com/mguentner/finance-scraper"
	"github.com/mguentner/finance-scraper/test"
)

func Test(t *testing.T) {
	testCases := []struct {
		name          string
		inputScrape   finance.ScrapeResult
		expectedQuote finance.Quote
	}{
		{name: "basic", inputScrape: finance.ScrapeResult{
			PERatio:          "1,118.09",
			EnterpriseValue:  "555.87B",
			MarketCap:        "555.24B",
			ReturnOnEquity:   "5.59%",
			InsiderOwnership: "20.02%",
		}, expectedQuote: finance.Quote{
			PERatio:          1118.09,
			EnterpriseValue:  555870000000,
			MarketCap:        555240000000,
			ReturnOnEquity:   5.59,
			InsiderOwnership: 20.02,
		}},
		{name: "trillion", inputScrape: finance.ScrapeResult{
			PERatio:          "1,118.09",
			EnterpriseValue:  "2.17T",
			MarketCap:        "555.24B",
			ReturnOnEquity:   "5.59%",
			InsiderOwnership: "20.02%",
		}, expectedQuote: finance.Quote{
			PERatio:          1118.09,
			EnterpriseValue:  2170000000000,
			MarketCap:        555240000000,
			ReturnOnEquity:   5.59,
			InsiderOwnership: 20.02,
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mapper := finance.YahooMapper{}

			quote, err := mapper.Map(tc.inputScrape)

			test.Ok(t, err)
			test.Equals(t, tc.expectedQuote, quote)
		})
	}
}
