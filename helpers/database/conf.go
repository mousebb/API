package database

import (
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
	ProductDatabase        = "product_data"
	CategoryDatabase       = "category_data"

	ProductMongoSession  *mgo.Session
	ProductMongoDatabase string
)

func Init() error {
	var err error
	if ProductMongoSession == nil {
		connectionString := ProductMongoConnectionString()
		ProductMongoSession, err = mgo.DialWithInfo(connectionString)
		if err != nil {
			return err
		}
		ProductMongoDatabase = connectionString.Database
	}

	return nil
}

func ConnectionString() string {
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

func MongoConnectionString() *mgo.DialInfo {
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
	if info.Database == "" {
		info.Database = "CurtCart"
	}
	info.Source = "admin"

	return &info
}

func ProductMongoConnectionString() *mgo.DialInfo {
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
	info.Database = "product_data"
	info.Source = "admin"

	return &info
}

func MongoPartConnectionString() *mgo.DialInfo {
	info := MongoConnectionString()
	info.Database = ProductDatabase
	return info
}

func AriesMongoConnectionString() *mgo.DialInfo {
	var info mgo.DialInfo
	addr := os.Getenv("MONGO_URL")
	if addr == "" {
		addr = "127.0.0.1"
	}
	addrs := strings.Split(addr, ",")
	info.Addrs = append(info.Addrs, addrs...)

	info.Username = os.Getenv("MONGO_ARIES_USERNAME")
	info.Password = os.Getenv("MONGO_ARIES_PASSWORD")
	info.Database = os.Getenv("MONGO_ARIES_DATABASE")
	info.Timeout = time.Second * 2
	info.FailFast = true
	if info.Database == "" {
		info.Database = "aries"
	}
	info.Source = "admin"

	return &info
}

func AriesConnectionString() string {
	return ConnectionString()
}

func GetCleanDBFlag() string {
	return *EmptyDb
}
