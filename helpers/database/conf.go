package database

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

// Scanner Provides the ability to scan SQL columns.
type Scanner interface {
	Scan(...interface{}) error
}

var (
	// EmptyDB Command line flag for testing
	EmptyDB = flag.String("clean", "", "bind empty database with structure defined")

	// ProductCollectionName Reference to MongoDB
	// for storing products.
	ProductCollectionName = "products"

	// CategoryCollectionName Reference to MongoDB
	// for storing categories.
	CategoryCollectionName = "categories"

	// CustomerCollectionName Reference to MongoDB
	// for storing customers.
	CustomerCollectionName = "customer"

	// ProductMongoSession MongoDB session for interacting
	// with products/category database.
	ProductMongoSession *mgo.Session

	// ProductMongoDatabase Name of product/category database.
	ProductMongoDatabase string

	// AriesMongoSession MongoDB session for interacting
	// with ARIES database.
	AriesMongoSession *mgo.Session

	// AriesMongoDatabase Name of ARIES application database.
	AriesMongoDatabase string

	// ErrorMongoSession MongoDB session for interacting
	// with errorDB database.
	ErrorMongoSession *mgo.Session

	// ErrorMongoDatabase Name of error application database.
	ErrorMongoDatabase string

	// DB Reference to MySQL database.
	DB *sql.DB

	// Driver SQL driver to use.
	Driver = "mysql"

	productDB = "product_data"
	ariesDB   = "aries"
	errorDB   = "errorDB"
)

// Init Initializes all database connections for SQL/NoSQL.
func Init() error {
	var err error
	if ProductMongoSession == nil {
		connectionString := mongoConnectionString(productDB)
		ProductMongoSession, err = mgo.DialWithInfo(connectionString)
		if err != nil {
			return err
		}
		ProductMongoDatabase = connectionString.Database
	}
	if AriesMongoSession == nil {
		connectionString := mongoConnectionString(ariesDB)
		AriesMongoSession, err = mgo.DialWithInfo(connectionString)
		if err != nil {
			return err
		}
		AriesMongoDatabase = connectionString.Database
	}
	if ErrorMongoSession == nil {
		connectionString := mongoConnectionString(errorDB)
		ErrorMongoSession, err = mgo.DialWithInfo(connectionString)
		if err != nil {
			return err
		}
		ErrorMongoDatabase = connectionString.Database
	}
	if DB == nil {
		DB, err = sql.Open(Driver, connectionString())
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() {
	AriesMongoSession.Close()
	ProductMongoSession.Close()
	ErrorMongoSession.Close()
	DB.Close()
}

func connectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("CURT_DEV_NAME")

		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	if EmptyDB != nil && *EmptyDB != "" {
		return "root:@tcp(127.0.0.1:3306)/CurtDev_Empty?parseTime=true&loc=America%2FChicago"
	}
	return "root:@tcp(127.0.0.1:3306)/CurtData?parseTime=true&loc=America%2FChicago"
}

// VcdbConnectionString Supplies connection string for the VCDB.
func VcdbConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("VCDB_NAME")
		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	return "root:@tcp(127.0.0.1:3306)/vcdb?parseTime=true&loc=America%2FChicago"
}

// VintelligencePass ...
func VintelligencePass() string {
	if vinPin := os.Getenv("VIN_PIN"); vinPin != "" {
		return fmt.Sprintf("%s", vinPin)
	}
	return "curtman:Oct2013!"
}

func mongoConnectionString(db string) *mgo.DialInfo {
	var info mgo.DialInfo
	addr := os.Getenv("MONGO_URL")
	if addr == "" {
		addr = "127.0.0.1"
	}
	addrs := strings.Split(addr, ",")
	info.Addrs = append(info.Addrs, addrs...)

	info.Username = os.Getenv("MONGO_CART_USERNAME")
	info.Password = os.Getenv("MONGO_CART_PASSWORD")
	info.Database = os.Getenv("MONGO_CART_DATABASE")
	info.Timeout = time.Second * 2
	info.FailFast = true
	info.Database = db
	info.Source = "admin"

	return &info
}

// GetCleanDBFlag ...
func GetCleanDBFlag() string {
	return *EmptyDB
}
