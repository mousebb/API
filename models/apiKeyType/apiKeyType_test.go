package apiKeyType

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/curt-labs/API/middleware"
	"github.com/julienschmidt/httprouter"
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
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{},
			DB:     db,
		}

		Convey("with invalid DB query", func() {
			tmp := getAllKeyTypes
			getAllKeyTypes = "invalid database query"

			as, err := GetAllKeyTypes(ctx)
			So(err, ShouldNotBeNil)
			So(as, ShouldBeNil)

			getAllKeyTypes = tmp
		})

		Convey("with no data", func() {
			as, err := GetAllKeyTypes(ctx)
			So(err, ShouldBeNil)
			So(as, ShouldBeNil)
		})

		Convey("with missing select columns", func() {
			tmp := getAllKeyTypes
			getAllKeyTypes = "SELECT id, type FROM ApiKeyType order by type"

			err := insertKeys()
			So(err, ShouldBeNil)

			as, err := GetAllKeyTypes(ctx)
			So(err, ShouldNotBeNil)
			So(as, ShouldBeNil)

			getAllKeyTypes = tmp
		})

		Convey("with valid DB connection", func() {

			err := insertKeys()
			So(err, ShouldBeNil)

			as, err := GetAllKeyTypes(ctx)
			So(err, ShouldBeNil)
			So(len(as), ShouldBeGreaterThan, 0)

		})

	})
}

func BenchmarkGetAllKeyTypes(b *testing.B) {
	ctx := &middleware.APIContext{
		DataContext: &middleware.DataContext{
			BrandID: 3,
		},
		Params: httprouter.Params{},
		DB:     db,
	}

	for i := 0; i < b.N; i++ {
		GetAllKeyTypes(ctx)
	}
}
