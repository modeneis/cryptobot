package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/modeneis/cryptobot/src/model"
	"gopkg.in/mgo.v2/bson"
)

func GenericHandler(res http.ResponseWriter, req *http.Request) {

	AllowOrigin(res, req)

	var result []interface{}
	var err error

	id := req.URL.Query().Get(":id")
	table := req.URL.Query().Get(":table")

	//set mime type to JSON, Its JSON REST API
	res.Header().Set("Content-type", "application/json")

	var thing interface{}

	// TODO: ADD all supported models here
	if table == "thing" {
		thing = new(model.Thing)
	} else if table == "market" {
		thing = new(model.Market)
	} else {
		log.Println("ERROR: GenericHandler, table/model not found", table, err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	// Db and collection name configs
	c := Session.DB(MONGODB_DATABASE).C(table)
	//if table == "market" {
	//	c = Session.DB(MONGODB_CRYPTO).C(table)
	//}
	// Handle the methods and behave accordingly
	switch req.Method {
	case "GET":

		/***************** BEGIN GET ****************************/
		if id != "" {
			log.Println("INFO: GET ONE ", table, id)
			err = c.Find(bson.M{"_id": id}).All(&result)
			if err != nil {
				log.Println("ERROR: GenericHandler, ", err)
				res.WriteHeader(http.StatusNotFound)
				return
			}
		} else {
			log.Println("INFO: GET ALL ", table)
			err = c.Find(nil).All(&result)
		}

		js, err := json.Marshal(result)
		if err != nil || len(js) == 0 {
			log.Println("ERROR: GenericHandler, ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Write(js)
		return

		/***************** END GET ****************************/

	case "POST":

		/***************** BEGIN POST ****************************/
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&thing)

		//// SET ID
		//// TODO: ADD all supported models here
		if table == "thing" {
			thing.(*model.Thing).ID = bson.NewObjectId().Hex()
		} else if table == "market" {
			thing.(*model.Market).ID = bson.NewObjectId().Hex()
		}

		// Insert thing into db
		err = c.Insert(thing)
		if err != nil {
			log.Println("ERROR: GenericHandler, ", err)
			res.WriteHeader(http.StatusNotFound)
			return
		}

		js, err := json.Marshal(thing)
		if err != nil || len(js) == 0 {
			log.Println("ERROR: GenericHandler, ", err)
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Write(js)
		return
		/***************** END POST ****************************/

	case "PUT":

		/***************** BEGIN PUT ****************************/
		decoder := json.NewDecoder(req.Body)
		// r.PostForm is a map of our POST form values
		err := decoder.Decode(&thing)

		if err != nil {
			log.Println("ERROR: GenericHandler, Decode 404 not found, ", err)
			res.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println("INFO: GenericHandler, table ", table, ", ID ", id, thing)

		// Update
		info, err := c.Upsert(bson.M{"_id": id}, thing)
		if err != nil {
			log.Println("ERROR: GenericHandler, ", err)
			res.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println("INFO: ", info)

		res.WriteHeader(http.StatusOK)
		return
		/***************** END PUT ****************************/

	case "DELETE":
		/***************** BEGIN DELETE ****************************/
		if id == "" {
			log.Println("ERROR: GenericHandler, ID is required", table, err)
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}
		// When a panda leaves :(, Delete from database
		err = c.Remove(bson.M{"_id": id})
		if err != nil {
			log.Println("ERROR: GenericHandler, ", err)
			res.WriteHeader(http.StatusNotFound)
			return
		}

		res.WriteHeader(http.StatusOK)
		return
		/***************** END DELETE ****************************/

	default:
		log.Println("INFO: GenericHandler, OPTIONS")
		res.WriteHeader(http.StatusOK)
		return
	}

}

func AllowOrigin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, CRYPTOBOT-KEY")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, HEAD, POST, PUT, OPTIONS")

	//TODO: add origin validation
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

}
