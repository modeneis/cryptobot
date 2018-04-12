package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/modeneis/cryptobot/src/model"
	"gopkg.in/mgo.v2/bson"
)

var (
	bittrexAPIV2URL = "https://bittrex.com/Api/v2.0/pub"
)

type BittrexAPIv2 struct{}
type Response struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

// GetTicker sends a public get request and returns current ticker information
// on the supplied currency.
// Example:
// currencyPair=USDT-BTC, interval=hour
// currencyPair=USDT-BTC, interval=oneMin
// currencyPair=USDT-BTC, interval=fiveMin
// currencyPair=USDT-BTC, interval=thirtyMin
// currencyPair=USDT-BTC, interval=hour
// currencyPair=USDT-BTC, interval=day
func (b *BittrexAPIv2) GetTicker(currencyPair, interval string) (candles model.Candles, err error) {
	//https://bittrex.com/Api/v2.0/pub/market/GetTicks?marketName=USDT-BTC&tickInterval=hour
	path := fmt.Sprintf("%s/market/GetTicks?marketName=%s&tickInterval=%s", bittrexAPIV2URL, currencyPair, interval)

	var tempCandles = model.Candles{}
	err = b.HTTPRequest(path, false, url.Values{}, &tempCandles.CandleLS)

	for _, candle := range tempCandles.CandleLS {
		if candle.ID == "" {
			candle.ID = bson.NewObjectId().Hex()
			candles.CandleLS = append(candles.CandleLS, candle)
		}
	}

	return
}

func (b *BittrexAPIv2) HTTPRequest(path string, auth bool, values url.Values, v interface{}) error {
	response := Response{}
	if err := SendHTTPGetRequest(path, true, &response); err != nil {
		return err
	}
	if response.Success {
		return json.Unmarshal(response.Result, &v)
	}
	return errors.New(response.Message)
}

func SendHTTPGetRequest(url string, jsonDecode bool, result interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		log.Printf("HTTP status code: %d\n", res.StatusCode)
		log.Printf("URL: %s\n", url)
		return errors.New("status code was not 200")
	}

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if jsonDecode {
		err := json.Unmarshal(contents, &result)
		if err != nil {
			log.Println(string(contents[:]))
			return err
		}
	} else {
		result = &contents
	}

	return nil
}
