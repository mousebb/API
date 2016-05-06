package customer

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/apiKeyType"
	"github.com/curt-labs/API/models/brand"
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

	TestUserID                  = uuid.NewV4()
	TestUserAPIKey              = uuid.NewV4()
	TestUserPrivateAPIKey       = uuid.NewV4()
	TestSingleUserAPIKey        = uuid.NewV4()
	TestSingleUserPrivateAPIKey = uuid.NewV4()
	TestSuperUserAPIKey         = uuid.NewV4()
	TestSuperUserPrivateAPIKey  = uuid.NewV4()
	TestEmail                   = time.Now().String() + "_test@example.com"
	TestSuperEmail              = time.Now().String() + "_super@example.com"
	TestSingleUserEmail         = time.Now().String() + "_single@example.com"
	TestPassword                = "TestPassword"
	TestSuperPassword           = "TestSuperPassword"
	TestSingleUserPassword      = "single_password"

	drops = map[string]string{
		`dropBrand`:             `DROP TABLE IF EXISTS Brand`,
		`dropCustomerToBrand`:   `DROP TABLE IF EXISTS CustomerToBrand`,
		`dropCustomer`:          `DROP TABLE IF EXISTS Customer`,
		`dropCustomerUser`:      `DROP TABLE IF EXISTS CustomerUser`,
		`dropCustomerLocations`: `DROP TABLE IF EXISTS CustomerLocations`,
		`dropApiKeyType`:        `DROP TABLE IF EXISTS ApiKeyType`,
		`dropApiKey`:            `DROP TABLE IF EXISTS ApiKey`,
		`dropApiKeyToBrand`:     `DROP TABLE IF EXISTS ApiKeyToBrand`,
	}

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
		`insertBrand`:                     `insert into Brand(ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) values (1, 'test brand', 'code','123','345','formal brand','long name','ffffff','auto')`,
		`insertCustomer`:                  `insert into Customer (cust_id, name, dealer_type, customerID) values (1, 'test', 1, 1)`,
		`insertSingleUserCustomer`:        `insert into Customer (cust_id, name, dealer_type, customerID) values (2, 'test single user', 1, 2)`,
		`insertCustomerToBrand`:           `insert into CustomerToBrand (ID, cust_id, brandID) values (1,1,1)`,
		`insertSingleUserCustomerToBrand`: `insert into CustomerToBrand (ID, cust_id, brandID) values (2,2,3)`,
		`insertKeyType`:                   `insert into ApiKeyType (id, type, date_added) VALUES ('a46ceab9-df0c-44a2-b21d-a859fc2c839c','random type', NOW())`,
		`insertPrivateKeyType`:            `insert into ApiKeyType (id, type, date_added) VALUES (UUID(),'Private', NOW())`,
		`insertPublicKeyType`:             `insert into ApiKeyType (id, type, date_added) VALUES (UUID(),'Public', NOW())`,
	}
)

func setupMySQL() {
	var err error
	for _, schema := range drops {
		_, err = db.Exec(schema)
		if err != nil {
			log.Fatal(err)
		}
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
}

func setupMongo() {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(TestPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	encryptedSuperPass, err := bcrypt.GenerateFromPassword([]byte(TestSuperPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	encryptedSinglePass, err := bcrypt.GenerateFromPassword([]byte(TestSingleUserPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	c := Customer{
		Identifier:     bson.NewObjectId(),
		CustomerNumber: 1,
		Users: []User{
			User{
				ID:             TestUserID.String(),
				Name:           "Test User",
				Email:          TestEmail,
				Password:       string(encryptedPass),
				CustomerNumber: 1,
				Keys: []APIKey{
					APIKey{
						Key: TestUserAPIKey.String(),
						Type: apiKeyType.KeyType{
							Type: PublicKeyType,
						},
					},
					APIKey{
						Key: TestUserPrivateAPIKey.String(),
						Type: apiKeyType.KeyType{
							Type: PrivateKeyType,
						},
					},
				},
			},
			User{
				ID:             uuid.NewV4().String(),
				Name:           "Test Super User",
				Email:          TestSuperEmail,
				Password:       string(encryptedSuperPass),
				CustomerNumber: 1,
				SuperUser:      true,
				Keys: []APIKey{
					APIKey{
						Key: TestSuperUserAPIKey.String(),
						Type: apiKeyType.KeyType{
							Type: PublicKeyType,
						},
					},
					APIKey{
						Key: TestSuperUserPrivateAPIKey.String(),
						Type: apiKeyType.KeyType{
							Type: PrivateKeyType,
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

	singleUser := Customer{
		Identifier:     bson.NewObjectId(),
		CustomerNumber: 2,
		Users: []User{
			User{
				ID:             uuid.NewV4().String(),
				Name:           "Test Single User",
				Email:          TestSingleUserEmail,
				Password:       string(encryptedSinglePass),
				CustomerNumber: 2,
				SuperUser:      true,
				Keys: []APIKey{
					APIKey{
						Key: TestSingleUserAPIKey.String(),
						Type: apiKeyType.KeyType{
							Type: PublicKeyType,
						},
					},
					APIKey{
						Key: TestSingleUserPrivateAPIKey.String(),
						Type: apiKeyType.KeyType{
							Type: PrivateKeyType,
						},
					},
				},
			},
		},
	}

	err = session.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName).Insert(&singleUser)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	var err error
	if os.Getenv("DOCKER_BIND_LOCALHOST") == "" {
		var mysql dockertest.ContainerID
		mysql, err = dockertest.ConnectToMySQL(15, time.Second*5, func(url string) bool {
			url = fmt.Sprintf("%s?parseTime=true&loc=%s", url, "America%2FChicago")

			db, err = sql.Open("mysql", url)
			if err != nil {
				log.Fatalf("MySQL connection failed, with address '%s'.", url)
			}

			setupMySQL()

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

			setupMongo()

			return session.Ping() == nil
		})

		defer func() {
			session.Close()
			mongo.KillRemove()
		}()

		if err != nil {
			log.Fatal(err)
		}
	} else {
		// travis
		db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/CurtData?parseTime=true")
		if err != nil {
			log.Fatalf("MySQL connection failed, with address '%s'.", "127.0.0.1:3306")
		}

		setupMySQL()

		session, err = mgo.Dial("mongodb://127.0.0.1:27017/mydb")
		if err != nil {
			log.Fatal(err)
		}

		setupMongo()
	}

	m.Run()
}

func TestGetUserByKey(t *testing.T) {
	Convey("Test GetUserByKey", t, func() {
		Convey("invalid mongo session", func() {
			tmp := database.CustomerCollectionName
			database.CustomerCollectionName = "example"

			user, err := GetUserByKey(session, TestUserAPIKey.String(), PublicKeyType)
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)

			database.CustomerCollectionName = tmp
		})
		Convey("valid key and type", func() {
			user, err := GetUserByKey(session, TestUserAPIKey.String(), PublicKeyType)
			So(err, ShouldBeNil)
			So(user, ShouldNotBeNil)
		})

		Convey("valid key and no type", func() {
			user, err := GetUserByKey(session, TestUserAPIKey.String(), "")
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

			user, err := AuthenticateUser(session, TestEmail, TestPassword)
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)

			database.CustomerCollectionName = tmp
		})

		Convey("valid email and invalid password", func() {
			user, err := AuthenticateUser(session, TestEmail, "bad_password")
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)
		})

		Convey("valid email and password", func() {
			user, err := AuthenticateUser(session, TestEmail, TestPassword)
			So(err, ShouldBeNil)
			So(user, ShouldNotBeNil)
		})
	})
}

func TestGetUsers(t *testing.T) {
	Convey("GetUsers(*mgo.Session, string)", t, func() {
		Convey("with invalid requestor key", func() {
			users, err := GetUsers(session, "")
			So(users, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "failed to retrieve the requesting users information")
		})

		Convey("with valid requestor key that isn't super user", func() {
			users, err := GetUsers(session, TestUserPrivateAPIKey.String())
			So(users, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "this information is only available for super users")
		})

		Convey("customer with one user requesting all users", func() {
			users, err := GetUsers(session, TestSingleUserPrivateAPIKey.String())
			So(err, ShouldBeNil)
			So(len(users), ShouldEqual, 0)
		})

		Convey("valid", func() {
			users, err := GetUsers(session, TestSuperUserPrivateAPIKey.String())
			So(err, ShouldBeNil)
			So(len(users), ShouldBeGreaterThan, 0)
		})
	})
}

func TestGetUser(t *testing.T) {
	Convey("testing GetUser(*mgo.Session, string, string, string)", t, func() {

		Convey("with nil mgo.Session", func() {
			user, err := GetUser(nil, "", "")
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)
		})

		Convey("with empty userID", func() {
			user, err := GetUser(session, "", "")
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)
		})

		Convey("with empty requestorKey", func() {
			user, err := GetUser(session, "asd;lfjas;fd", "")
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)
		})

		Convey("with invalid requestorKey", func() {
			user, err := GetUser(session, uuid.NewV4().String(), uuid.NewV4().String())
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)
		})

		Convey("with requestorKey for non-super user", func() {
			user, err := GetUser(session, uuid.NewV4().String(), TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)
		})

		Convey("with invalid userID", func() {
			user, err := GetUser(session, uuid.NewV4().String(), TestSuperUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)
			So(user, ShouldBeNil)
		})

		Convey("valid", func() {
			user, err := GetUser(session, TestUserID.String(), TestSuperUserPrivateAPIKey.String())
			So(err, ShouldBeNil)
			So(user, ShouldNotBeNil)
		})
	})
}

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
			err := AddUser(session, db, nil, TestUserAPIKey.String())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "user object was null")
		})

		Convey("with no name and no email", func() {
			u := &User{}
			err := AddUser(session, db, u, TestUserAPIKey.String())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "name is required,e-mail is required")
		})

		Convey("with invalid requestor key", func() {
			u := &User{
				Name:  "Test User",
				Email: TestEmail,
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
				Email:    TestEmail,
				Password: "",
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)

			PasswordCharset = set
		})

		Convey("with bad database connection", func() {

			badDb, err := sql.Open("mysql", "localhost:3306")
			So(err, ShouldBeNil)

			u := &User{
				Name:     "Test User",
				Email:    TestEmail,
				Password: "",
			}
			err = AddUser(session, badDb, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)
		})

		Convey("with invalid location", func() {

			u := &User{
				Name:     "Test User",
				Email:    TestEmail,
				Password: "",
				Location: &Location{},
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)
		})

		Convey("with invalid insertUser query", func() {
			tmp := insertUser
			insertUser = "bad query"

			u := &User{
				Name:     "Test User",
				Email:    TestEmail,
				Password: "",
				Location: &validLocation,
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)

			insertUser = tmp
		})

		Convey("with invalid insert params", func() {

			tmp := insertUser
			insertUser = strings.Replace(insertUser, ", passwordConverted", "", 1)
			insertUser = strings.Replace(insertUser, ", 1)", ")", 1)

			u := &User{
				Name:     "Test User",
				Email:    TestEmail,
				Password: "",
				Location: &validLocation,
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)

			insertUser = tmp
		})

		Convey("with invalid getNewUserID query", func() {

			tmp := getNewUserID
			getNewUserID = "bad query"

			u := &User{
				Name:           "Test User",
				Email:          TestEmail,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)

			getNewUserID = tmp
		})

		Convey("with invalid where params in getNewUserID", func() {

			tmp := getNewUserID
			getNewUserID = strings.Replace(getNewUserID, "email = ? && ", "", 1)

			u := &User{
				Name:           "Test User",
				Email:          TestEmail,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)

			getNewUserID = tmp
		})

		Convey("invalid insertAPIKey query", func() {
			tmp := insertAPIKey
			insertAPIKey = "bad query"

			u := &User{
				Name:           "Test User",
				Email:          "valid" + TestEmail,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)

			insertAPIKey = tmp
		})

		Convey("invalid NSQ host", func() {

			tmp := NsqHost
			NsqHost = "1.2.3.4:4150"

			u := &User{
				Name:           "Test User",
				Email:          TestEmail,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)

			NsqHost = tmp

		})

		Convey("valid", func() {

			u := &User{
				Name:           "Test User",
				Email:          "valid" + TestEmail,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldBeNil)

		})

		Convey("duplicate e-mail", func() {

			u := &User{
				Name:           "Test User",
				Email:          "valid" + TestEmail,
				Password:       "",
				CustomerNumber: 1,
				Location:       &validLocation,
			}
			err := AddUser(session, db, u, TestUserPrivateAPIKey.String())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "there is a user registered with this e-mail")

		})

	})
}

func TestUpdateUser(t *testing.T) {
	Convey("Test UpdateUser(*sql.DB, *User)", t, func() {
		u := &User{
			Name:           "Test User",
			Email:          time.Now().String() + "update_valid" + TestEmail,
			Password:       "",
			CustomerNumber: 1,
			Location: &Location{
				Address: Address{
					StreetAddress: "1401 Meyer Rd",
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
				Email:         "test_meyer@example.com",
			},
		}
		e := AddUser(session, db, u, TestUserPrivateAPIKey.String())
		So(e, ShouldBeNil)

		Convey("with nil user", func() {
			err := UpdateUser(db, nil)
			So(err, ShouldNotBeNil)
		})

		Convey("invalid db connection", func() {
			err := UpdateUser(&sql.DB{}, nil)
			So(err, ShouldNotBeNil)
		})

		Convey("invalid location", func() {
			u.Location.ID = 0
			u.Location.Address = Address{}
			u.Location.Name = ""
			err := UpdateUser(db, u)
			So(err, ShouldNotBeNil)
		})

		Convey("invalid properties", func() {
			u.Name = ""
			err := UpdateUser(db, u)
			So(err, ShouldNotBeNil)

			u.Name = "Test User"
		})

		Convey("invalid updateUser query", func() {
			tmp := updateUser
			updateUser = "bad query"

			err := UpdateUser(db, u)
			So(err, ShouldNotBeNil)

			updateUser = tmp
		})

		Convey("invalid parameter count in updateUser query", func() {
			tmp := updateUser
			updateUser = strings.Replace(updateUser, "name = ?, ", "", 1)

			err := UpdateUser(db, u)
			So(err, ShouldNotBeNil)

			updateUser = tmp
		})

		Convey("valid", func() {

			err := UpdateUser(db, u)
			So(err, ShouldBeNil)
		})
	})
}

func TestValidate(t *testing.T) {
	Convey("validate()", t, func() {
		Convey("should fail and return two errors when empty", func() {
			var u User
			errs := u.validate(db)
			So(len(errs), ShouldEqual, 2)
			So(errs[0], ShouldEqual, "name is required")
			So(errs[1], ShouldEqual, "e-mail is required")
		})

		Convey("should fail on name with valid e-mail", func() {
			u := User{
				Email: "test@example.com",
			}
			errs := u.validate(db)
			So(len(errs), ShouldEqual, 1)
			So(errs[0], ShouldEqual, "name is required")
		})

		Convey("should fail on email with valid name", func() {
			u := User{
				Name: "Test User",
			}
			errs := u.validate(db)
			So(len(errs), ShouldEqual, 1)
			So(errs[0], ShouldEqual, "e-mail is required")
		})

		Convey("with invalid database connection", func() {
			u := User{
				Name:  "Test User",
				Email: "test@example.com",
			}
			errs := u.validate(nil)
			So(len(errs), ShouldEqual, 1)
		})

		Convey("with invalid checkForEmail query", func() {
			tmp := checkForEmail
			checkForEmail = strings.Replace(checkForEmail, " where email = ?", "", 1)
			u := User{
				Name:  "Test User",
				Email: "test@example.com",
			}
			errs := u.validate(db)
			So(len(errs), ShouldEqual, 1)

			checkForEmail = tmp
		})

		Convey("should pass", func() {
			u := User{
				Name:  "Test User",
				Email: "test@example.com",
			}
			errs := u.validate(db)
			So(len(errs), ShouldEqual, 0)
		})
	})
}

func TestStoreLocation(t *testing.T) {
	Convey("storeLocation(*sql.Tx)", t, func() {
		Convey("nil Location", func() {
			var tx *sql.Tx
			var u User
			err := u.storeLocation(tx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "location cannot be null")
		})

		Convey("invalid name", func() {
			var tx *sql.Tx
			u := User{
				Location: &Location{},
			}
			err := u.storeLocation(tx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid name")
		})

		Convey("invalid email", func() {
			var tx *sql.Tx
			u := User{
				Location: &Location{
					Name: "Test Location",
				},
			}
			err := u.storeLocation(tx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid email")
		})

		Convey("invalid phone", func() {
			var tx *sql.Tx
			u := User{
				Location: &Location{
					Name:  "Test Location",
					Email: "test@example.com",
				},
			}
			err := u.storeLocation(tx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid phone")
		})

		Convey("invalid city-state/postal", func() {
			var tx *sql.Tx
			u := User{
				Location: &Location{
					Name:  "Test Location",
					Email: "test@example.com",
					Phone: "515-555-4444",
				},
			}
			err := u.storeLocation(tx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid address")
		})

		Convey("invalid address", func() {
			var tx *sql.Tx
			u := User{
				Location: &Location{
					Name:  "Test Location",
					Email: "test@example.com",
					Phone: "515-555-4444",
					Address: Address{
						City: "Eau Claire",
						State: geography.State{
							State:        "Wisconsin",
							Abbreviation: "WI",
						},
					},
				},
			}
			err := u.storeLocation(tx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid address")
		})

		Convey("invalid address data", func() {
			var tx *sql.Tx
			u := User{
				Location: &Location{
					Name:  "Test Location",
					Email: "test@example.com",
					Phone: "515-555-4444",
					Address: Address{
						StreetAddress: "123 Test Street",
						City:          "Example",
						State: geography.State{
							State:        "Example State",
							Abbreviation: "ES",
						},
					},
				},
			}
			err := u.storeLocation(tx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "failed to get geospatial data")
		})

		Convey("invalid customer number", func() {

			tx, err := db.Begin()
			So(err, ShouldBeNil)

			u := User{
				Location: &Location{
					Name:  "Test Location",
					Email: "test@example.com",
					Phone: "515-555-4444",
					Address: Address{
						StreetAddress: "6208 Industrial Drive",
						City:          "Eau Claire",
						State: geography.State{
							State:        "Wisconsin",
							Abbreviation: "WI",
						},
					},
				},
			}

			err = u.storeLocation(tx)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "'cust_id' cannot be null")
		})

		Convey("valid", func() {

			tx, err := db.Begin()
			So(err, ShouldBeNil)

			u := User{
				CustomerNumber: 1,
				Location: &Location{
					Name:  "Test Location",
					Email: "test@example.com",
					Phone: "515-555-4444",
					Address: Address{
						StreetAddress: "6208 Industrial Drive",
						City:          "Eau Claire",
						State: geography.State{
							State:        "Wisconsin",
							Abbreviation: "WI",
						},
					},
				},
			}

			err = u.storeLocation(tx)
			So(err, ShouldBeNil)
		})
	})
}

func TestResetAuth(t *testing.T) {
	Convey("resetAuth()", t, func() {
		Convey("should return nil", func() {
			var u User
			err := u.resetAuth()
			So(err, ShouldBeNil)
		})
	})
}

func TestGenerateKeys(t *testing.T) {
	Convey("Test generateKeys(*sql.Tx, []brand.Brand)", t, func() {
		Convey("with bad transaction", func() {

			tmp := apiKeyType.GetAllTypes
			apiKeyType.GetAllTypes = "bad query"

			tx, err := db.Begin()
			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)

			u := User{}
			err = u.generateKeys(tx, []brand.Brand{})
			So(err, ShouldNotBeNil)

			apiKeyType.GetAllTypes = tmp
		})

		Convey("with limit 0 on key type query", func() {

			tmp := apiKeyType.GetAllTypes
			apiKeyType.GetAllTypes = fmt.Sprintf("%s limit 0", apiKeyType.GetAllTypes)

			tx, err := db.Begin()
			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)

			u := User{}
			err = u.generateKeys(tx, []brand.Brand{})
			So(err, ShouldNotBeNil)

			apiKeyType.GetAllTypes = tmp
		})

		Convey("with blank PrivateKeyType", func() {

			tmp := PrivateKeyType
			PrivateKeyType = ""

			tx, err := db.Begin()
			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)

			u := User{}
			err = u.generateKeys(tx, []brand.Brand{})
			So(err, ShouldNotBeNil)

			PrivateKeyType = tmp
		})

		Convey("with blank PublicKeyType", func() {

			tmp := PublicKeyType
			PublicKeyType = ""

			tx, err := db.Begin()
			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)

			u := User{}
			err = u.generateKeys(tx, []brand.Brand{})
			So(err, ShouldNotBeNil)

			PublicKeyType = tmp
		})

		Convey("with bad insertAPIKey", func() {

			tmp := insertAPIKey
			insertAPIKey = "bad query"

			tx, err := db.Begin()
			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)

			u := User{}
			err = u.generateKeys(tx, []brand.Brand{})
			So(err, ShouldNotBeNil)

			insertAPIKey = tmp
		})

		Convey("with invalid insert parameters", func() {

			tmp := insertAPIKey
			insertAPIKey = strings.Replace(insertAPIKey, "api_key, ", "", 1)
			insertAPIKey = strings.Replace(insertAPIKey, "(?, ", "(", 1)

			tx, err := db.Begin()
			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)

			u := User{}
			err = u.generateKeys(tx, []brand.Brand{})
			So(err, ShouldNotBeNil)

			insertAPIKey = tmp
		})

		Convey("valid", func() {

			tx, err := db.Begin()
			So(err, ShouldBeNil)
			So(tx, ShouldNotBeNil)

			u := User{}
			err = u.generateKeys(tx, []brand.Brand{})
			So(err, ShouldBeNil)
			So(len(u.Keys), ShouldEqual, 2)
		})

	})
}
