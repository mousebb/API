package customer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	nsq "github.com/bitly/go-nsq"
	"gopkg.in/mgo.v2/bson"
)

var (
	// NsqHost Represents a remote IP address for our NSQ host.
	NsqHost = os.Getenv("NSQ_HOST")
)

type nopLogger struct{}

type message struct {
	ModificationType string        `json:"modification_type"`
	ChangeType       string        `json:"change_type"`
	Identifier       string        `json:"id"`
	TransitID        bson.ObjectId `json:"_id" bson:"_id"`
	Error            error         `json:"error" bson:"error"`
}

func (*nopLogger) Output(int, string) error {
	return nil
}

// PushCustomer Sends a message to NSQ for fanning out a given customer.
func PushCustomer(db *sql.DB, customerID int, change string, custOperationID bson.ObjectId) error {

	id, err := numberToID(db, customerID)
	if err != nil {
		return err
	}

	config := nsq.NewConfig()
	w, err := nsq.NewProducer(getDaemonHosts(), config)
	if w == nil && err == nil {
		return fmt.Errorf("%s", "failed to connect to producer")
	}
	w.SetLogger(&nopLogger{}, nsq.LogLevelError)
	defer w.Stop()

	if err != nil {
		return err
	}

	pm := message{
		ModificationType: "Customer",
		ChangeType:       change,
		Identifier:       string(id),
		TransitID:        custOperationID,
	}

	js, err := json.Marshal(pm)
	if err != nil {
		return nil
	}

	err = w.Publish("admin_change", js)
	if err != nil {
		return err
	}

	return nil
}

func getDaemonHosts() string {
	if NsqHost == "" {
		return "192.168.99.100:4150"
	}
	return NsqHost
}
