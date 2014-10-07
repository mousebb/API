package video

import (
	"database/sql"
	"net/url"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

type Video struct {
	YouTubeId   string
	DateAdded   time.Time
	Sort        int
	Title       string
	Description string
	Watchpage   *url.URL
	Screenshot  *url.URL
}

var (
	uniqueVideoStmt = `select distinct embed_link, dateAdded, sort, title, description, watchpage, screenshot
				from Video
				order by sort`
)

func UniqueVideos() (videos []Video, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(uniqueVideoStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var v Video
		err = rows.Scan(
			&v.YouTubeId,
			&v.DateAdded,
			&v.Sort,
			&v.Title,
			&v.Description,
			&v.Watchpage,
			&v.Screenshot,
		)
		if err != nil {
			return
		}
		videos = append(videos, v)
	}

	return
}
