package markets

import (
	"log"

	"github.com/modeneis/cryptobot/src/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func SaveCandles(candlesLS []model.Candle, c *mgo.Collection) (err error) {
	for _, candle := range candlesLS {
		err = c.Insert(candle)
		if err != nil {
			if mgo.IsDup(err) {
				//log.Println("WARN: SaveCandles, IsDup, ", err)
				err = nil // ignore this err type
			} else {
				log.Println("ERROR: SaveCandles, ", err)
			}
		}
	}
	return err
}

func FindCandleByIDs(ids []string, c *mgo.Collection) (candleLS []model.Candle, err error) {
	for _, id := range ids {
		responseObj := model.Candle{}
		err = c.Find(bson.M{"_id": id}).One(&responseObj)
		if err != nil {
			log.Println("ERROR: FindCandleByIDs, ", err)
			return candleLS, err
		}
		candleLS = append(candleLS, responseObj)
	}
	return candleLS, err
}

func SaveMarkets(marketLS []model.Market, c *mgo.Collection) (err error) {
	for _, market := range marketLS {
		err = c.Insert(market)
		if err != nil {
			log.Println("ERROR: SaveMarkets, ", err)
			return err
		}
	}
	return err
}

func FindMarketsByIDs(ids []string, c *mgo.Collection) (markets []model.Market, err error) {
	for _, id := range ids {
		responseObj := model.Market{}
		err = c.Find(bson.M{"_id": id}).One(&responseObj)
		if err != nil {
			log.Println("ERROR: ConsultMarkets, ", err)
			return markets, err
		}
		markets = append(markets, responseObj)
	}
	return markets, err
}
