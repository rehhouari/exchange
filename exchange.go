package exchange

import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidCode is returned when the currency code is invalid
var ErrInvalidCode = errors.New("Invalid currency code")

// ErrInvalidDate is returned when the date is too old
var ErrInvalidDate = errors.New("Oldest possible date is 1999-01-04")

// ErrInvalidDateFormat is returned when the date isn't formatted as YYYY-MM-DD
var ErrInvalidDateFormat = errors.New("Date format must be YYYY-MM-DD")

// ErrTimeframeExceeded is returned when the time between start_date and end_date is bigger than 365 days
var ErrTimeframeExceeded = errors.New("Maximum allowed timeframe is 365 days")

// ErrInvalidTimeFrame is returned when the to date is older than to date. For example flipped the arguments
var ErrInvalidTimeFrame = errors.New("From date must be older than To date")

// ErrInvalidAPIResponse is returned when the API return success: false
var ErrInvalidAPIResponse = errors.New("Unknown API error")

const (
	baseURL             string = "https://api.exchangerate.host"
	symbolsURL          string = baseURL + "/symbols"
	cryptocurrenciesURL string = baseURL + "/cryptocurrencies"
	latestURL           string = baseURL + "/latest"
	convertURL          string = baseURL + "/convert"
	historicalURL       string = baseURL + "/"
	timeseriesURL       string = baseURL + "/timeseries"
	fluctuationURL      string = baseURL + "/fluctuation"
	sourcesURL          string = baseURL + "/sources"
)

// Exchange is returned by New() and allows access to the methods
type Exchange struct {
	Base          string
	isInitialized bool // is set to true if used via New
}

type query struct {
	From      string
	To        string
	Base      string
	Amount    float64
	Symbols   []string
	Date      string
	TimeFrame [2]string
}

var client http.Client = http.Client{}

// New creates a new instance of Exchange
func New(base string) *Exchange {
	x := &Exchange{
		Base:          base,
		isInitialized: true,
	}
	return x
}

// SetBase sets a new base currency for the exchange rates
func (exchange *Exchange) SetBase(base string) error {
	if err := ValidateCode(base); err != nil {
		return err
	}
	exchange.Base = base
	return nil
}

// ValidateCode validates a single symbol code
func ValidateCode(code string) error {
	if len(code) != 3 {
		return ErrInvalidCode
	}
	return nil
}

// ValidateSymbols validates all symbols' codes in an array
func ValidateSymbols(symbols []string) error {
	for code := range symbols {
		if err := ValidateCode(symbols[code]); err != nil {
			return err
		}
	}

	return nil
}

// ValidateDate validates date string according to YYYY-MM-DD format and if it's
func ValidateDate(date string) error {
	matched, err := regexp.Match("[0-9]{4,4}-((0[1-9])|(1[0-2]))-([0-3]{1}[0-9]{1})", []byte(date))
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidDateFormat
	}
	oldestDate, _ := time.Parse("2006-01-02", "1999-01-03")
	selectedDate, _ := time.Parse("2006-01-02", date)
	if selectedDate.Before(oldestDate) {
		return ErrInvalidDate
	}
	return nil
}

// ValidateTimeFrame checks if the from and to date are not more than 365 days apart and they're not mixed
func ValidateTimeFrame(TimeFrame [2]string) error {
	from, _ := time.Parse("2006-01-02", TimeFrame[0])
	to, _ := time.Parse("2006-01-02", TimeFrame[1])
	if to.Before(from) {
		return ErrInvalidTimeFrame
	}

	if to.Sub(from).Hours() > 8759.992992006 {
		return ErrTimeframeExceeded
	}

	return nil
}

func (exchange *Exchange) get(url string, q query) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	processQuery(req, q)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	success := result["success"]

	if !success.(bool) {
		return nil, ErrInvalidAPIResponse
	}

	return result, nil
}

func addToQuery(req *http.Request, key string, value string) {
	q := req.URL.Query()          // Get a copy of the query values.
	q.Add(key, value)             // Add a new value to the set.
	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
}

func processQuery(req *http.Request, q query) error {
	if q.Base != "" {
		if err := ValidateCode(q.Base); err != nil {
			return err
		}
		addToQuery(req, "base", q.Base)
	}

	if q.From != "" {
		if err := ValidateCode(q.From); err != nil {
			return err
		}
		addToQuery(req, "from", q.From)
	}

	if q.To != "" {
		if err := ValidateCode(q.To); err != nil {
			return err
		}
		addToQuery(req, "to", q.To)
	}

	if q.Amount > 1 {
		addToQuery(req, "amount", strconv.Itoa(q.Amount))
	}

	if len(q.Symbols) != 0 {
		addToQuery(req, "symbols", strings.Join(q.Symbols, ","))
	}

	if q.Date != "" {
		if err := ValidateDate(q.Date); err != nil {
			return err
		}
		addToQuery(req, "date", q.Date)
	}

	if q.TimeFrame != [2]string{} {
		for i := 0; i < 1; i++ {
			if err := ValidateDate(q.TimeFrame[i]); err != nil {
				return err
			}
		}
		if err := ValidateTimeFrame(q.TimeFrame); err != nil {
			return err
		}
		addToQuery(req, "start_date", string(q.TimeFrame[0]))
		addToQuery(req, "end_date", string(q.TimeFrame[1]))
	}

	return nil
}

func (exchange *Exchange) apiSymbols() (map[string]map[string]string, error) {
	resp, err := exchange.get(symbolsURL, query{})
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]string)
	for code, data := range resp["symbols"].(map[string]interface{}) {
		values := make(map[string]string)
		for name, value := range data.(map[string]interface{}) {
			values[name] = value.(string)
		}
		result[code] = values
	}
	return result, nil
}

func (exchange *Exchange) apiSource() (map[string]map[string]string, error) {
	resp, err := exchange.get(sourcesURL, query{})
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]string)
	for code, data := range resp["sources"].(map[string]interface{}) {
		values := make(map[string]string)
		for name, value := range data.(map[string]interface{}) {
			values[name] = value.(string)
		}
		result[code] = values
	}
	return result, nil
}

func (exchange *Exchange) apiCryptocurrencies() (map[string]map[string]string, error) {
	resp, err := exchange.get(cryptocurrenciesURL, query{})
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]string)
	for code, data := range resp["cryptocurrencies"].(map[string]interface{}) {
		values := make(map[string]string)
		for name, value := range data.(map[string]interface{}) {
			values[name] = value.(string)
		}
		result[code] = values
	}
	return result, nil
}

func (exchange *Exchange) apiLatest(q query) (map[string]*decimal.Decimal, error) {
	resp, err := exchange.get(latestURL, q)
	if err != nil {
		return nil, err
	}
	result := resp["rates"].(map[string]interface{})
	rates := make(map[string]*decimal.Decimal, len(result))
	for key := range result {
		rates[key] = decimal.NewFromFloat(result[key].(float64))
	}
	return rates, nil
}

func (exchange *Exchange) apiConvert(q query) (*decimal.Decimal, error) {
	resp, err := exchange.get(convertURL, q)
	if err != nil {
		return nil, err
	}
	result := resp["result"].(float64)
	return decimal.NewFromFloat(result), nil
}

func (exchange *Exchange) apiHistorical(q query) (map[string]*big.Float, error) {
	if err := ValidateDate(q.Date); err != nil {
		return nil, err
	}
	url := historicalURL + q.Date
	q.Date = ""
	resp, err := exchange.get(url, q)
	if err != nil {
		return nil, err
	}
	result := resp["rates"].(map[string]interface{})
	rates := make(map[string]*big.Float, len(result))
	for key := range result {
		rates[key] = big.NewFloat(result[key].(float64))
	}
	return rates, nil
}

func (exchange *Exchange) apiTimeseriesAndFuctuation(url string, q query) (map[string]map[string]*big.Float, error) {
	resp, err := exchange.get(url, q)
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]*big.Float)
	for date, rates := range resp["rates"].(map[string]interface{}) {
		ratemap := make(map[string]*big.Float)
		for symbol, rate := range rates.(map[string]interface{}) {
			frate := big.NewFloat(rate.(float64))
			ratemap[symbol] = frate
			result[date] = ratemap
		}
	}
	return result, nil
}

// ForexCodes returns and array of supported forex/fiat currency codes
func (exchange *Exchange) ForexCodes() ([]string, error) {
	var codes []string

	result, err := exchange.apiSymbols()
	if err != nil {
		return nil, err
	}

	for k := range result {
		codes = append(codes, k)
	}

	sort.Strings(codes)
	return codes, nil
}

// ForexData returns a map of supported forex/fiat currencies data (code & description)
func (exchange *Exchange) ForexData() (map[string]map[string]string, error) {
	return exchange.apiSymbols()
}

// CryptoCodes returns and array of supported cryptocurrency codes
func (exchange *Exchange) CryptoCodes() ([]string, error) {
	var codes []string

	result, err := exchange.apiCryptocurrencies()
	if err != nil {
		return nil, err
	}

	for k := range result {
		codes = append(codes, k)
	}

	sort.Strings(codes)
	return codes, nil
}

// CryptoData returns a map of supported cryptocurrencies data (name and symbol)
func (exchange *Exchange) CryptoData() (map[string]map[string]string, error) {
	return exchange.apiCryptocurrencies()
}

// Sources returns and array of supported sources
func (exchange *Exchange) Sources() ([]string, error) {
	var codes []string

	result, err := exchange.apiSources()
	if err != nil {
		return nil, err
	}

	for k := range result {
		codes = append(codes, k)
	}

	sort.Strings(codes)
	return codes, nil
}

// SourcesData returns a map of supported forex/fiat currencies data (code & description)
func (exchange *Exchange) SourcesData() (map[string]map[string]string, error) {
	return exchange.apiSources()
}
