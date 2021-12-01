package finance

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	thousand = 1000.0
	million  = thousand * thousand
	billion  = million * thousand
	trillion = billion * thousand
)

// YahooMapper maps ScrapeResult, if its source is also Yahoo finance, to a quote.
type YahooMapper struct {
}

// Map maps ScrapeResult to a quote.
func (ym *YahooMapper) Map(result ScrapeResult) (Quote, error) {
	peRatio, err := parseFloat(result.PERatio)
	if err != nil {
		return Quote{}, errors.Wrapf(err, "failed to parse PE ratio %s", result.PERatio)
	}

	enterpriseValue, err := parseFloat(result.EnterpriseValue)
	if err != nil {
		return Quote{}, errors.Wrapf(err, "failed to parse enterprise value %s", result.EnterpriseValue)
	}

	marketCap, err := parseFloat(result.MarketCap)
	if err != nil {
		return Quote{}, errors.Wrapf(err, "failed to parse market cap value %s", result.MarketCap)
	}

	returnOnEquity, err := parseFloat(result.ReturnOnEquity)
	if err != nil {
		return Quote{}, errors.Wrapf(err, "failed to parse return on equity value %s", result.ReturnOnEquity)
	}

	insiderOwnership, err := parseFloat(result.InsiderOwnership)
	if err != nil {
		return Quote{}, errors.Wrapf(err, "failed to parse insider ownership value %s", result.InsiderOwnership)
	}

	return Quote{
		PERatio:          peRatio,
		EnterpriseValue:  enterpriseValue,
		MarketCap:        marketCap,
		ReturnOnEquity:   returnOnEquity,
		InsiderOwnership: insiderOwnership,
	}, nil
}

func parseFloat(toParse string) (float64, error) {
	if len(toParse) == 0 || toParse == "N/A" {
		return 0.0, nil
	}
	value := strings.ReplaceAll(strings.ReplaceAll(toParse, ",", ""), "%", "")
	if strings.HasSuffix(value, "B") {
		parsed, err := strconv.ParseFloat(value[0:len(value)-1], 64)
		if err != nil {
			return 0, err
		}
		return parsed * billion, nil
	}
	if strings.HasSuffix(value, "M") {
		parsed, err := strconv.ParseFloat(value[0:len(value)-1], 64)
		if err != nil {
			return 0, err
		}
		return parsed * million, nil
	}
	if strings.HasSuffix(value, "T") {
		parsed, err := strconv.ParseFloat(value[0:len(value)-1], 64)
		if err != nil {
			return 0, err
		}
		return parsed * trillion, nil
	}
	return strconv.ParseFloat(value, 64)
}
