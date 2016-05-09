package brand

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	db *sql.DB

	drops = map[string]string{
		`dropBrand`:           `DROP TABLE IF EXISTS Brand`,
		`dropCustomerToBrand`: `DROP TABLE IF EXISTS CustomerToBrand`,
		`dropWebsiteToBrand`:  `DROP TABLE IF EXISTS WebsiteToBrand`,
		`dropWebsite`:         `DROP TABLE IF EXISTS Website`,
		`dropCustomer`:        `DROP TABLE IF EXISTS Customer`,
	}

	schemas = map[string]string{
		`brandSchema`:           `CREATE TABLE Brand (ID int(11) NOT NULL AUTO_INCREMENT,name varchar(255) NOT NULL,code varchar(255) NOT NULL,logo varchar(255) DEFAULT NULL,logoAlt varchar(255) DEFAULT NULL,formalName varchar(255) DEFAULT NULL,longName varchar(255) DEFAULT NULL,primaryColor varchar(10) DEFAULT NULL,autocareID varchar(4) DEFAULT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`customerToBrandSchema`: `CREATE TABLE CustomerToBrand (ID int(11) NOT NULL AUTO_INCREMENT,cust_id int(11) NOT NULL,brandID int(11) NOT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=54486 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT`,
		`websiteToBrandSchema`:  `CREATE TABLE WebsiteToBrand (ID int(11) NOT NULL AUTO_INCREMENT,WebsiteID int(11) NOT NULL,brandID int(11) NOT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT`,
		`websiteSchema`:         `CREATE TABLE Website (ID int(11) NOT NULL AUTO_INCREMENT,url varchar(255) DEFAULT NULL,description varchar(255) DEFAULT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`customerSchema`:        `CREATE TABLE Customer (cust_id int(11) NOT NULL AUTO_INCREMENT,name varchar(255) DEFAULT NULL,email varchar(255) DEFAULT NULL,address varchar(500) DEFAULT NULL,city varchar(150) DEFAULT NULL,stateID int(11) DEFAULT NULL,phone varchar(50) DEFAULT NULL,fax varchar(50) DEFAULT NULL,contact_person varchar(300) DEFAULT NULL,dealer_type int(11) NOT NULL,latitude varchar(200) DEFAULT NULL,longitude varchar(200) DEFAULT NULL,password varchar(255) DEFAULT NULL,website varchar(500) DEFAULT NULL,customerID int(11) DEFAULT NULL,isDummy tinyint(1) NOT NULL DEFAULT '0',parentID int(11) DEFAULT NULL,searchURL varchar(500) DEFAULT NULL,eLocalURL varchar(500) DEFAULT NULL,logo varchar(500) DEFAULT NULL,address2 varchar(500) DEFAULT NULL,postal_code varchar(25) DEFAULT NULL,mCodeID int(11) NOT NULL DEFAULT '1',salesRepID int(11) DEFAULT NULL,APIKey varchar(64) DEFAULT NULL,tier int(11) NOT NULL DEFAULT '1',showWebsite tinyint(1) NOT NULL DEFAULT '0',PRIMARY KEY (cust_id)) ENGINE=InnoDB AUTO_INCREMENT=10444525 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
	}

	dataInserts = map[string]string{
		`insertBrand`:           `insert into Brand(ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) values (1, 'test brand', 'code','123','345','formal brand','long name','ffffff','auto')`,
		`insertWebsite`:         `insert into Website(ID, url, description) values (1, 'www.website', 'test site')`,
		`insertWebsiteToBrand`:  `insert into WebsiteToBrand(ID, WebsiteID, brandID) values(1,1,1)`,
		`insertCustomer`:        `insert into Customer (cust_id, name, dealer_type) values (1, 'test', 1)`,
		`insertCustomerToBrand`: `insert into CustomerToBrand (ID, cust_id, brandID) values (1,1,1)`,
	}
)

func TestMain(m *testing.M) {
	var err error
	if os.Getenv("CI") == "" {
		var mysql dockertest.ContainerID
		mysql, err = dockertest.ConnectToMySQL(15, time.Second*5, func(url string) bool {
			db, err = sql.Open("mysql", url)
			if err != nil {
				log.Fatalf("MySQL connection failed, with address '%s'.", url)
			}

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

			return db.Ping() == nil
		})

		defer func() {
			db.Close()
			mysql.KillRemove()
		}()

		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = sql.Open(
			"mysql",
			fmt.Sprintf(
				"root:%s@tcp(%s:%s)%s?parseTime=true",
				os.Getenv("MARIADB_ENV_MYSQL_ROOT_PASSWORD"),
				os.Getenv("MARIADB_PORT_3306_TCP_ADDRESS"),
				os.Getenv("MARIADB_PORT_3306_TCP_PORT"),
				os.Getenv("MARIADB_NAME"),
			),
		)
		if err != nil {
			log.Fatalf("MySQL connection failed, with address '%s'.", "127.0.0.1:3306")
		}

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

		defer db.Close()
	}

	m.Run()

}

func TestGetAllBrands(t *testing.T) {

	Convey("Testing GetAllBrands", t, func() {

		Convey("with invalid db query", func() {
			tmp := getAllBrandsStmt
			getAllBrandsStmt = "invalid database query"

			brands, err := GetAllBrands(db)
			So(err, ShouldNotBeNil)
			So(brands, ShouldHaveSameTypeAs, []Brand{})

			getAllBrandsStmt = tmp
		})

		Convey("with missing select columns", func() {
			tmp := getAllBrandsStmt
			getAllBrandsStmt = strings.Replace(getAllBrandsStmt, "ID, ", "", 1)

			brands, err := GetAllBrands(db)
			So(err, ShouldNotBeNil)
			So(brands, ShouldHaveSameTypeAs, []Brand{})

			getAllBrandsStmt = tmp
		})

		Convey("with valid db connection", func() {

			brands, err := GetAllBrands(db)
			So(err, ShouldBeNil)
			So(brands, ShouldHaveSameTypeAs, []Brand{})
		})

	})
}

func TestGet(t *testing.T) {

	Convey("Testing Get", t, func() {

		Convey("with invalid db query", func() {
			tmp := getBrandStmt
			getBrandStmt = "invalid database query"

			b := Brand{}
			err := b.Get(db)
			So(err, ShouldNotBeNil)
			So(b.ID, ShouldEqual, 0)

			getBrandStmt = tmp
		})

		Convey("with missing select columns", func() {
			tmp := getBrandStmt
			getBrandStmt = strings.Replace(getBrandStmt, "ID, ", "", 1)

			b := Brand{
				ID: 1,
			}
			err := b.Get(db)
			So(err, ShouldNotBeNil)
			So(b.ID, ShouldEqual, 0)

			getBrandStmt = tmp
		})

		Convey("with valid db connection", func() {

			b := Brand{
				ID: 1,
			}

			err := b.Get(db)
			So(err, ShouldBeNil)
			So(b.Code, ShouldNotEqual, "")
		})

	})
}

func TestGetCustomerBrands(t *testing.T) {

	Convey("Testing GetCustomerBrands", t, func() {

		Convey("with invalid db query", func() {
			tmp := getCustomerBrands
			getCustomerBrands = "invalid database query"

			brands, err := GetCustomerBrands(0, db)
			So(err, ShouldNotBeNil)
			So(brands, ShouldHaveSameTypeAs, []Brand{})

			getCustomerBrands = tmp
		})

		Convey("with missing select columns", func() {
			tmp := getCustomerBrands
			getCustomerBrands = strings.Replace(getCustomerBrands, "ID, ", "", 1)

			brands, err := GetCustomerBrands(-1, db)
			So(err, ShouldNotBeNil)
			So(brands, ShouldHaveSameTypeAs, []Brand{})

			getCustomerBrands = tmp
		})

		Convey("with valid db connection", func() {

			brands, err := GetCustomerBrands(1, db)
			So(err, ShouldBeNil)
			So(brands, ShouldHaveSameTypeAs, []Brand{})
		})

	})
}

func TestGetWebsites(t *testing.T) {
	Convey("Testing with brand > 0", t, func() {
		ws, err := getWebsites(1, db)
		So(err, ShouldBeNil)
		So(len(ws), ShouldBeGreaterThan, 0)
	})
}
