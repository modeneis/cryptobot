package markets

import (
	"log"

	"os"
	"time"

	"github.com/modeneis/cryptobot/src/model"
	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/exchanges/anx"
	"github.com/thrasher-/gocryptotrader/exchanges/bittrex"
	"gopkg.in/mgo.v2/bson"
)

var DEBUG = os.Getenv("DEBUG")

func GoGetMarkets(name string, SetupConfig *config.Config) {

	if name == "ANX" {
		GetANX(SetupConfig)
	} else if name == "Bittrex" {
		GetBitTrexTicker(SetupConfig)
	}
}

func GetBitTrexTicker(SetupConfig *config.Config) ([]model.Market, error) {
	var err error
	name := "Bittrex"
	setup := bittrex.Bittrex{}
	setup.Name = name

	bittrexConfig, err := SetupConfig.GetExchangeConfig(name)
	if err != nil {
		log.Println("ERROR: GetBitTrex, Setup() init error, ", name)
	}
	setup.Setup(bittrexConfig)

	if setup.Enabled != true {
		log.Println("ERROR: GetBitTrex, Setup() incorrect values set, ", name)
	}

	getTicker := bittrex.Bittrex{}
	var markets []model.Market
	for _, v := range setup.EnabledPairs { //AvailablePairs

		ticker, err := getTicker.GetTicker(v)
		if err != nil {
			log.Printf("ERROR: GetBitTrex,  GetTicker() invalid pair %s, : %s ", v, err)
			err = nil
		} else {
			market := model.Market{}

			market.Name = "Bittrex"
			market.Ask = ticker.Result.Ask
			market.Bid = ticker.Result.Bid
			market.Last = ticker.Result.Last
			market.Pair = v
			market.UpdatedAt = time.Now()

			market.ID = bson.NewObjectId().Hex()

			Debug("***************** ")
			Debug("Name, Pair, Ask, Bid, Last, Date -> ", market.Name, market.Pair, market.Ask, market.Bid, market.Last, market.UpdatedAt)
			Debug("***************** ")

			markets = append(markets, market)
		}
	}
	return markets, err

}

func GetANX(SetupConfig *config.Config) {
	name := "ANX"
	setup := anx.ANX{}
	setup.Name = name
	anxConfig, err := SetupConfig.GetExchangeConfig(name)
	if err != nil {
		log.Println("ERROR: GetANX, Setup() init error, ", name)
	}
	setup.Setup(anxConfig)

	if setup.Enabled != true {
		log.Println("ERROR: GetANX, Setup() incorrect values set, ", name)
	}

	getTicker := anx.ANX{}

	log.Println("INFO: GetANX, EXCHANGES =", name)

	for i, v := range setup.AvailablePairs {
		//fmt.Print("2**%d = %d\n", i, v)
		log.Println("INFO: GetANX, ", i)
		//fmt.Println(v)

		log.Printf("PAIR = %v \n", v)

		ticker, err := getTicker.GetTicker(v)

		if err != nil {
			log.Printf("ERROR: GetANX, %s, GetTicker() error: %s ", name, err)
		}
		if ticker.Result != "success" {
			log.Println("ERROR: GetANX, GetTicker() unsuccessful, ", name)
		}
		if ticker.Result == "success" {
			Debug("***************** ")
			Debug("High -> ", ticker.Data.High.Value)
			Debug("Low  -> ", ticker.Data.Low.Value)
			Debug("Sell -> ", ticker.Data.Sell.Value)
			Debug("Vol  -> ", ticker.Data.Vol.Value)
			Debug("***************** ")

			//fmt.Println(ticker.Data)
		}
	}
}

func Debug(v ...interface{}) {
	if DEBUG == "true" {
		log.Println("DEBUG: ", v)
	}
}
