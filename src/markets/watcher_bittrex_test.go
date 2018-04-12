package markets_test

import (
	"log"
	"os"
	"testing"

	"github.com/modeneis/cryptobot/src/markets"
	"github.com/modeneis/cryptobot/src/model"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/thrasher-/gocryptotrader/config"
	"gopkg.in/mgo.v2"
)

//DB eg: localhost:28015 : or website.com
var DB = ""

var MONGODB_HOST = os.Getenv("MONGODB_HOST")
var MONGODB_PORT = os.Getenv("MONGODB_PORT")
var MONGODB_USERNAME = os.Getenv("MONGODB_USERNAME")
var MONGODB_PASSWORD = os.Getenv("MONGODB_PASSWORD")
var MONGODB_DATABASE = os.Getenv("MONGODB_DATABASE")

//var MONGODB_CRYPTO = "crypto"

var Session *mgo.Session

var SetupConfig *config.Config

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

	SetupConfig = config.GetConfig()
	//SetupConfig.EncryptConfig = -1
	SetupConfig.LoadConfig("../config.dat")

}

func TestGetBitTrex(t *testing.T) {

	var err error
	var marke []model.Market
	Convey("Given I get all markets from BitTrex", t, func() {

		marke, err = markets.GetBitTrexTicker(SetupConfig)
		So(err, ShouldBeNil)
		So(len(marke), ShouldBeGreaterThan, 0)
		/******** TABLES *************/

		Convey("Then I save all markets into MongoDB", func() {
			c := Session.DB(MONGODB_DATABASE).C("market")
			err := markets.SaveMarkets(marke, c)
			So(err, ShouldBeNil)

			Convey("Then I make sure they are really saved", func() {
				var ids []string

				for _, mark := range marke {
					ids = append(ids, mark.ID)
				}

				marketsConsult, err := markets.FindMarketsByIDs(ids, c)
				So(err, ShouldBeNil)
				So(marketsConsult, ShouldNotBeNil)
				So(len(marketsConsult), ShouldBeGreaterThan, 0)

			})

		})

	})
}
