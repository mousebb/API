package applicationGuide

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/customer"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	db     *sql.DB
	apiKey string

	drops = map[string]string{
		`dropApplicationGuides`: `DROP TABLE IF EXISTS ApplicationGuides`,
		`dropCategories`:        `DROP TABLE IF EXISTS Categories`,
		`dropApiKey`:            `DROP TABLE IF EXISTS ApiKey`,
		`dropApiKeyToBrand`:     `DROP TABLE IF EXISTS ApiKeyToBrand`,
	}

	schemas = map[string]string{
		`applicationGuideSchema`: `CREATE TABLE ApplicationGuides (
			  ID int(11) unsigned NOT NULL AUTO_INCREMENT,
			  url varchar(255) NOT NULL DEFAULT '',
			  websiteID int(11) NOT NULL,
			  fileType varchar(15) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
			  catID int(11) NOT NULL,
			  icon varchar(255) NOT NULL DEFAULT '',
			  brandID int(11) NOT NULL DEFAULT '1',
			  PRIMARY KEY (ID)
			) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`,
		`categories`: `CREATE TABLE Categories (
			  catID int(11) NOT NULL AUTO_INCREMENT,
			  dateAdded timestamp NULL DEFAULT CURRENT_TIMESTAMP,
			  parentID int(11) NOT NULL,
			  catTitle varchar(100) DEFAULT NULL,
			  shortDesc varchar(255) DEFAULT NULL,
			  longDesc longtext,
			  image varchar(255) DEFAULT NULL,
			  isLifestyle int(11) NOT NULL,
			  codeID int(11) NOT NULL DEFAULT '0',
			  sort int(11) NOT NULL DEFAULT '1',
			  vehicleSpecific tinyint(1) NOT NULL DEFAULT '0',
			  vehicleRequired tinyint(1) NOT NULL DEFAULT '0',
			  metaTitle text,
			  metaDesc text,
			  metaKeywords text,
			  icon varchar(255) DEFAULT NULL,
			  path varchar(255) DEFAULT NULL,
			  brandID int(11) NOT NULL DEFAULT '1',
			  isDeleted tinyint(1) NOT NULL DEFAULT '0',
			  PRIMARY KEY (catID),
			  KEY IX_Categories_ParentID (parentID),
			  KEY IX_Categories_Sort (sort),
			  KEY brandID (brandID)
			) ENGINE=InnoDB AUTO_INCREMENT=345 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`,
		`apiKey`: `CREATE TABLE ApiKey (
			  id int(11) NOT NULL AUTO_INCREMENT,
			  api_key varchar(64) NOT NULL,
			  type_id varchar(64) NOT NULL,
			  user_id varchar(64) NOT NULL,
			  date_added timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  UNIQUE KEY id (id),
			  KEY FK__ApiKey__type_id__5AEE1AF6 (type_id),
			  KEY FK__ApiKey__user_id__5BE23F2F (user_id),
			  KEY api_key (api_key)
			) ENGINE=InnoDB AUTO_INCREMENT=14433 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT;`,
		`apiKeyToBrand`: `CREATE TABLE ApiKeyToBrand (
			  ID int(11) NOT NULL AUTO_INCREMENT,
			  keyID int(11) NOT NULL,
			  brandID int(11) NOT NULL,
			  PRIMARY KEY (ID),
			  KEY FK_ApiKeyToBrand_ApiKey (keyID),
			  KEY FK_ApiKeyToBrand_Brand (brandID)
			) ENGINE=InnoDB AUTO_INCREMENT=38350 DEFAULT CHARSET=latin1 ROW_FORMAT=COMPACT;`,
	}

	dataInserts = map[string]string{
		`insertAppGuide`: `INSERT INTO ApplicationGuides (ID, url, websiteID, fileType, catID, icon, brandID)
			VALUES (1, 'http://imgur.com/gallery/anQ7zvr', 1, 'jpg', 1, 'www.curtmfg.com/assets/434da33a-2abd-4821-a236-562d38be3e79.png', 3)`,
		`insertApiKey`: `insert into ApiKey (id, api_key, type_id, user_id, date_added)	values(1, UUID(), UUID(), UUID(), NOW())`,
		`insertApiKeyToBrand`: `insert into ApiKeyToBrand (keyID, brandID) values(1, 3)`,
	}
)

func TestMain(m *testing.M) {
	var err error
	if os.Getenv("CI") == "" {
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

			err = insertApplicationGuides()
			if err != nil {
				log.Fatal(err)
			}

			err = getAPIKey()
			if err != nil {
				log.Fatal(err)
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

		err = insertApplicationGuides()
		if err != nil {
			log.Fatal(err)
		}

		err = getAPIKey()
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
	}

	m.Run()
}

func insertApplicationGuides() error {
	for _, insert := range dataInserts {
		_, err := db.Exec(insert)
		if err != nil {
			return err
		}
	}
	return nil
}

func getAPIKey() error {
	return db.QueryRow("select api_key from ApiKey where id = 1").Scan(&apiKey)
}

func TestGetBySite(t *testing.T) {
	Convey("GetBySite", t, func() {

		ctx := &middleware.APIContext{
			DataContext: &customer.DataContext{
				BrandID: 3,
				APIKey:  apiKey,
			},
			Params: httprouter.Params{},
			DB:     db,
		}
		Convey("Bad query", func() {
			ag := ApplicationGuide{
				Website: Website{
					ID: 1,
				},
			}

			tmp := getApplicationGuidesBySite
			getApplicationGuidesBySite = "bogus query"
			ags, err := ag.GetBySite(ctx)
			So(err, ShouldNotBeNil)
			So(len(ags), ShouldBeZeroValue)
			getApplicationGuidesBySite = tmp
		})

		Convey("with missing select columns", func() {
			ag := ApplicationGuide{
				Website: Website{
					ID: 1,
				},
			}

			tmp := getApplicationGuidesBySite
			getApplicationGuidesBySite = "select ag.ID from ApplicationGuides as ag"
			ags, err := ag.GetBySite(ctx)
			So(err, ShouldNotBeNil)
			So(len(ags), ShouldBeZeroValue)
			getApplicationGuidesBySite = tmp
		})

		Convey("invalid website", func() {
			ag := ApplicationGuide{
				Website: Website{
					ID: 10,
				},
			}

			ags, err := ag.GetBySite(ctx)
			So(err, ShouldBeNil)
			So(len(ags), ShouldEqual, 0)
		})

		Convey("valid", func() {
			ag := ApplicationGuide{
				Website: Website{
					ID: 1,
				},
			}

			ags, err := ag.GetBySite(ctx)
			So(err, ShouldBeNil)
			So(len(ags), ShouldBeGreaterThan, 0)
		})
	})
}

func TestGet(t *testing.T) {
	Convey("Get", t, func() {

		ctx := &middleware.APIContext{
			DataContext: &customer.DataContext{
				BrandID: 3,
				APIKey:  apiKey,
			},
			Params: httprouter.Params{},
			DB:     db,
		}
		Convey("Bad query", func() {
			ag := ApplicationGuide{
				Website: Website{
					ID: 1,
				},
			}

			tmp := getApplicationGuide
			getApplicationGuide = "bogus query"
			err := ag.Get(ctx)
			So(err, ShouldNotBeNil)
			getApplicationGuide = tmp
		})

		Convey("with missing select columns", func() {
			ag := ApplicationGuide{
				Website: Website{
					ID: 1,
				},
			}

			tmp := getApplicationGuide
			getApplicationGuide = "select ag.ID from ApplicationGuides as ag"
			err := ag.Get(ctx)
			So(err, ShouldNotBeNil)
			getApplicationGuide = tmp
		})

		Convey("valid", func() {
			ag := ApplicationGuide{
				ID: 1,
			}
			err := ag.Get(ctx)
			So(err, ShouldBeNil)
		})
	})
}
