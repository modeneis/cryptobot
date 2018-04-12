package markets_test

import (
	"log"
	"os"
	"testing"

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

	Session.SetMode(mgo.Monotonic, true)

}

func TestIndexMarketCollection(t *testing.T) {

	Convey("Given I create indexes for market collection", t, func() {

		c := Session.DB("test").C("market")

		var indexes []mgo.Index

		//name
		indexes = append(indexes, mgo.Index{
			Key:        []string{"name"},
			Background: true,
			Sparse:     false,
			Unique:     false,
			DropDups:   false,
		})

		//ask
		indexes = append(indexes, mgo.Index{
			Key:        []string{"ask"},
			Background: true,
		})

		//bid
		indexes = append(indexes, mgo.Index{
			Key:        []string{"bid"},
			Background: true,
		})

		//last
		indexes = append(indexes, mgo.Index{
			Key:        []string{"last"},
			Background: true,
		})

		//pair
		indexes = append(indexes, mgo.Index{
			Key:        []string{"pair"},
			Background: true,
		})

		//updatedat
		indexes = append(indexes, mgo.Index{
			Key:        []string{"updatedat"},
			Background: true,
		})

		ensureAllIndexes(t, c, indexes)

		Convey("Then all indexes should be created for market collection", func() {
			verifyAllIndex(t, c)
		})
	})
}

func TestIndexCandleCollection(t *testing.T) {

	Convey("Given I create indexes for candle collections", t, func() {

		var indexes []mgo.Index

		//open
		indexes = append(indexes, mgo.Index{
			Key:        []string{"open"},
			Background: true,
		})

		//close
		indexes = append(indexes, mgo.Index{
			Key:        []string{"close"},
			Background: true,
		})

		//high
		indexes = append(indexes, mgo.Index{
			Key:        []string{"high"},
			Background: true,
		})

		//low
		indexes = append(indexes, mgo.Index{
			Key:        []string{"low"},
			Background: true,
		})

		//volume
		indexes = append(indexes, mgo.Index{
			Key:        []string{"volume"},
			Background: true,
		})

		//basevolume
		indexes = append(indexes, mgo.Index{
			Key:        []string{"basevolume"},
			Background: true,
		})

		//timestamp
		indexes = append(indexes, mgo.Index{
			Key:        []string{"timestamp"},
			Background: true,
			DropDups:   true,
			Unique:     true,
		})

		var collections []string
		collections = append(collections, "onemin_candle")

		collections = append(collections, "fivemin_candle")

		collections = append(collections, "thirtymin_candle")

		collections = append(collections, "hour_candle")

		collections = append(collections, "day_candle")

		for _, collectionName := range collections {
			c := Session.DB("test").C(collectionName)
			ensureAllIndexes(t, c, indexes)
		}

		Convey("Then all indexes should be created for each collection", func() {
			for _, collectionName := range collections {
				c := Session.DB("test").C(collectionName)
				verifyAllIndex(t, c)
			}
		})
	})
}

func verifyAllIndex(t *testing.T, c *mgo.Collection) {
	inds, err := c.Indexes()
	if err != nil {
		t.Fatalf("Could not create index %v to MONGODB, Error:  %v", inds, err.Error())
	}
	So(len(inds), ShouldBeGreaterThan, 0)
	for i, j := range inds {
		//log.Printf("INFO: Index %d: %+v", i, j)
		So(i, ShouldBeBetween, -1, 8)
		So(j, ShouldNotBeNil)
	}
}

func ensureAllIndexes(t *testing.T, c *mgo.Collection, indexes [](mgo.Index)) {
	for _, index := range indexes {
		ensureIndex(t, c, index)
	}
}

func ensureIndex(t *testing.T, c *mgo.Collection, index mgo.Index) {
	err := c.EnsureIndex(index)
	if err != nil {
		t.Fatalf("Could not create index %+v to MONGODB collection %+v, Error:  %v", index, c, err.Error())
	}
	return
}
