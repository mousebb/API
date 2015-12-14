package apiKeyType

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/apiKeyType"
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

	m.Run()

}

func TestGetApiKeyTypes(t *testing.T) {

	Convey("Testing GetApiKeyTypes", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{},
			DB:     db,
		}

		Convey("with valid db connection", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/api/types", nil)
			So(err, ShouldBeNil)

			resp, err := GetApiKeyTypes(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []apiKeyType.ApiKeyType{})
		})

	})
}
