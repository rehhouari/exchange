package exchange

import (
	"fmt"
	"testing"
)

func TestExchange(t *testing.T) {
	// Create a new Exchange instance and set USD as the base currency for the exchange rates and conversion
	ex := New("USD")

	fmt.Println("ForexCodes:")
	// Get the available forex/fiat codes ([]string)
	forexCodes, _ := ex.ForexCodes()
	fmt.Println(forexCodes)
	fmt.Println("======================")
/* 	// Get the available crypto codes ([]string)
	// Warning: +6000
	fmt.Println("CryptoCodes:")
	cryptoCodes, _ := ex.CryptoCodes()
	fmt.Println(cryptoCodes)
	fmt.Println("======================")

	// Get the forex codes data, includes code and description.
	forexData, _ := ex.ForexData()

	// Get the crypto codes data, includes code and description.
	// Warning: +6000
	cryptoData, _ := ex.CryptoData()
	fmt.Println("cryptoData[\"BTC\"][\"code\"]:")
	fmt.Println(cryptoData["BTC"]["code"])
	fmt.Println("======================")

	fmt.Println("forexData:")
	// loop through the forex cpdes
	for _, code := range forexCodes {
		// print the forex codes data in the format: USD: US Dollar
		fmt.Println(code+":", forexData[code]["description"])
	}
	fmt.Println("======================")

	fmt.Println("Sources:")
	sources, _ := ex.Sources()
	fmt.Println(sources)
	fmt.Println("======================")

	fmt.Println("Sources:")
	sources, _ := ex.Sources()
	fmt.Println(sources)
	fmt.Println("======================")

	fmt.Println("SourcesData:")
	// loop through the forex cpdes
	sourcesData, _ := ex.SourcesData()
	for _, code := range sources {
		fmt.Println(code+":", sourcesData[code]["description"])
	}
	fmt.Println("======================") */
	/* 	// Change the base currency to euro
	   	ex.SetBase("EUR")
	   	// Get the latest exchange rates with all currencies (Base is EUR)
	   	fmt.Println(ex.LatestRatesAll())

	   	// Get the latest rates again, this time it will be loaded from in-memory cache
	   	// Cache last till midnight GMT because it's the time exchangerate.host update the rates
	   	fmt.Println(ex.LatestRatesAll())
	   	// disable caching
	   	ex.SetCache(false)
	   	// Get the latest rates with multiple currencies, not all (USD and JPY only)
	   	fmt.Println(ex.LatestRatesMultiple([]string{"USD", "JPY"}))

	   	// Get the exchange rates at 2012-12-12 but only with USD
	   	fmt.Println(ex.HistoricalRatesSingle("2012-12-12", "USD"))

	   	// Get historical rates between 2012 12 10 and 2012 12 12 for JPY and GBP
	   	fmt.Println(ex.TimeseriesMultiple("2012-12-10", "2012-12-12", []string{"USD", "JPY"}))

	   	// Get the fluctuation between 2012 12 10 and 2012 12 12 with USD
	   	fluctuation, _ := ex.FluctuationSingle("2012-12-10", "2012-12-12", "USD")
	   	// Print the change
	   	fmt.Println(fluctuation["change"]) */
}
