package geography

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
		`dropStates`:  `DROP TABLE IF EXISTS States`,
		`dropCountry`: `DROP TABLE IF EXISTS Country`,
	}

	schemas = map[string]string{
		`statesSchema`:  `CREATE TABLE States (stateID int(11) NOT NULL AUTO_INCREMENT,state varchar(100) NOT NULL,abbr varchar(3) NOT NULL,countryID int(11) NOT NULL DEFAULT '1',PRIMARY KEY (stateID)) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`countrySchema`: `CREATE TABLE Country (countryID int(11) NOT NULL AUTO_INCREMENT,name varchar(255) DEFAULT NULL,abbr varchar(10) DEFAULT NULL,PRIMARY KEY (countryID)) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
	}

	dataInserts = map[string]string{
		`insertCountry`: `INSERT INTO Country (name, abbr) VALUES('United States', 'USA')`,
		`insertState`:   `INSERT INTO States (state,abbr,countryID) VALUES('Wisconsin','WI',1)`,
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
}

func insertGeoData() error {
	for _, insert := range dataInserts {
		_, err := db.Exec(insert)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestMain(m *testing.M) {
	var err error
	if os.Getenv("CI") == "" {
		var mysql dockertest.ContainerID
		mysql, err = dockertest.ConnectToMySQL(15, time.Second*5, func(url string) bool {
			db, err = sql.Open("mysql", url+"?parseTime=true")
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
		defer db.Close()

		setupMySQL()
	}

	m.Run()
}

func TestGetAllCountriesAndstates(t *testing.T) {
	Convey("GetAllCountriesAndStates(db *sql.DB)", t, func() {
		Convey("with bad sql", func() {
			tmp := getAllCountriesAndStatesStmt
			getAllCountriesAndStatesStmt = "bad query"

			countries, err := GetAllCountriesAndStates(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldEqual, 0)
			So(err, ShouldNotBeNil)

			getAllCountriesAndStatesStmt = tmp
		})

		Convey("with bad query params", func() {
			tmp := getAllCountriesAndStatesStmt
			getAllCountriesAndStatesStmt = strings.Replace(getAllCountriesAndStatesStmt, "order by c.countryID", "where s.abbr = ? order by c.countryID", 1)

			countries, err := GetAllCountriesAndStates(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldEqual, 0)
			So(err, ShouldNotBeNil)

			getAllCountriesAndStatesStmt = tmp
		})

		Convey("with no data", func() {
			countries, err := GetAllCountriesAndStates(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldEqual, 0)
			So(err, ShouldBeNil)
		})

		Convey("with bad select columns", func() {

			err := insertGeoData()
			So(err, ShouldBeNil)

			tmp := getAllCountriesAndStatesStmt
			getAllCountriesAndStatesStmt = strings.Replace(getAllCountriesAndStatesStmt, "select c.countryID,", "select", 1)

			countries, err := GetAllCountriesAndStates(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldEqual, 0)
			So(err, ShouldBeNil)

			getAllCountriesAndStatesStmt = tmp
		})

		Convey("success", func() {
			err := insertGeoData()
			So(err, ShouldBeNil)
			countries, err := GetAllCountriesAndStates(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)
		})
	})
}

func TestGetAllCountries(t *testing.T) {
	Convey("GetAllCountries(db *sql.DB)", t, func() {
		Convey("with bad sql", func() {
			tmp := getAllCountriesStmt
			getAllCountriesStmt = "bad query"

			countries, err := GetAllCountries(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldEqual, 0)
			So(err, ShouldNotBeNil)

			getAllCountriesStmt = tmp
		})

		Convey("with bad query params", func() {
			tmp := getAllCountriesStmt
			getAllCountriesStmt = strings.Replace(getAllCountriesStmt, "order by c.countryID", "where c.abbr = ? order by c.countryID", 1)

			countries, err := GetAllCountries(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldEqual, 0)
			So(err, ShouldNotBeNil)

			getAllCountriesStmt = tmp
		})

		Convey("with bad select columns", func() {

			err := insertGeoData()
			So(err, ShouldBeNil)

			tmp := getAllCountriesStmt
			getAllCountriesStmt = strings.Replace(getAllCountriesStmt, "select c.countryID,", "select", 1)

			countries, err := GetAllCountries(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldEqual, 0)
			So(err, ShouldBeNil)

			getAllCountriesStmt = tmp
		})

		Convey("success", func() {
			countries, err := GetAllCountries(db)
			So(countries, ShouldHaveSameTypeAs, []Country{})
			So(len(countries), ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)
		})
	})
}

// func TestGeography(t *testing.T) {
// 	Convey("Testing Gets", t, func() {
// 		Convey("Testing GetAllCountriesAndStates()", func() {
// 			countrystates, err := GetAllCountriesAndStates(db)
// 			So(len(countrystates), ShouldBeGreaterThanOrEqualTo, 0)
// 			So(err, ShouldBeNil)
// 		})
//
// 		Convey("Testing GetAllCountries", func() {
// 			countries, err := GetAllCountries(db)
// 			So(len(countries), ShouldBeGreaterThanOrEqualTo, 0)
// 			So(err, ShouldBeNil)
// 		})
//
// 		Convey("Testing GetAllStates", func() {
// 			states, err := GetAllStates(db)
// 			So(len(states), ShouldBeGreaterThanOrEqualTo, 0)
// 			So(err, ShouldBeNil)
// 		})
// 	})
// }
//
// func BenchmarkGetAllCountriesAndStates(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		GetAllCountriesAndStates(db)
// 	}
// }
//
// func BenchmarkGetAllCountries(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		GetAllCountries(db)
// 	}
// }
//
// func BenchmarkGetAllStates(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		GetAllStates(db)
// 	}
// }
