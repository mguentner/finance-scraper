# Yahoo Finance HTML Scraper
Scrape data for stock symbol and analyze result.

## Usage
* ```go build cmd/finance-scraper/finance-scraper.go```
* ```./finance-scraper -analyzer=asana -symbol=NCNO```
* symbol beeing stock symbol and asana is the only analzyer supported right now

example output: 
```
PE Ratio 0.00 -> ⭕
Enterprise Value 6.28 B < Market Cap 6.50 B -> ✅
Return on Equity -13.74 -> ⭕
Insider Ownership 2.17 -> ⭕
```
