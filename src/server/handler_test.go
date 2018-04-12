package server_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/drewolson/testflight"
	"github.com/modeneis/cryptobot/src/model"
	"github.com/modeneis/cryptobot/src/server"
	. "github.com/smartystreets/goconvey/convey"
)

// go test -v server/handler_test.go

func TestCRUDCommentaries(t *testing.T) {
	targetObject := model.Thing{
		Name: "test",
	}

	mux := server.InitRouting()

	targetURL := "/v1/create/thing/"
	testflight.WithServer(mux, func(r *testflight.Requester) {
		t.Log("testflight.WithServer --> ")
		Convey("Given I create a new thing on api and get ID back", t, func() {
			raw, err := json.Marshal(&targetObject)
			So(err, ShouldBeNil)
			response := r.Post(targetURL, "application/json", string(raw))
			So(response, ShouldNotBeNil)
			if response.StatusCode != 200 {
				log.Println("ERROR: TestCRUDCommentaries, ", response.StatusCode)
			}
			So(response.StatusCode, ShouldEqual, 200)
			So(len(response.Body), ShouldBeGreaterThan, 0)

			//marshal to response

			responseObj := model.Response{}

			err = json.Unmarshal(response.RawBody, &responseObj)
			So(err, ShouldBeNil)
			So(responseObj.ID, ShouldNotBeBlank)

			//

			Convey("Then I Get this thing by ID", func() {

				targetURL = "/v1/read/thing/" + responseObj.ID + "/"
				log.Println(targetURL)

				response := r.Get(targetURL)
				So(response, ShouldNotBeNil)
				So(response.StatusCode, ShouldEqual, 200)
				So(len(response.Body), ShouldBeGreaterThan, 0)

				responseObj1 := []model.Thing{}

				err = json.Unmarshal(response.RawBody, &responseObj1)
				So(err, ShouldBeNil)

				So(len(responseObj1), ShouldBeGreaterThan, 0)

				//ASSERT Name
				So(responseObj1[0].Name, ShouldEqual, targetObject.Name)

				Convey("Then I Update this thing by ID", func() {

					targetURL = "/v1/update/thing/" + responseObj.ID + "/"

					//if bson.IsObjectIdHex(responseObj.ID) {
					//	targetObject.ID = bson.ObjectIdHex(responseObj.ID)
					//}

					//UPDATE FIELDS FOR FURTHER ASSERTION
					targetObject.Name = "test1234"

					raw, err := json.Marshal(&targetObject)
					So(err, ShouldBeNil)

					So(err, ShouldBeNil)
					response := r.Put(targetURL, "application/json", string(raw))
					So(response, ShouldNotBeNil)
					if response.StatusCode != 200 {
						log.Println("ERROR: -> ", response.StatusCode)
					}
					So(response.StatusCode, ShouldEqual, 200)

					Convey("Then I Get this thing by ID", func() {

						targetURL = "/v1/read/thing/" + responseObj.ID + "/"
						log.Println(targetURL)

						response := r.Get(targetURL)
						So(response, ShouldNotBeNil)
						So(response.StatusCode, ShouldEqual, 200)
						So(len(response.Body), ShouldBeGreaterThan, 0)

						responseObj := []model.Thing{}

						err = json.Unmarshal(response.RawBody, &responseObj)
						So(err, ShouldBeNil)
						So(len(responseObj), ShouldBeGreaterThan, 0)

						//ASSERT FIELDS YOU HAVE UPDATED HERE
						//TODO: Assert all the fields o/
						So(responseObj[0].ID, ShouldNotBeBlank)

						//ASSERT Lock
						So(responseObj[0].Name, ShouldEqual, targetObject.Name)

						Convey("Then I Get Delete thing by ID", func() {

							targetURL = "/v1/delete/thing/" + responseObj[0].ID + "/"
							log.Println(targetURL)

							response := r.Get(targetURL)
							So(response, ShouldNotBeNil)
							So(response.StatusCode, ShouldEqual, 200)
							So(len(response.Body), ShouldBeGreaterThan, 0)

						})
					})

				})

			})
		})
	})
}
