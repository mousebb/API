package customer

import (
	"database/sql"
	"log"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/geography"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/twinj/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	session *mgo.Session
	db      *sql.DB

	API_KEY         = uuid.NewV4()
	PRIVATE_API_KEY = uuid.NewV4()
	TEST_EMAIl      = "test@example.com"
	TEST_PASSWORD   = "test_password"

	schemas = map[string]string{
		`apiKeySchema`:           `CREATE TABLE ApiKey (id int(11) NOT NULL AUTO_INCREMENT,api_key varchar(64) NOT NULL,type_id varchar(64) NOT NULL,user_id varchar(64) NOT NULL,date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,UNIQUE KEY id (id)) ENGINE=InnoDB AUTO_INCREMENT=14489 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`apiKeyBrandSchema`:      `CREATE TABLE ApiKeyToBrand (ID int(11) NOT NULL AUTO_INCREMENT,keyID int(11) NOT NULL,brandID int(11) NOT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=38361 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT`,
		`apiKeyTypeSchema`:       `CREATE TABLE ApiKeyType (id varchar(64) NOT NULL,type varchar(500) DEFAULT NULL,date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`brandSchema`:            `CREATE TABLE Brand (ID int(11) NOT NULL AUTO_INCREMENT,name varchar(255) NOT NULL,code varchar(255) NOT NULL,logo varchar(255) DEFAULT NULL,logoAlt varchar(255) DEFAULT NULL,formalName varchar(255) DEFAULT NULL,longName varchar(255) DEFAULT NULL,primaryColor varchar(10) DEFAULT NULL,autocareID varchar(4) DEFAULT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`customerToBrandSchema`:  `CREATE TABLE CustomerToBrand (ID int(11) NOT NULL AUTO_INCREMENT,cust_id int(11) NOT NULL,brandID int(11) NOT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=54486 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT`,
		`customerSchema`:         `CREATE TABLE Customer (cust_id int(11) NOT NULL AUTO_INCREMENT,name varchar(255) DEFAULT NULL,email varchar(255) DEFAULT NULL,address varchar(500) DEFAULT NULL,city varchar(150) DEFAULT NULL,stateID int(11) DEFAULT NULL,phone varchar(50) DEFAULT NULL,fax varchar(50) DEFAULT NULL,contact_person varchar(300) DEFAULT NULL,dealer_type int(11) NOT NULL,latitude varchar(200) DEFAULT NULL,longitude varchar(200) DEFAULT NULL,password varchar(255) DEFAULT NULL,website varchar(500) DEFAULT NULL,customerID int(11) DEFAULT NULL,isDummy tinyint(1) NOT NULL DEFAULT '0',parentID int(11) DEFAULT NULL,searchURL varchar(500) DEFAULT NULL,eLocalURL varchar(500) DEFAULT NULL,logo varchar(500) DEFAULT NULL,address2 varchar(500) DEFAULT NULL,postal_code varchar(25) DEFAULT NULL,mCodeID int(11) NOT NULL DEFAULT '1',salesRepID int(11) DEFAULT NULL,APIKey varchar(64) DEFAULT NULL,tier int(11) NOT NULL DEFAULT '1',showWebsite tinyint(1) NOT NULL DEFAULT '0',PRIMARY KEY (cust_id)) ENGINE=InnoDB AUTO_INCREMENT=10444525 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`customerUserSchema`:     `CREATE TABLE CustomerUser (id varchar(64) NOT NULL,name varchar(255) DEFAULT NULL,email varchar(255) NOT NULL,password varchar(255) NOT NULL,customerID int(11) DEFAULT NULL,date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,active tinyint(1) NOT NULL DEFAULT '0',locationID int(11) NOT NULL DEFAULT '0',isSudo tinyint(1) NOT NULL DEFAULT '0',cust_ID int(11) NOT NULL,NotCustomer tinyint(1) DEFAULT NULL,passwordConverted tinyint(1) NOT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`customerLocationSchema`: `CREATE TABLE CustomerLocations (locationID int(11) NOT NULL AUTO_INCREMENT,name varchar(500) DEFAULT NULL,address varchar(500) DEFAULT NULL,city varchar(500) DEFAULT NULL,stateID int(11) NOT NULL,email varchar(500) DEFAULT NULL,phone varchar(20) DEFAULT NULL,fax varchar(20) DEFAULT NULL,latitude double NOT NULL,longitude double NOT NULL,cust_id int(11) NOT NULL DEFAULT '0',contact_person varchar(300) DEFAULT NULL,isprimary tinyint(1) NOT NULL DEFAULT '0',postalCode varchar(30) DEFAULT NULL,ShippingDefault tinyint(1) NOT NULL DEFAULT '0',PRIMARY KEY (locationID)) ENGINE=InnoDB AUTO_INCREMENT=11001 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
	}

	dataInserts = map[string]string{
		`insertBrand`:           `insert into Brand(ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) values (1, 'test brand', 'code','123','345','formal brand','long name','ffffff','auto')`,
		`insertCustomer`:        `insert into Customer (cust_id, name, dealer_type, customerID) values (1, 'test', 1, 1)`,
		`insertCustomerToBrand`: `insert into CustomerToBrand (ID, cust_id, brandID) values (1,1,1)`,
		`insertKeyType`:         `INSERT INTO ApiKeyType (id, type, date_added) VALUES ('a46ceab9-df0c-44a2-b21d-a859fc2c839c','random type', NOW())`,
	}
)

func TestMain(m *testing.M) {

	mysql, err := dockertest.ConnectToMySQL(15, time.Second*5, func(url string) bool {
		var err error
		db, err = sql.Open("mysql", url)
		if err != nil {
			log.Fatalf("MySQL connection failed, with address '%s'.", url)
		}

		for _, schema := range schemas {
			_, err = db.Exec(schema)
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, insert := range dataInserts {
			_, err = db.Exec(insert)
			if err != nil {
				log.Fatal(err)
			}
		}

		return db.Ping() == nil
	})

	defer func() {
		db.Close()
		mysql.KillRemove()
	}()

	if err != nil {
		log.Fatal(err)
	}

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
					Name:           "Test User",
					Email:          TEST_EMAIl,
					Password:       string(encryptedPass),
					CustomerNumber: 1,
					Keys: []APIKey{
						APIKey{
							Key: API_KEY.String(),
							Type: APIKeyType{
								Type: "Public",
							},
						},
						APIKey{
							Key: PRIVATE_API_KEY.String(),
							Type: APIKeyType{
								Type: "Private",
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

// func TestGetUserByKey(t *testing.T) {
// 	Convey("Test GetUserByKey", t, func() {
// 		Convey("invalid mongo session", func() {
// 			tmp := database.CustomerCollectionName
// 			database.CustomerCollectionName = "example"
//
// 			user, err := GetUserByKey(session, API_KEY.String(), "Public")
// 			So(err, ShouldNotBeNil)
// 			So(user, ShouldBeNil)
//
// 			database.CustomerCollectionName = tmp
// 		})
// 		Convey("valid key and type", func() {
// 			user, err := GetUserByKey(session, API_KEY.String(), "Public")
// 			So(err, ShouldBeNil)
// 			So(user, ShouldNotBeNil)
// 		})
// 	})
// }
//
// func TestAuthenticateUser(t *testing.T) {
// 	Convey("Test AuthenticateUser", t, func() {
//
// 		Convey("invalid mongo session", func() {
// 			tmp := database.CustomerCollectionName
// 			database.CustomerCollectionName = "example"
//
// 			user, err := AuthenticateUser(session, TEST_EMAIl, TEST_PASSWORD)
// 			So(err, ShouldNotBeNil)
// 			So(user, ShouldBeNil)
//
// 			database.CustomerCollectionName = tmp
// 		})
//
// 		Convey("valid email and invalid password", func() {
// 			user, err := AuthenticateUser(session, TEST_EMAIl, "bad_password")
// 			So(err, ShouldNotBeNil)
// 			So(user, ShouldBeNil)
// 		})
//
// 		Convey("valid email and password", func() {
// 			user, err := AuthenticateUser(session, TEST_EMAIl, TEST_PASSWORD)
// 			So(err, ShouldBeNil)
// 			So(user, ShouldNotBeNil)
// 		})
// 	})
// }
//
// func TestAuthenticateUserByKey(t *testing.T) {
// 	Convey("Test AuthenticateUserByKey", t, func() {
//
// 		Convey("invalid mongo session", func() {
// 			tmp := database.CustomerCollectionName
// 			database.CustomerCollectionName = "example"
//
// 			user, err := AuthenticateUserByKey(session, PRIVATE_API_KEY.String())
// 			So(err, ShouldNotBeNil)
// 			So(user, ShouldBeNil)
//
// 			database.CustomerCollectionName = tmp
// 		})
//
// 		Convey("invalid key", func() {
// 			user, err := AuthenticateUserByKey(session, uuid.NewV4().String())
// 			So(err, ShouldNotBeNil)
// 			So(user, ShouldBeNil)
// 		})
//
// 		Convey("valid key", func() {
// 			user, err := AuthenticateUserByKey(session, PRIVATE_API_KEY.String())
// 			So(err, ShouldBeNil)
// 			So(user, ShouldNotBeNil)
// 		})
// 	})
// }

func TestAddUser(t *testing.T) {
	Convey("Test AddUser", t, func() {

		validLocation := Location{
			Address: Address{
				StreetAddress: "6208 Industrial Drive",
				City:          "Eau Claire",
				PostalCode:    "54701",
				State: geography.State{
					Abbreviation: "WI",
					State:        "Wisconsin",
					Country: &geography.Country{
						Abbreviation: "USA",
						Country:      "United States of Ameria",
					},
				},
			},
			ContactPerson: "Test User",
			Phone:         "7155555555",
			Fax:           "7154444444",
			Name:          "Test Location",
			Email:         "test@example.com",
		}

		Convey("with nil User", func() {
			err := AddUser(session, db, nil, API_KEY.String())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "user object was null")
		})

		Convey("with no name and no email", func() {
			u := &User{}
			err := AddUser(session, db, u, API_KEY.String())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "name is required,e-mail is required")
		})

		Convey("with invalid requestor key", func() {
			u := &User{
				Name:  "Test User",
				Email: TEST_EMAIl,
			}
			err := AddUser(session, db, u, "")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "failed to retrieve the requesting users information")
		})

		Convey("with valid requestor key and empty password", func() {
			set := PasswordCharset
			PasswordCharset = " "
			u := &User{
				Name:     "Test User",
				Email:    TEST_EMAIl,
				Password: "",
			}
			err := AddUser(session, db, u, PRIVATE_API_KEY.String())
			So(err, ShouldNotBeNil)

			PasswordCharset = set
		})

		Convey("with bad database connection", func() {

			badDb, err := sql.Open("mysql", "localhost:3306")
			So(err, ShouldBeNil)

			u := &User{
				Name:     "Test User",
				Email:    TEST_EMAIl,
				Password: "",
			}
			err = AddUser(session, badDb, u, PRIVATE_API_KEY.String())
			So(err, ShouldNotBeNil)
		})

		Convey("with invalid location", func() {

			u := &User{
				Name:     "Test User",
				Email:    TEST_EMAIl,
				Password: "",
				Location: &Location{},
			}
			err := AddUser(session, db, u, PRIVATE_API_KEY.String())
			So(err, ShouldNotBeNil)
		})

		Convey("with invalid insertUser query", func() {
			tmp := insertUser
			insertUser = "bad query"

			u := &User{
				Name:     "Test User",
				Email:    TEST_EMAIl,
				Password: "",
				Location: &validLocation,
			}
			err := AddUser(session, db, u, PRIVATE_API_KEY.String())
			So(err, ShouldNotBeNil)

			insertUser = tmp
		})

		Convey("with invalid insert params", func() {

			tmp := insertUser
			insertUser = strings.Replace(insertUser, ", passwordConverted", "", 1)
			insertUser = strings.Replace(insertUser, ", 1)", ")", 1)

			u := &User{
				Name:     "Test User",
				Email:    TEST_EMAIl,
				Password: "",
				Location: &validLocation,
			}
			err := AddUser(session, db, u, PRIVATE_API_KEY.String())
			So(err, ShouldNotBeNil)

			insertUser = tmp
		})

		Convey("with invalid getNewUserID query", func() {

			tmp := getNewUserID
			getNewUserID = "bad query"

			u := &User{
				Name:           "Test User",
				Email:          TEST_EMAIl,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, PRIVATE_API_KEY.String())
			So(err, ShouldNotBeNil)

			getNewUserID = tmp
		})

		Convey("with invalid where params in getNewUserID", func() {

			tmp := getNewUserID
			getNewUserID = strings.Replace(getNewUserID, "email = ? && ", "", 1)

			u := &User{
				Name:           "Test User",
				Email:          TEST_EMAIl,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, PRIVATE_API_KEY.String())
			So(err, ShouldNotBeNil)

			getNewUserID = tmp
		})

		Convey("valid", func() {

			u := &User{
				Name:           "Test User",
				Email:          TEST_EMAIl,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, PRIVATE_API_KEY.String())
			So(err, ShouldBeNil)

		})

	})
}
