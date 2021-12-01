package main

import (
	"flag"
	"github.com/hill-daniel/finance-scraper"
	"github.com/hill-daniel/finance-scraper/colly"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	lvl, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = log.InfoLevel
	}
	customFormatter := &log.TextFormatter{}
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)
}

func main() {
	symbol := flag.String("symbol", "", "The symbol to check")
	analyzerType := flag.String("analyzer", "", "The analyzer for the result")
	flag.Parse()

	collector, err := colly.NewYahooCollector()
	if err != nil {
		log.Fatalf("failed to create collector: %v", err)
	}

	result, err := collector.Scrape(*symbol)
	if err != nil {
		log.Fatalf("failed to scrape: %v", err)
	}

	mapper := finance.YahooMapper{}
	quote, err := mapper.Map(result)
	if err != nil {
		log.Fatalf("failed to map: %v", err)
	}

	analyzer := finance.CreateAnalyzer(*analyzerType)
	analysisResult, err := analyzer.Analyze(quote)
	if err != nil {
		log.Fatalf("failed to analyze: %v", err)
	}
	log.Infof("Analysis result:\n%s", analysisResult)
}
