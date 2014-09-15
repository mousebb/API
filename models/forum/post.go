package forum

import (
	"database/sql"
	"errors"
	"time"

	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllForumPosts    = `select * from ForumPost`
	getForumPost        = `select * from ForumPost where postID = ?`
	getForumThreadPosts = `select * from ForumPost where threadID = ?`
	addForumPost        = `insert into ForumPost(parentID,threadID,createdDate,title,post,name,email,company,notify,approved,active,IPAddress,flag,sticky) values(?,?,UTC_TIMESTAMP(),?,?,?,?,?,?,?,1,?,?,?)`
	updateForumPost     = `update ForumPost set parentID = ?, threadID = ?, title = ?, post = ?, name = ?, email = ?, company = ?, notify = ?, approved = ?, IPAddress = ?, flag = ?, sticky = ? where postID = ?`
	deleteForumPost     = `delete from ForumPost where postID = ?`
)

type Posts []Post
type Post struct {
	ID        int
	ParentID  int
	ThreadID  int
	Created   time.Time
	Title     string
	Post      string
	Name      string
	Email     string
	Company   string
	Notify    bool
	Approved  bool
	Active    bool
	IPAddress string
	Flag      bool
	Sticky    bool
}

func GetAllPosts() (posts Posts, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllForumPosts)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var post Post
		if err = rows.Scan(&post.ID, &post.ParentID, &post.ThreadID, &post.Created, &post.Title, &post.Post, &post.Name, &post.Email, &post.Company, &post.Notify, &post.Approved, &post.Active, &post.IPAddress, &post.Flag, &post.Sticky); err == nil {
			posts = append(posts, post)
		}
	}

	return
}

func (p *Post) Get() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumPost)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var post Post
	row := stmt.QueryRow(p.ID)
	err = row.Scan(&post.ID, &post.ParentID, &post.ThreadID, &post.Created, &post.Title, &post.Post, &post.Name, &post.Email, &post.Company, &post.Notify, &post.Approved, &post.Active, &post.IPAddress, &post.Flag, &post.Sticky)

	if row == nil || err != nil {
		if row == nil {
			return errors.New("Invalid reference to Forum Post")
		}
		return err
	}

	p.ID = post.ID
	p.ParentID = post.ParentID
	p.ThreadID = post.ThreadID
	p.Created = post.Created
	p.Title = post.Title
	p.Post = post.Post
	p.Name = post.Name
	p.Email = post.Email
	p.Company = post.Company
	p.Notify = post.Notify
	p.Approved = post.Approved
	p.Active = post.Active
	p.IPAddress = post.IPAddress
	p.Flag = post.Flag
	p.Sticky = post.Sticky

	return nil
}

func (t *Thread) GetPosts() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getForumThreadPosts)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(t.ID)
	for rows.Next() {
		var post Post
		if err = rows.Scan(&post.ID, &post.ParentID, &post.ThreadID, &post.Created, &post.Title, &post.Post, &post.Name, &post.Email, &post.Company, &post.Notify, &post.Approved, &post.Active, &post.IPAddress, &post.Flag, &post.Sticky); err == nil {
			t.Posts = append(t.Posts, post)
		}
	}

	return nil
}

func (p *Post) Add() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(addForumPost)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.ParentID, p.ThreadID, p.Title, p.Post, p.Name, p.Email, p.Company, p.Notify, p.Approved, p.IPAddress, p.Flag, p.Sticky)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		p.ID = int(id)
	}

	return nil
}

func (p *Post) Update() error {
	if p.ID == 0 {
		return errors.New("Invalid Post ID")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updateForumPost)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(p.ParentID, p.ThreadID, p.Title, p.Post, p.Name, p.Email, p.Company, p.Notify, p.Approved, p.IPAddress, p.Flag, p.Sticky, p.ID); err != nil {
		return err
	}

	return nil
}

func (p *Post) Delete() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deleteForumPost)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(p.ID); err != nil {
		return err
	}

	return nil
}