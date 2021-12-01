package colly

import (
	"container/list"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/hill-daniel/finance-scraper"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	scrapeDelay         = 5 * time.Second
	domainSchema        = "finance.yahoo.com/*"
	statisticsURLFormat = "https://finance.yahoo.com/quote/%s/key-statistics?p=%s"
	summaryURLFormat    = "https://finance.yahoo.com/quote/%s?p=%s"
	keyPERatio          = "peRatio"
	keyMarketCap        = "marketCap"
	keyEnterpriseValue  = "enterpriseValue"
	keyReturnOnEquity   = "returnOnEquity"
	keyInsiderOwnership = "insiderOwnership"
)

// Collector is a Web Scraper.
type Collector struct {
	*colly.Collector
}

type scrapeValue struct {
	key   string
	value string
}

// NewYahooCollector creates a new Collector.
func NewYahooCollector() (*Collector, error) {
	c := colly.NewCollector()
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  domainSchema,
		RandomDelay: scrapeDelay,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create collector")
	}
	return &Collector{Collector: c}, nil
}

// Scrape scrapes data for given symbol.
func (c *Collector) Scrape(symbol string) (finance.ScrapeResult, error) {
	if len(strings.TrimSpace(symbol)) == 0 {
		return finance.ScrapeResult{}, errors.New("failed to scrape. Symbol cannot be empty")
	}
	scrapeStack := list.New()
	symbol = strings.ToUpper(symbol)

	c.OnRequest(func(r *colly.Request) {
		url := r.URL.String()
		log.Infof("Visiting %s", url)
	})

	c.OnHTML("td[data-test]", func(e *colly.HTMLElement) {
		attr := e.Attr("data-test")
		if attr == "PE_RATIO-value" {
			text := e.DOM.Find("span").Text()
			scrapeStack.PushFront(&scrapeValue{
				key:   keyPERatio,
				value: text,
			})
		}
	})

	url := fmt.Sprintf(summaryURLFormat, symbol, symbol)
	err := c.Visit(url)
	if err != nil {
		log.Errorf("failed to visit url %s: %v", url, err)
	}

	c.OnHTML("tr.fi-row", func(e *colly.HTMLElement) {
		nextIsValue := false
		e.ForEach("td", func(i int, element *colly.HTMLElement) {
			if nextIsValue {
				nextIsValue = false
				scrapeElement := scrapeStack.Back()
				scrapeValue := scrapeElement.Value.(*scrapeValue)
				scrapeValue.value = element.Text
				return
			}
			span := element.DOM.Find("span")
			if strings.Contains(strings.ToLower(span.Text()), "market cap") {
				nextIsValue = true
				scrapeStack.PushBack(&scrapeValue{
					key: keyMarketCap,
				})
			} else if strings.EqualFold(span.Text(), "enterprise value") {
				nextIsValue = true
				scrapeStack.PushBack(&scrapeValue{
					key: keyEnterpriseValue,
				})
			}
		})
	})

	nextIsValue := false
	c.OnHTML("td", func(e *colly.HTMLElement) {
		if nextIsValue {
			nextIsValue = false
			scrapeElement := scrapeStack.Back()
			scrapeValue := scrapeElement.Value.(*scrapeValue)
			scrapeValue.value = e.Text
			return
		}
		classes := e.Attr("class")
		if strings.Contains(classes, "fi-row") {
			span := e.DOM.Find("span")
			if strings.Contains(strings.ToLower(span.Text()), "return on equity") {
				nextIsValue = true
				scrapeStack.PushBack(&scrapeValue{
					key: keyReturnOnEquity,
				})
			} else if strings.Contains(strings.ToLower(span.Text()), "held by insiders") {
				nextIsValue = true
				scrapeStack.PushBack(&scrapeValue{
					key: keyInsiderOwnership,
				})
			}
		}
	})

	url = fmt.Sprintf(statisticsURLFormat, symbol, symbol)
	err = c.Visit(url)
	if err != nil {
		log.Errorf("failed to visit url %s: %v", url, err)
	}

	keysToValue := make(map[string]string)
	for e := scrapeStack.Front(); e != nil; e = e.Next() {
		sv := e.Value.(*scrapeValue)
		keysToValue[sv.key] = sv.value
	}

	return finance.ScrapeResult{
		PERatio:          keysToValue[keyPERatio],
		EnterpriseValue:  keysToValue[keyEnterpriseValue],
		MarketCap:        keysToValue[keyMarketCap],
		ReturnOnEquity:   keysToValue[keyReturnOnEquity],
		InsiderOwnership: keysToValue[keyInsiderOwnership],
	}, nil
}
