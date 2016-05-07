package helpers

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

type DB struct {

}

var Session *mgo.Session

var connected bool

func ConnectToMongo() {
	if connected {
		fmt.Println("Error, already connected to MongoDB")
		return
	}
	var err error
	Session, err = mgo.Dial("")
	if err != nil {
		fmt.Println("ConnectToMongo dial error", err)
	}
}