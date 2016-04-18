package pubsub

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/cloud"
	ps "google.golang.org/cloud/pubsub"
)

var (
	pubsubCtx  context.Context
	EMAIL      string
	CLIENT_KEY string
	projectID  = "curt-applications"
)

func NewContext() error {
	conf := &jwt.Config{
		Email:      EMAIL,
		PrivateKey: []byte(CLIENT_KEY),
		Scopes: []string{
			ps.ScopeCloudPlatform,
			ps.ScopePubSub,
		},
		TokenURL: google.JWTTokenURL,
	}

	pubsubCtx = cloud.NewContext(projectID, conf.Client(oauth2.NoContext))

	return nil
}

func PushMessage(topic string, msgs ...*ps.Message) error {
	var err error
	if pubsubCtx == nil {
		err = NewContext()
		if err != nil {
			return err
		}
	}

	err = createTopic(topic)
	if err != nil {
		return err
	}

	_, err = ps.Publish(pubsubCtx, topic, msgs...)
	return err

}

func createTopic(topic string) error {
	exists, err := ps.TopicExists(pubsubCtx, topic)
	if err != nil || exists {
		return err
	}

	return ps.CreateTopic(pubsubCtx, topic)
}
