package customer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"

	"gopkg.in/mgo.v2/bson"
)

var (
	pubsubContext context.Context
	ProjectID     = os.Getenv("PROJECT_ID")
	scopes        = []string{
		pubsub.ScopeCloudPlatform,
		pubsub.ScopePubSub,
	}
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

func NewPubSubContext() error {
	if ProjectID == "" {
		return fmt.Errorf("invalid ProjectID: %s", ProjectID)
	}

	if os.Getenv("GOOGLE_OAUTH_EMAIL") != "" {
		conf := &jwt.Config{
			Email:      os.Getenv("GOOGLE_OAUTH_EMAIL"),
			PrivateKey: []byte(os.Getenv("GOOGLE_CLIENT_KEY")),
			Scopes:     scopes,
			TokenURL:   google.JWTTokenURL,
		}
		pubsubContext = cloud.NewContext(ProjectID, conf.Client(oauth2.NoContext))
	} else {
		ctx := context.Background()
		c, err := google.DefaultClient(ctx, scopes...)
		if err != nil {
			return err
		}

		pubsubContext = cloud.WithContext(ctx, ProjectID, c)
	}

	return nil
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

	return publish(op, CustomerIndexTopic)
}

func publish(op IndexOperation, topic string) error {
	var err error
	if pubsubContext == nil {
		err = NewPubSubContext()
		if err != nil {
			return err
		}
	}

	err = createTopic(topic)
	if err != nil {
		return err
	}

	var msg pubsub.Message
	msg.Data, err = json.Marshal(&op)
	if err != nil {
		return err
	}

	pubsub.Publish(pubsubContext, topic, &msg)

	return nil
}

func createTopic(t string) error {
	exists, err := pubsub.TopicExists(pubsubContext, t)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return pubsub.CreateTopic(pubsubContext, t)
}
