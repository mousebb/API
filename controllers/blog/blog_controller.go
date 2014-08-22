package blog_controller

import (
	"errors"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/blog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func GetAll(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	blogs, err := blog_model.GetAll()
	if err != nil {
		return err.Error()
	}
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(blogs, sort)
		} else {
			sortutil.AscByField(blogs, sort)
		}

	}
	return encoding.Must(enc.Encode(blogs))
}

func GetAllCategories(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	cats, err := blog_model.GetAllCategories()
	if err != nil {
		return err.Error()
	}
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(cats, sort)
		} else {
			sortutil.AscByField(cats, sort)
		}

	}
	return encoding.Must(enc.Encode(cats))
}

func GetBlog(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var b blog_model.Blog
	var err error
	b.ID, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return err.Error()
	}
	err = b.Get()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(b))
}

func CreateBlog(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var b blog_model.Blog
	var err error

	b.Title = r.FormValue("title")
	b.Slug = r.FormValue("slug")
	b.Text = r.FormValue("text")
	b.PublishedDate, err = time.Parse(timeFormat, r.FormValue("publishedDate"))
	b.UserID, err = strconv.Atoi(r.FormValue("userID"))
	b.MetaTitle = r.FormValue("metaTitle")
	b.MetaDescription = r.FormValue("metaDescription")
	b.Keywords = r.FormValue("keywords")
	b.Active, err = strconv.ParseBool(r.FormValue("active"))
	categoryIDs := r.Form["categoryID"]
	for _, v := range categoryIDs {
		var bc blog_model.BlogCategory
		bc.Category.ID, err = strconv.Atoi(v)
		b.BlogCategories = append(b.BlogCategories, bc)
	}

	err = b.Create()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(b))
}
func GetBlogCategory(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var c blog_model.Category
	var err error
	c.ID, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return err.Error()
	}
	err = c.Get()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(c))
}
func CreateBlogCategory(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var c blog_model.Category
	var err error

	c.Name = r.FormValue("name")
	c.Slug = r.FormValue("slug")
	c.Active, err = strconv.ParseBool(r.FormValue("active"))

	err = c.Create()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(c))
}

func UpdateBlog(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var b blog_model.Blog
	var err error

	id := r.FormValue("id")
	if id == "" || err != nil {
		return fmt.Sprint(errors.New("Invalid ID supplied."), err)
	}
	b.ID, err = strconv.Atoi(id)
	b.Get()

	var tempBC []blog_model.BlogCategory

	title := r.FormValue("title")
	slug := r.FormValue("slug")
	text := r.FormValue("text")
	publishedDate := r.FormValue("publishedDate")
	userID := r.FormValue("userID")
	metaTitle := r.FormValue("metaTitle")
	metaDescription := r.FormValue("metaDescription")
	keywords := r.FormValue("keywords")
	active := r.FormValue("active")
	categoryIDs := r.Form["categoryID"]
	for _, v := range categoryIDs {
		var bc blog_model.BlogCategory
		bc.Category.ID, err = strconv.Atoi(v)
		tempBC = append(tempBC, bc)
	}

	if err != nil {
		return err.Error()
	}
	if title != "" {
		b.Title = title
	}
	if slug != "" {
		b.Slug = slug
	}
	if text != "" {
		b.Text = text
	}
	if publishedDate != "" {
		b.PublishedDate, err = time.Parse(timeFormat, publishedDate)
	}
	if userID != "" {
		b.UserID, err = strconv.Atoi(userID)
	}
	if metaTitle != "" {
		b.MetaTitle = metaTitle
	}
	if metaDescription != "" {
		b.MetaDescription = metaDescription
	}
	if keywords != "" {
		b.Keywords = keywords
	}
	if active != "" {
		b.Active, err = strconv.ParseBool(active)
	}
	if categoryIDs != nil {
		b.BlogCategories = tempBC
	}

	err = b.Update()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(b))
}

func DeleteBlog(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var b blog_model.Blog
	var err error
	b.ID, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return err.Error()
	}
	err = b.Delete()
	if err != nil {
		return err.Error()
	}
	return encoding.Must(enc.Encode(b))
}

func Search(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	var err error

	title := r.FormValue("title")
	slug := r.FormValue("slug")
	text := r.FormValue("text")
	createdDate := r.FormValue("createdDate")
	publishedDate := r.FormValue("publishedDate")
	lastModified := r.FormValue("lastModified")
	userID := r.FormValue("userID")
	metaTitle := r.FormValue("metaTitle")
	metaDescription := r.FormValue("metaDescription")
	keywords := r.FormValue("keywords")
	active := r.FormValue("active")

	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err := blog_model.Search(title, slug, text, publishedDate, createdDate, lastModified, userID, metaTitle, metaDescription, keywords, active, page, results)
	if err != nil {
		return err.Error()
	}

	return encoding.Must(enc.Encode(l))
}