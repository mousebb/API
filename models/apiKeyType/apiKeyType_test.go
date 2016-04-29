package apiKeyType

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	db *sql.DB

	schemas = map[string]string{
		`apiKeyTypeSchema`: `CREATE TABLE ApiKeyType (id varchar(64) NOT NULL,type varchar(500) DEFAULT NULL,date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
	}

	dataInserts = map[string]string{
		`insertKeyType`: `INSERT INTO ApiKeyType (id, type, date_added) VALUES (UUID(),'random type', NOW())`,
	}
)

func TestMain(m *testing.M) {

	mysql, err := dockertest.ConnectToMySQL(15, time.Second*5, func(url string) bool {
		var err error
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

		return db.Ping() == nil
	})

	defer func() {
		db.Close()
		mysql.KillRemove()
	}()

	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}

func insertKeys() error {
	for _, insert := range dataInserts {
		_, err := db.Exec(insert)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestGetAllKeyTypes(t *testing.T) {
	Convey("Test GetAllKeyTypes", t, func() {
		tx, err := db.Begin()
		So(err, ShouldBeNil)

		Convey("with invalid DB query", func() {
			tmp := GetAllTypes
			GetAllTypes = "invalid database query"

			as, err := GetAllKeyTypes(tx)
			So(err, ShouldNotBeNil)
			So(as, ShouldBeNil)

			GetAllTypes = tmp
		})

		Convey("with no data", func() {
			as, err := GetAllKeyTypes(tx)
			So(err, ShouldBeNil)
			So(as, ShouldBeNil)
		})

		Convey("with missing select columns", func() {
			tmp := GetAllTypes
			GetAllTypes = "SELECT id, type FROM ApiKeyType order by type"

			err := insertKeys()
			So(err, ShouldBeNil)

			as, err := GetAllKeyTypes(tx)
			So(err, ShouldNotBeNil)
			So(as, ShouldBeNil)

			GetAllTypes = tmp
		})

		Convey("with valid DB connection", func() {

			err := insertKeys()
			So(err, ShouldBeNil)

			as, err := GetAllKeyTypes(tx)
			So(err, ShouldBeNil)
			So(len(as), ShouldBeGreaterThan, 0)

		})

	})
}

func BenchmarkGetAllKeyTypes(b *testing.B) {
	tx, err := db.Begin()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		GetAllKeyTypes(tx)
	}
}
