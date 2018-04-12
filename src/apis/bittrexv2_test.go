package apis_test

import (
	"github.com/modeneis/cryptobot/src/apis"
	"log"
	"os"
	"testing"

	"github.com/modeneis/cryptobot/src/markets"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

var DB = ""

var MONGODB_HOST = os.Getenv("MONGODB_HOST")
var MONGODB_PORT = os.Getenv("MONGODB_PORT")
var MONGODB_USERNAME = os.Getenv("MONGODB_USERNAME")
var MONGODB_PASSWORD = os.Getenv("MONGODB_PASSWORD")
var MONGODB_DATABASE = os.Getenv("MONGODB_DATABASE")
var Session *mgo.Session

func init() {
	var err error

	if MONGODB_HOST == "" {
		MONGODB_HOST = "localhost"
	}

	if MONGODB_PORT == "" {
		MONGODB_PORT = "27017"
	}

	if MONGODB_DATABASE == "" {
		MONGODB_DATABASE = "test"
	}

	DB = MONGODB_HOST + ":" + MONGODB_PORT

	log.Println("INFO: Loading standalone mgo ", DB)

	if MONGODB_USERNAME != "" {
		Session, err = mgo.Dial("mongodb://" + MONGODB_USERNAME + ":" + MONGODB_PASSWORD + "@" + DB)
	} else {
		Session, err = mgo.Dial("mongodb://" + DB)
	}

	Session.SetMode(mgo.Eventual, false)

	if err != nil {
		log.Fatalf("Could not connect to rethinkdb, Error:  %v", err.Error())
	}
}

func TestGetBitTrex_hour_candle(t *testing.T) {

	var bittrexv2 = apis.BittrexAPIv2{}
	Convey("Given I call BitTrex Ticker APIv2, Then I should get array of Ticks ", t, func() {

		candles, err := bittrexv2.GetTicker("USDT-BTC", "hour")
		So(err, ShouldBeNil)
		So(len(candles.CandleLS), ShouldBeGreaterThan, 0)

		Convey("Then I should be able to insert it into mongodb", func() {
			c := Session.DB(MONGODB_DATABASE).C("hour_candle")
			err := markets.SaveCandles(candles.CandleLS, c)
			So(err, ShouldBeNil)

			Convey("Then I make sure they are really saved", func() {
				var ids []string

				for _, c := range candles.CandleLS {
					ids = append(ids, c.ID)
				}

				CandleLS, err := markets.FindCandleByIDs(ids, c)
				So(err, ShouldBeNil)
				So(CandleLS, ShouldNotBeNil)
				So(len(CandleLS), ShouldBeGreaterThan, 0)

			})
		})
	})
}
