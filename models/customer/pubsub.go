package customer

import (
	"database/sql"
	"log"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

const (
	CustomerIndexTopic        = "customer_index"
	INTRANSIT          Status = "In Transit"
)

type Status string

type IndexOperation struct {
	ModificationType string        `json:"modification_type"`
	ChangeType       string        `json:"change_type"`
	Identifier       string        `json:"id"`
	TransitID        bson.ObjectId `json:"_id" bson:"_id"`
	Status           Status        `json:"status" bson:"status"`
	Error            string        `json:"error" bson:"error"`
}

// PushCustomer Sends a message to NSQ for fanning out a given customer.
func PushCustomer(db *sql.DB, customerID int, change string, custOperationID bson.ObjectId) error {

	id, err := numberToID(db, customerID)
	if err != nil {
		return err
	}

	op := IndexOperation{
		TransitID:  bson.NewObjectId(),
		Status:     INTRANSIT,
		Identifier: strconv.Itoa(id),
		ChangeType: change,
	}

	log.Println(op)
	// return publish(op, CustomerIndexTopic)
	// return pubsub.PushMessage(CustomerIndexTopic, &msg)

	return nil
}
