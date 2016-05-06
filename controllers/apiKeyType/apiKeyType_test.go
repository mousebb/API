package apiKeyType

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	db *sql.DB

	drops = map[string]string{
		`dropApiKeyType`: `DROP TABLE IF EXISTS ApiKeyType`,
	}

	schemas = map[string]string{
		`apiKeyTypeSchema`: `CREATE TABLE ApiKeyType (id varchar(64) NOT NULL,type varchar(500) DEFAULT NULL,date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
	}

	dataInserts = map[string]string{
		`insertKeyType`: `INSERT INTO ApiKeyType (id, type, date_added) VALUES (UUID(),'random type', NOW())`,
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
		db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/CurtData?parseTime=true")
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

		defer func() {
			db.Close()
		}()
	}

	m.Run()

}

func TestGetApiKeyTypes(t *testing.T) {

	Convey("Testing GetApiKeyTypes", t, func() {
		// ctx := &middleware.APIContext{
		// 	DataContext: &customer.DataContext{
		// 		BrandID: 3,
		// 	},
		// 	Params: httprouter.Params{},
		// 	DB:     db,
		// }

		Convey("with valid db connection", func() {
			// rec := httptest.NewRecorder()
			// req, err := http.NewRequest("GET", "http://localhost:8080/api/types", nil)
			// So(err, ShouldBeNil)

			// resp, err := GetApiKeyTypes(ctx, rec, req)
			// So(err, ShouldBeNil)
			// So(resp, ShouldHaveSameTypeAs, []apiKeyType.ApiKeyType{})
		})

	})
}
