package customer

import (
	"github.com/curt-labs/API/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetUserByKey Retrieves a User by using the APIKey associated
// with a User.
func GetUserByKey(sess *mgo.Session, key, t string) (*User, error) {
	var u User
	var err error

	c := sess.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName)

	err = c.Find(bson.M{"keys.$.key": key, "keys.$.key.type.type": t, "active": 1}).One(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
