package applicationGuide

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/applicationGuide"
	"github.com/curt-labs/API/models/customer"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	Key = "f50452b9-2fb0-4008-864f-b2c58359d3ad"
)

var (
	db *sql.DB

	drops = map[string]string{
		`dropApiKeyType`:        `DROP TABLE IF EXISTS ApiKeyType`,
		`dropApiKey`:            `DROP TABLE IF EXISTS ApiKey`,
		`dropApiKeyToBrand`:     `DROP TABLE IF EXISTS ApiKeyToBrand`,
		`dropCategories`:        `DROP TABLE IF EXISTS Categories`,
		`dropApplicationGuides`: `DROP TABLE IF EXISTS ApplicationGuides`,
	}

	schemas = map[string]string{
		`apiKeyTypeSchema`:       `CREATE TABLE ApiKeyType (id varchar(64) NOT NULL,type varchar(500) DEFAULT NULL,date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`categorySchema`:         `CREATE TABLE Categories (catID int(11) NOT NULL AUTO_INCREMENT,dateAdded timestamp NULL DEFAULT CURRENT_TIMESTAMP,parentID int(11) NOT NULL,catTitle varchar(100) DEFAULT NULL,shortDesc varchar(255) DEFAULT NULL,longDesc longtext,image varchar(255) DEFAULT NULL,isLifestyle int(11) NOT NULL,codeID int(11) NOT NULL DEFAULT '0',sort int(11) NOT NULL DEFAULT '1',vehicleSpecific tinyint(1) NOT NULL DEFAULT '0',vehicleRequired tinyint(1) NOT NULL DEFAULT '0',metaTitle text,metaDesc text,metaKeywords text,icon varchar(255) DEFAULT NULL,path varchar(255) DEFAULT NULL,brandID int(11) NOT NULL DEFAULT '1',isDeleted tinyint(1) NOT NULL DEFAULT '0',PRIMARY KEY (catID)) ENGINE=InnoDB AUTO_INCREMENT=345 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`apiKeyBrandSchema`:      `CREATE TABLE ApiKeyToBrand (ID int(11) NOT NULL AUTO_INCREMENT,keyID int(11) NOT NULL,brandID int(11) NOT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=38361 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT`,
		`apiKeySchema`:           `CREATE TABLE ApiKey (id int(11) NOT NULL AUTO_INCREMENT,api_key varchar(64) NOT NULL,type_id varchar(64) NOT NULL,user_id varchar(64) NOT NULL,date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,UNIQUE KEY id (id)) ENGINE=InnoDB AUTO_INCREMENT=14489 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`applicationGuideSchema`: `CREATE TABLE ApplicationGuides (ID int(11) unsigned NOT NULL AUTO_INCREMENT,url varchar(255) NOT NULL DEFAULT '',websiteID int(11) NOT NULL,fileType varchar(15) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,catID int(11) NOT NULL,icon varchar(255) NOT NULL DEFAULT '',brandID int(11) NOT NULL DEFAULT '1',PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
	}

	dataInserts = map[string]string{
		`insertKeyType`:          `INSERT INTO ApiKeyType (id, type, date_added) VALUES ('a46ceab9-df0c-44a2-b21d-a859fc2c839c','random type', NOW())`,
		`insertApiKey`:           `INSERT INTO ApiKey(id, api_key, type_id, user_id) values(1, '` + Key + `', 'a46ceab9-df0c-44a2-b21d-a859fc2c839c', UUID())`,
		`insertKeyToBrand`:       `INSERT INTO ApiKeyToBrand(ID, keyID, brandID) VALUES(1, 1, 3)`,
		`insertCategory`:         `INSERT INTO Categories (catID, dateAdded, parentID, catTitle, shortDesc, longDesc, image, isLifestyle, codeID, sort, vehicleSpecific, vehicleRequired, metaTitle, metaDesc, metaKeywords, icon, path, brandID, isDeleted) VALUES(1, NOW(), 0, 'Random Category', 'Random Category Short Desc', 'Random Category Long Desc', 'https://www.curtmfg.com/masterlibrary/01CategoryImages/DH/Receiver_Hitches.png', 0, 1, 0, 1, 1, 'Random Category meta title', 'Random Category meta description', '', 'https://storage.googleapis.com/curt-icons/category/Trailer_Hitches.png', 'Random Category', 3, 0)`,
		`insertApplicationGuide`: `INSERT INTO ApplicationGuides (ID, url, websiteID, fileType, catID, icon, brandID) VALUES(1, 'https://storage.googleapis.com/aries-applicationguides/ARIES%20Complete%20App%20Guide.csv', 3, 'csv', 1, 'storage.googleapis.com/aries-website/csv_file.png', 3)`,
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

func TestGetApplicationGuide(t *testing.T) {

	Convey("Testing GetApplicationGuide", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &customer.DataContext{
				BrandID: 3,
				APIKey:  Key,
			},
			Params: httprouter.Params{},
			DB:     db,
		}

		Convey("with a non-numeric identifier", func() {
			ctx.Params = httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "a",
				},
			}
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/applicationGuide/0", nil)
			So(err, ShouldBeNil)

			resp, err := GetApplicationGuide(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with a zero identifier", func() {
			ctx.Params = httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "0",
				},
			}
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/applicationGuide/0", nil)
			So(err, ShouldBeNil)

			resp, err := GetApplicationGuide(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with a valid identifier", func() {
			ctx.Params = httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "1",
				},
			}
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/applicationGuide/1", nil)
			So(err, ShouldBeNil)

			resp, err := GetApplicationGuide(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, applicationGuide.ApplicationGuide{})
		})

	})
}

func TestGetApplicationGuidesByWebsite(t *testing.T) {

	Convey("Testing GetApplicationGuidesByWebsite", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &customer.DataContext{
				BrandID: 3,
				APIKey:  Key,
			},
			Params: httprouter.Params{},
			DB:     db,
		}

		Convey("with a invalid identifier", func() {
			ctx.Params = httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "a",
				},
			}
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/applicationGuide/a/website", nil)
			So(err, ShouldBeNil)

			resp, err := GetApplicationGuidesByWebsite(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with a valid identifier", func() {
			ctx.Params = httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "3",
				},
			}
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/applicationGuide/1/website", nil)
			So(err, ShouldBeNil)

			resp, err := GetApplicationGuidesByWebsite(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []applicationGuide.ApplicationGuide{})
		})

	})
}
