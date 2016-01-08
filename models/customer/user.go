package customer

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/curt-labs/API/helpers/database"
	"github.com/jmcvetta/randutil"
	"github.com/kellydunn/golang-geo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	insertUser = `insert into CustomerUser(id, name, email, password, customerID, date_added, active, locationID, isSudo, cust_ID, NotCustomer, passwordConverted)
					values(UUID(),?, ?, ?, ?, ?, 1, ?, ?, (
						select cust_id from Customer where customerID = ? limit 1
					), 1, 1)`
	getNewUserID = `select id from CustomerUser where email = ? && password = ? limit 1`

	// GeocodingAPIKey API Key for Google Maps Geocoding API.
	GeocodingAPIKey = `AIzaSyAKINnVskCaZkQhhh6I2D6DOpeylY1G1-Q`
	// PasswordCharset The character set that we want to use for
	// generating random passwords.
	PasswordCharset = `ABCDEFGHJKMNPQRTUVWXYZabcdefghijkmnpqrtuvwxyz12346789`
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

// AuthenticateUser Retrieves a User based of an email/password authentication model.
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

	// checkout the TODO on resetAuth() func definition
	// err = u.resetAuth()
	// if err != nil {
	// 	return nil, err
	// }

	return &u, nil
}

// AddUser Will commit a new user to the same Customer object as
// the requestor's Customer reference. It will not update the following
// fields from the submitted User object: `ID`, `CustomerNumber`, `DateAdded`, `Keys`, or `ComnetAccounts`.
//
// Required fields are: `Name`, `Email`.
//
// `Password` may be supplied, if not supplied it will be randomly generated.
//
// Order of Operation:
// 1. Validate required fields
// 2. Retrieve `CustomerNumber` from `requestorKey`
// 3. Determine if `Passsword` generation is necessary
// 4. Add to MySQL
// 5. Generate geospatial data for `Location`, if provided.
// 6. Generate `DateAdded` timestamp
// 7. Call fanner process for the associated `Customer`
func AddUser(sess *mgo.Session, db *sql.DB, user *User, requestorKey string) error {

	if user == nil {
		return fmt.Errorf("user object was null")
	}

	// validate
	errors := user.validate()
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ","))
	}

	// fetch requestor
	requestor, err := GetUserByKey(sess, requestorKey, "Private")
	if err != nil || requestor.CustomerNumber == 0 {
		return fmt.Errorf("failed to retrieve the requesting users information %v", err)
	}

	// set customer number from requestor
	user.CustomerNumber = requestor.CustomerNumber

	if user.Password == "" {
		pass, err := randutil.String(8, PasswordCharset)
		if err != nil || strings.TrimSpace(pass) == "" {
			if err == nil {
				err = fmt.Errorf("generated password was empty")
			}
			return fmt.Errorf("failed to generate password %s", err.Error())
		}

		enc, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to generate password %s", err.Error())
		}

		user.Password = string(enc)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if user.Location != nil {
		exists, err := user.Location.Exists(db, user.CustomerNumber)
		if !exists && err == nil {
			err = user.storeLocation(tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	user.DateAdded = time.Now()

	stmt, err := tx.Prepare(insertUser)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		user.Name,
		user.Email,
		user.Password,
		user.CustomerNumber,
		user.DateAdded,
		user.Location.ID,
		user.SuperUser,
		user.CustomerNumber,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err = tx.Prepare(getNewUserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	var userID *string
	err = stmt.QueryRow(user.Email, user.Password).Scan(&userID)
	if err != nil || userID == nil || *userID == "" {
		tx.Rollback()
		return fmt.Errorf("failed to insert user %s", err.Error())
	}

	user.ID = *userID

	err = tx.Commit()
	if err != nil {
		return err
	}

	return PushCustomer(db, user.CustomerNumber, "update", bson.NewObjectId())
}

func (user User) validate() []string {
	var errs []string
	if user.Name == "" {
		errs = append(errs, "name is required")
	}

	if user.Email == "" {
		errs = append(errs, "e-mail is required")
	}

	return errs
}

func (user User) storeLocation(tx *sql.Tx) error {
	if user.Location.Address.City == "" || (user.Location.Address.State.State == "" && user.Location.Address.State.Abbreviation == "") {
		if user.Location.Address.PostalCode == "" {
			return fmt.Errorf("invalid address")
		}
	}
	if user.Location.Address.StreetAddress == "" {
		return fmt.Errorf("invalid address")
	}

	// get geospatial data
	geo.SetGoogleAPIKey(GeocodingAPIKey)
	coder := geo.GoogleGeocoder{}
	point, err := coder.Geocode(
		fmt.Sprintf(
			"%s %s %s %s %s",
			user.Location.Address.StreetAddress,
			user.Location.Address.StreetAddress2,
			user.Location.Address.City,
			user.Location.Address.State.Abbreviation,
			user.Location.Address.PostalCode,
		),
	)
	if err != nil || point == nil {
		return fmt.Errorf("failed to get geospatial data %s", err.Error())
	}
	user.Location.Address.Coordinates = Coordinates{
		Latitude:  point.Lat(),
		Longitude: point.Lng(),
	}

	if user.Location.Name == "" {
		return fmt.Errorf("invalid name")
	}
	if user.Location.Email == "" {
		return fmt.Errorf("invalid email")
	}
	if user.Location.Phone == "" {
		return fmt.Errorf("invalid phone")
	}

	user.Location.insert(tx, user.CustomerNumber)

	return nil
}

// AuthenticateUserByKey Allows a user to authenticate using the `Private` APIKey.
// If a user tries to authenticate against this with their `Public` APIKey,
// they should fail, since the `Public` APIKey is something you would embed
// in public eyes.
func AuthenticateUserByKey(sess *mgo.Session, key string) (*User, error) {
	var err error

	c := sess.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName)

	qry := bson.M{
		"users.keys.key":       key,
		"users.keys.type.type": "Private",
	}

	var res userResult
	err = c.Find(qry).Select(bson.M{"_id": 0, "users.$": 1}).One(&res)
	if err != nil {
		if err == mgo.ErrNotFound {
			err = fmt.Errorf("authentication failed for %s", key)
		}
		return nil, err
	}

	u := res.Users[0]

	// checkout the TODO on resetAuth() func definition
	// err = u.resetAuth()
	// if err != nil {
	// 	return nil, err
	// }

	return &u, nil
}

// TODO: I'd like to bring up the argument that
// we don't need this if we leverage other API keys
// properly. If this does end up being implemented remember
// the following Order of Operation
//
// 1. Generate new UUID
// 2. Update Authentication Key in MySQL
// 3. Fan that customer by making call to fanner.
func (user *User) resetAuth() error {

	return nil
}
