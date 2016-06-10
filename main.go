package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bitly/go-nsq"
	"github.com/gorilla/mux"
)

var producer *nsq.Producer

//NSQstream is the stream name used in NSQ by Location Service
var NSQstream = "topic_location"

// Location is used to store a driver's location, as received by the driver
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DriverLocation contain the position of one driver at agiven time
type DriverLocation struct {
	DriverID int `json:"driverID"`
	Location
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	//NSQconnnection is the connection string to NSQ
	NSQconnnection := "172.17.0.1:4150"

	config := nsq.NewConfig()
	var err error
	producer, err = nsq.NewProducer(NSQconnnection, config)
	if err != nil {
		log.Fatal("Could not create a producer for nsq. Quit.")
	}
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/drivers/{id:[0-9]+}", ReceiveDriverLocation).Methods("PATCH")
	http.Handle("/", r)
	log.Printf("Server started and listening on port %d.", 3141)
	log.Fatal(http.ListenAndServe(":1337", nil))
	producer.Stop()
}

// ReceiveDriverLocation handles a driver's request to register its position at a given time
func ReceiveDriverLocation(w http.ResponseWriter, r *http.Request) {
	log.Printf("\t%s",
		r.RequestURI)
	// Read route parameter
	vars := mux.Vars(r)
	param := vars["id"]
	driverID, err := strconv.Atoi(param)
	if err != nil {
		log.Printf("Received bad request with driver id %s.", param)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var location Location
	if err := decoder.Decode(&location); err != nil {
		log.Printf("Unprocessable entity received. Err details: %s.", err.Error())
		// Return Unprocessable entity
		w.WriteHeader(422)
		return
	}

	if err := publish(driverID, &location); err != nil {
		log.Printf("Could not publish to NSQ.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func publish(driverID int, location *Location) error {

	publishedLocation := &DriverLocation{
		DriverID:  driverID,
		Location:  *location,
		UpdatedAt: time.Now(),
	}

	message, err := json.Marshal(publishedLocation)
	if err != nil {
		return err
	}

	if err := producer.Publish(NSQstream, message); err != nil {
		return err
	}
	return nil
}