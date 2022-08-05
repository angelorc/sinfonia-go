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
	CoingeckoEndpoint = "https://api.coingecko.com/api/v3/coins/{id}/market_chart/range?vs_currency={currency}&from={from}&to={to}"
)

type HistoryResponse struct {
	Prices [][2]float64 `json:"prices"`
}

func GetHistoricalCoinPrice(id, vsCurrency string, startTime, endTime time.Time) ([][2]float64, error) {
	price, err := getPriceFromAPI(id, vsCurrency, startTime, endTime)
	if err != nil {
		return [][2]float64{{0, 0}}, err
	}

	return price, nil
}

func getPriceFromAPI(id, vsCurrency string, startTime, endTime time.Time) ([][2]float64, error) {
	endpoint := strings.ReplaceAll(CoingeckoEndpoint, "{id}", id)
	endpoint = strings.ReplaceAll(endpoint, "{currency}", vsCurrency)
	endpoint = strings.ReplaceAll(endpoint, "{from}", fmt.Sprintf("%d", startTime.Unix()))
	endpoint = strings.ReplaceAll(endpoint, "{to}", fmt.Sprintf("%d", endTime.Unix()))

	res, err := http.Get(endpoint)
	if err != nil {
		return [][2]float64{{0, 0}}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return [][2]float64{{0, 0}}, fmt.Errorf("bad token price history response: status %d", res.StatusCode)
	}

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return [][2]float64{{0, 0}}, err
	}

	var response HistoryResponse
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return [][2]float64{{0, 0}}, err
	}

	return response.Prices, nil
}
