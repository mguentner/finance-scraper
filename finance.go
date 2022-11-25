package finance

// ScrapeResult holds the result data from a scraping run.
type ScrapeResult struct {
	PERatio           string
	EnterpriseValue   string
	MarketCap         string
	ReturnOnEquity    string
	InsiderOwnership  string
	OutstandingShares string
}

// Quote is the mapped result of ScrapeResult, for calculation purposes.
type Quote struct {
	PERatio           float64
	EnterpriseValue   float64
	MarketCap         float64
	ReturnOnEquity    float64
	InsiderOwnership  float64
	OutstandingShares float64
}

// Scraper scrapes a data source.
type Scraper interface {
	Scrape(symbol string) (ScrapeResult, error)
}

// Analyzer analyzes given quote and returns a result.
type Analyzer interface {
	Analyze(quote Quote) (string, error)
}
