package brandCtlr

import (
	// "bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/customer"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	db *sql.DB

	schemas = map[string]string{
		`dropBrand`:             `DROP TABLE IF EXISTS Brand`,
		`dropCustomerToBrand`:   `DROP TABLE IF EXISTS CustomerToBrand`,
		`dropWebsiteToBrand`:    `DROP TABLE IF EXISTS WebsiteToBrand`,
		`dropCustomer`:          `DROP TABLE IF EXISTS Customer`,
		`brandSchema`:           `CREATE TABLE Brand (ID int(11) NOT NULL AUTO_INCREMENT,name varchar(255) NOT NULL,code varchar(255) NOT NULL,logo varchar(255) DEFAULT NULL,logoAlt varchar(255) DEFAULT NULL,formalName varchar(255) DEFAULT NULL,longName varchar(255) DEFAULT NULL,primaryColor varchar(10) DEFAULT NULL,autocareID varchar(4) DEFAULT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`customerToBrandSchema`: `CREATE TABLE CustomerToBrand (ID int(11) NOT NULL AUTO_INCREMENT,cust_id int(11) NOT NULL,brandID int(11) NOT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=54486 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT`,
		`websiteToBrandSchema`:  `CREATE TABLE WebsiteToBrand (ID int(11) NOT NULL AUTO_INCREMENT,WebsiteID int(11) NOT NULL,brandID int(11) NOT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT`,
		`websiteSchema`:         `CREATE TABLE Website (ID int(11) NOT NULL AUTO_INCREMENT,url varchar(255) DEFAULT NULL,description varchar(255) DEFAULT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`customerSchema`:        `CREATE TABLE Customer (cust_id int(11) NOT NULL AUTO_INCREMENT,name varchar(255) DEFAULT NULL,email varchar(255) DEFAULT NULL,address varchar(500) DEFAULT NULL,city varchar(150) DEFAULT NULL,stateID int(11) DEFAULT NULL,phone varchar(50) DEFAULT NULL,fax varchar(50) DEFAULT NULL,contact_person varchar(300) DEFAULT NULL,dealer_type int(11) NOT NULL,latitude varchar(200) DEFAULT NULL,longitude varchar(200) DEFAULT NULL,password varchar(255) DEFAULT NULL,website varchar(500) DEFAULT NULL,customerID int(11) DEFAULT NULL,isDummy tinyint(1) NOT NULL DEFAULT '0',parentID int(11) DEFAULT NULL,searchURL varchar(500) DEFAULT NULL,eLocalURL varchar(500) DEFAULT NULL,logo varchar(500) DEFAULT NULL,address2 varchar(500) DEFAULT NULL,postal_code varchar(25) DEFAULT NULL,mCodeID int(11) NOT NULL DEFAULT '1',salesRepID int(11) DEFAULT NULL,APIKey varchar(64) DEFAULT NULL,tier int(11) NOT NULL DEFAULT '1',showWebsite tinyint(1) NOT NULL DEFAULT '0',PRIMARY KEY (cust_id)) ENGINE=InnoDB AUTO_INCREMENT=10444525 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
	}

	dataInserts = map[string]string{
		`insertBrand`: `insert into Brand(name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) values ('test brand', 'code','123','345','formal brand','long name','ffffff','auto')`,
	}
)

func TestMain(m *testing.M) {
	var err error
	if os.Getenv("DOCKER_BIND_LOCALHOST") == "" {
		var mysql dockertest.ContainerID
		mysql, err = dockertest.ConnectToMySQL(15, time.Second*5, func(url string) bool {
			db, err = sql.Open("mysql", url+"?parseTime=true")
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
	} else {
		db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/CurtData?parseTime=true")
		if err != nil {
			log.Fatalf("MySQL connection failed, with address '%s'.", "127.0.0.1:3306")
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

		defer func() {
			db.Close()
		}()
	}

	m.Run()

}

func TestGetAllBrands(t *testing.T) {

	Convey("Testing GetAllBrands", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &customer.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{},
			DB:     db,
		}

		Convey("with valid db connection", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/brands", nil)
			So(err, ShouldBeNil)

			resp, err := GetAllBrands(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []brand.Brand{})
		})

	})
}

func TestGetBrand(t *testing.T) {

	Convey("Testing GetBrand", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &customer.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{},
			DB:     db,
		}

		Convey("with invalid brand", func() {
			ctx.Params = httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "a",
				},
			}
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/brands/1", nil)
			So(err, ShouldBeNil)

			resp, err := GetBrand(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with valid brand", func() {

			ctx.Params = httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "1",
				},
			}
			recA := httptest.NewRecorder()
			reqA, err := http.NewRequest("GET", "http://localhost:8080/brands", nil)
			So(err, ShouldBeNil)

			resp, err := GetAllBrands(ctx, recA, reqA)
			So(err, ShouldBeNil)

			id := resp.([]brand.Brand)[0].ID
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/brands/%d", id), nil)
			So(err, ShouldBeNil)

			ctx.Params = httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: strconv.Itoa(id),
				},
			}

			resp, err = GetBrand(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, brand.Brand{})
		})

	})
}

func BenchmarkBrands(b *testing.B) {

}
