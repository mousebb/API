package customer

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/curt-labs/API/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type userResult struct {
	Users []User `bson:"users"`
}

// GetUserByKey Retrieves a User by using the APIKey associated
// with a User.
func GetUserByKey(sess *mgo.Session, key, t string) (*User, error) {
	var err error

	c := sess.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName)

	qry := bson.M{
		"users.keys.key": key,
	}
	if t != "" {
		qry["users.keys.type.type"] = t
	}

	var res userResult
	err = c.Find(qry).Select(bson.M{"_id": 0, "users.$": 1}).One(&res)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = fmt.Errorf("failed to locate user using %s %s", t, key)
		}
		return nil, err
	}

	return &res.Users[0], nil
}

func AuthenticateUser(sess *mgo.Session, email, pass string) (*User, error) {
	var err error

	c := sess.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName)

	qry := bson.M{
		"users.email": email,
	}

	var res userResult
	err = c.Find(qry).Select(bson.M{"_id": 0, "users.$": 1}).One(&res)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = fmt.Errorf("authentication failed for %s", email)
		}
		return nil, err
	}

	u := res.Users[0]
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass)) != nil {
		return nil, fmt.Errorf("authentication failed for %s", email)
	}

	return &u, nil
}
