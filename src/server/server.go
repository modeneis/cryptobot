package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/pat"
	"github.com/modeneis/cryptobot/src/apis"
	"github.com/modeneis/cryptobot/src/markets"
	"github.com/thrasher-/gocryptotrader/config"
	"gopkg.in/mgo.v2"
)

var (
	router *pat.Router
	//Session defines the mongodb session

	//Cache
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
var Database *mgo.Database
var Collection *mgo.Collection

var RESTART_TIMEOUT_ENABLED = os.Getenv("RESTART_TIMEOUT_ENABLED")
var RESTART_TIMEOUT_FACTOR_STR = os.Getenv("RESTART_TIMEOUT_FACTOR")
var RESTART_TIMEOUT_FACTOR = 1000

var WATCH_MARKETS = os.Getenv("WATCH_MARKETS")
var WATCH_TIME_STR = os.Getenv("WATCH_TIME")
var WATCH_TIME = 10

var Table = os.Getenv("table")

var DEBUG = os.Getenv("DEBUG")

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

	if RESTART_TIMEOUT_FACTOR_STR == "" {
		RESTART_TIMEOUT_FACTOR_STR = "1000"
	}

	DB = MONGODB_HOST + ":" + MONGODB_PORT

	log.Println("INFO: Loading standalone MONGODB ", DB)
	Session, err = mgo.Dial("mongodb://" + DB)

	if err != nil {
		log.Fatalf("Could not connect to MONGODB, Error:  %v", err.Error())
	}

	if MONGODB_USERNAME != "" {
		err = Session.Login(&mgo.Credential{Username: MONGODB_USERNAME, Password: MONGODB_PASSWORD})
		if err != nil {
			log.Fatalf("Could not connect to MONGODB, Error:  %v", err.Error())
		}
	}

	Session.SetMode(mgo.Eventual, false)

	SetupConfig = config.GetConfig()
	//SetupConfig.EncryptConfig = -1
	SetupConfig.LoadConfig("./config.dat")

	if WATCH_TIME_STR == "" {
		WATCH_TIME_STR = "10"
	}

	WATCH_TIME, _ = strconv.Atoi(WATCH_TIME_STR)

}

//NewServer return pointer to new created server object
func NewServer(Port string) *http.Server {
	router = InitRouting()
	return &http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}
}

//StartServer start and listen @server
func StartServer(Port string) {
	log.Println("**************** Starting server *****************")
	s := NewServer(Port)
	log.Println("Server starting --> " + Port)

	//if RESTART_TIMEOUT_ENABLED == "true" {
	//	var RESTART_TIMEOUT = getRandomValueFromInterval(0.5, rand.New(rand.NewSource(time.Now().UnixNano())).Float64(), time.Duration(RESTART_TIMEOUT_FACTOR)*time.Second)
	//	//enable auto restart
	//	log.Println("******************************* RESTART_TIMEOUT_ENABLED, timeout is -> ", RESTART_TIMEOUT)
	//	time.AfterFunc(RESTART_TIMEOUT, func() {
	//		log.Println("******************************* WARN: Node will restart after timeout, (time.Minute) ", RESTART_TIMEOUT)
	//		os.Exit(0)
	//	})
	//}

	if WATCH_MARKETS == "true" {

		marketCollection := Session.DB(MONGODB_DATABASE).C("market")
		go watchMarkets(marketCollection, WATCH_TIME)

		oneMinCandleCollection := Session.DB(MONGODB_DATABASE).C("onemin_candle")
		go watchCandles(oneMinCandleCollection, 1, "USDT-BTC", "oneMin")

		fiveminCandleCollection := Session.DB(MONGODB_DATABASE).C("fivemin_candle")
		go watchCandles(fiveminCandleCollection, 5, "USDT-BTC", "fiveMin")

		thirtyminCandleCollection := Session.DB(MONGODB_DATABASE).C("thirtymin_candle")
		go watchCandles(thirtyminCandleCollection, 30, "USDT-BTC", "thirtyMin")

		hourCandleCollection := Session.DB(MONGODB_DATABASE).C("hour_candle")
		go watchCandles(hourCandleCollection, 60, "USDT-BTC", "hour")

		dayCandleCollection := Session.DB(MONGODB_DATABASE).C("day_candle")
		go watchCandles(dayCandleCollection, 60*24, "USDT-BTC", "day")

	}

	//enable graceful shutdown
	err := gracehttp.Serve(
		s,
	)

	if err != nil {
		log.Fatalln("Error: %v", err)
		os.Exit(0)
	}

}

func getRandomValueFromInterval(randomizationFactor, random float64, currentInterval time.Duration) time.Duration {
	var delta = randomizationFactor * float64(currentInterval)
	var minInterval = float64(currentInterval) - delta
	var maxInterval = float64(currentInterval) + delta

	// Get a random value from the range [minInterval, maxInterval].
	// The formula used below has a +1 because if the minInterval is 1 and the maxInterval is 3 then
	// we want a 33% chance for selecting either 1, 2 or 3.
	return time.Duration(minInterval + (random * (maxInterval - minInterval + 1)))
}

func InitRouting() *pat.Router {

	r := pat.New()

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	assets := http.StripPrefix("/assets/", http.FileServer(http.Dir("./static/site/startup-kit/assets/")))
	r.PathPrefix("/assets/").Handler(assets)

	ss := http.StripPrefix("/test", http.FileServer(http.Dir("./templates/")))
	r.PathPrefix("/test").Handler(ss)

	r.Handle("/{page}", TemplateHandler("static/site/startup-kit/"))
	//r.Handle("/site/{page}", ShowCommentsHTML("static/site/startup-kit/"))

	r.Handle("/", TemplateHandler("static/site/startup-kit/"))

	r.HandleFunc("/v1/create/{table}/", GenericHandler)
	r.HandleFunc("/v1/read/{table}/{id}/", GenericHandler)
	r.HandleFunc("/v1/update/{table}/{id}/", GenericHandler)
	r.HandleFunc("/v1/delete/{table}/{id}/", GenericHandler)

	return r
}

func watchMarkets(c *mgo.Collection, t int) {

	time.AfterFunc(time.Duration(t)*time.Second, func() {

		marketLS, err := markets.GetBitTrexTicker(SetupConfig)
		if err != nil {
			log.Println("ERROR: markets.GetBitTrexTicker, ", err)
		} else {
			err = markets.SaveMarkets(marketLS, c)
			if err != nil {
				log.Println("ERROR: markets.SaveMarkets, ", err)
			}
		}
		watchMarkets(c, t)
	})
}

var bittrexv2 = apis.BittrexAPIv2{}

func watchCandles(c *mgo.Collection, t int, currencyPair, interval string) {
	Debug("START watchCandles, ", t, currencyPair, interval)

	time.AfterFunc(time.Duration(t)*time.Minute, func() {

		candles, err := bittrexv2.GetTicker(currencyPair, interval)
		if err != nil {
			log.Println("ERROR: bittrexv2.GetTicker, ", err)
		} else {
			err = markets.SaveCandles(candles.CandleLS, c)
			if err != nil {
				log.Println("ERROR: markets.SaveCandles", err)
			}
		}
		Debug("END, watchCandles, ", t, currencyPair, interval)
		watchCandles(c, t, currencyPair, interval)
	})
}

func Debug(v ...interface{}) {
	if DEBUG == "true" {
		log.Println("DEBUG: ", v)
	}
}
