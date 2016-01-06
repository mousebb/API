package customer

import (
	"log"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/curt-labs/API/helpers/database"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/twinj/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session *mgo.Session

	API_KEY       = uuid.NewV4()
	TEST_EMAIl    = "test@example.com"
	TEST_PASSWORD = "test_password"
)

func TestMain(m *testing.M) {

	mongo, err := dockertest.ConnectToMongoDB(15, time.Second, func(url string) bool {
		var err error
		session, err = mgo.Dial(url)
		if err != nil {
			log.Fatalf("MongoDB connection failed, with address '%s'.", url)
		}

		session.SetMode(mgo.Monotonic, true)

		encryptedPass, err := bcrypt.GenerateFromPassword([]byte(TEST_PASSWORD), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}

		c := Customer{
			Identifier: bson.NewObjectId(),
			Users: []User{
				User{
					Name:     "Test User",
					Email:    TEST_EMAIl,
					Password: string(encryptedPass),
					Keys: []APIKey{
						APIKey{
							Key: API_KEY.String(),
							Type: APIKeyType{
								Type: "Public",
							},
						},
					},
				},
			},
		}

		err = session.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName).Insert(&c)
		if err != nil {
			log.Fatal(err)
		}

		return session.Ping() == nil
	})

	defer func() {
		session.Close()
		mongo.KillRemove()
	}()

	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func TestGetUserByKey(t *testing.T) {
	Convey("Test GetUserByKey", t, func() {
		Convey("invalid mongo session", func() {
			tmp := database.CustomerCollectionName
			database.CustomerCollectionName = "example"

			user, err := GetUserByKey(session, API_KEY.String(), "Public")
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)

			database.CustomerCollectionName = tmp
		})
		Convey("valid key and type", func() {
			user, err := GetUserByKey(session, API_KEY.String(), "Public")
			So(err, ShouldBeNil)
			So(user, ShouldNotBeNil)
		})
	})
}

func TestAuthenticateUser(t *testing.T) {
	Convey("Test AuthenticateUser", t, func() {

		Convey("invalid mongo session", func() {
			tmp := database.CustomerCollectionName
			database.CustomerCollectionName = "example"

			user, err := AuthenticateUser(session, TEST_EMAIl, TEST_PASSWORD)
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)

			database.CustomerCollectionName = tmp
		})

		Convey("valid email and invalid password", func() {
			user, err := AuthenticateUser(session, TEST_EMAIl, "bad_password")
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)
		})

		Convey("valid email and password", func() {
			user, err := AuthenticateUser(session, TEST_EMAIl, TEST_PASSWORD)
			So(err, ShouldBeNil)
			So(user, ShouldNotBeNil)
		})
	})
}
