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

type Scanner interface {
	Scan(...interface{}) error
}

var (
	EmptyDb = flag.String("clean", "", "bind empty database with structure defined")

	ProductCollectionName  = "products"
	CategoryCollectionName = "categories"
	productDB              = "product_data"
	ariesDB                = "aries"

	ProductMongoSession  *mgo.Session
	ProductMongoDatabase string

	AriesMongoSession  *mgo.Session
	AriesMongoDatabase string

	DB     *sql.DB
	Driver = "mysql"
)

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
	if DB == nil {
		DB, err = sql.Open(Driver, connectionString())
		if err != nil {
			return err
		}
	}

	return nil
}

func connectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("CURT_DEV_NAME")

		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	if EmptyDb != nil && *EmptyDb != "" {
		return "root:@tcp(127.0.0.1:3306)/CurtDev_Empty?parseTime=true&loc=America%2FChicago"
	}
	return "root:@tcp(127.0.0.1:3306)/CurtData?parseTime=true&loc=America%2FChicago"
}

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

func GetCleanDBFlag() string {
	return *EmptyDb
}
