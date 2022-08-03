// thanks to riccardo and briatore :) https://github.com/RiccardoM/briatore/blob/main/reporter/coingecko.go

package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	CoingeckoEndpoint           = "https://api.coingecko.com/api/v3/coins/{id}"
	CoingeckoHistoricalEndpoint = "https://api.coingecko.com/api/v3/coins/{id}/history?date={date}"
)

type HistoryResponse struct {
	MarketData *MarketData `json:"market_data"`
}

func (h HistoryResponse) GetCoinPrice(currency string) (float64, error) {
	if h.MarketData == nil {
		return 0, nil
	}

	price, ok := h.MarketData.CurrentPrice[currency]
	if !ok {
		return 0, fmt.Errorf("invalid currency: %s", currency)
	}

	return price, nil
}

type MarketData struct {
	CurrentPrice map[string]float64 `json:"current_price"`
}

type PriceData struct {
	CoinGeckoID string    `json:"coinGeckoID"`
	Price       float64   `json:"price"`
	Timestamp   time.Time `json:"timestamp"`
	Currency    string    `json:"currency"`
}

func NewPriceData(coinGeckoID string, price float64, currency string, timestamp time.Time) PriceData {
	return PriceData{
		CoinGeckoID: coinGeckoID,
		Price:       price,
		Currency:    currency,
		Timestamp:   timestamp,
	}
}

// GetCoinPrice gets the current price of the coin having the given CoinGecko ID,
// measured in the given currency.
func GetCoinPrice(id string, currency string) (float64, error) {
	price, err := getPriceFromAPI(id, CoingeckoEndpoint, time.Time{}, currency)
	if err != nil {
		return 0, err
	}

	priceData := NewPriceData(id, price, currency, time.Time{})

	return priceData.Price, nil
}

// GetHistoricalCoinPrice gets the historical price of the coin having the given CoinGecko ID,
// measured in the given currency.
func GetHistoricalCoinPrice(id string, timestamp time.Time, currency string) (float64, error) {
	price, err := getPriceFromAPI(id, CoingeckoHistoricalEndpoint, timestamp, currency)
	if err != nil {
		return 0, err
	}

	priceData := NewPriceData(id, price, currency, timestamp)

	return priceData.Price, nil
}

// getPriceFromAPI returns the price for the coin having the given id for the given timestamp and currency
func getPriceFromAPI(id, coingeckoEndpoint string, timestamp time.Time, currency string) (float64, error) {
	endpoint := strings.ReplaceAll(coingeckoEndpoint, "{id}", id)
	if coingeckoEndpoint == CoingeckoHistoricalEndpoint {
		endpoint = strings.ReplaceAll(endpoint, "{date}", timestamp.Format("02-01-2006"))
	}

	res, err := http.Get(endpoint)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad token price history response: status %d", res.StatusCode)
	}

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var response HistoryResponse
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return 0, err
	}

	return response.GetCoinPrice(currency)
}
