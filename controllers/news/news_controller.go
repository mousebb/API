package news_controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/helpers/pagination"
	"github.com/curt-labs/API/helpers/sortutil"
	"github.com/curt-labs/API/models/news"
	"github.com/go-martini/martini"
)

const (
	timeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

func GetAll(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var fs news_model.Newses
	var err error

	fs, err = news_model.GetAll(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting all news", err, rw, r)
		return ""
	}
	sort := r.FormValue("sort")
	direction := r.FormValue("direction")
	if sort != "" {
		if strings.ContainsAny(direction, "esc") {
			sortutil.DescByField(fs, sort)
		} else {
			sortutil.AscByField(fs, sort)
		}

	}
	return encoding.Must(enc.Encode(fs))
}

func Get(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var f news_model.News
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		f.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting news ID", err, rw, r)
			return ""
		}
	} else {
		f.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting news ID", err, rw, r)
			return ""
		}
	}
	err = f.Get(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting news", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(f))
}

func Create(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var n news_model.News
	var err error

	n.Title = r.FormValue("title")
	n.Lead = r.FormValue("lead")
	n.Content = r.FormValue("content")
	start := r.FormValue("start")
	end := r.FormValue("end")
	active := r.FormValue("active")
	n.Slug = r.FormValue("slug")
	if start != "" {
		n.PublishStart, err = time.Parse(timeFormat, start)
	}
	if end != "" {
		n.PublishEnd, err = time.Parse(timeFormat, end)
	}
	if active != "" {
		n.Active, err = strconv.ParseBool(active)
	}
	err = n.Create(dtx)
	if err != nil {
		apierror.GenerateError("Trouble creating news", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(n))
}

func Update(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var n news_model.News
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		n.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting news ID", err, rw, r)
			return ""
		}
	} else {
		n.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting news ID", err, rw, r)
			return ""
		}
	}
	n.Get(dtx)
	title := r.FormValue("title")
	lead := r.FormValue("lead")
	content := r.FormValue("content")
	start := r.FormValue("start")
	end := r.FormValue("end")
	active := r.FormValue("active")
	slug := r.FormValue("slug")

	if title != "" {
		n.Title = title
	}
	if lead != "" {
		n.Lead = lead
	}
	if content != "" {
		n.Content = content
	}
	if start != "" {
		n.PublishStart, err = time.Parse(timeFormat, start)
	}
	if end != "" {
		n.PublishEnd, err = time.Parse(timeFormat, end)
	}
	if active != "" {
		n.Active, err = strconv.ParseBool(active)
	}
	if slug != "" {
		n.Slug = slug
	}

	err = n.Update(dtx)
	if err != nil {
		apierror.GenerateError("Trouble updating news", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(n))
}

func Delete(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, params martini.Params, dtx *apicontext.DataContext) string {
	var n news_model.News
	var err error

	idStr := r.FormValue("id")
	if idStr != "" {
		n.ID, err = strconv.Atoi(idStr)
		if err != nil {
			apierror.GenerateError("Trouble getting news ID", err, rw, r)
			return ""
		}
	} else {
		n.ID, err = strconv.Atoi(params["id"])
		if err != nil {
			apierror.GenerateError("Trouble getting news ID", err, rw, r)
			return ""
		}
	}
	err = n.Delete(dtx)
	if err != nil {
		apierror.GenerateError("Trouble deleting news", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(n))
}

func GetTitles(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var l pagination.Objects
	var err error
	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err = news_model.GetTitles(page, results, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting news titles", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(l))
}

func GetLeads(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var l pagination.Objects
	var err error
	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err = news_model.GetLeads(page, results, dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting news leads", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(l))
}

func Search(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var err error

	title := r.FormValue("title")
	lead := r.FormValue("lead")
	content := r.FormValue("content")
	publishStart := r.FormValue("publishStart")
	publishEnd := r.FormValue("publishEnd")
	active := r.FormValue("active")
	slug := r.FormValue("slug")

	page := r.FormValue("page")
	results := r.FormValue("results")

	l, err := news_model.Search(title, lead, content, publishStart, publishEnd, active, slug, page, results, dtx)
	if err != nil {
		apierror.GenerateError("Trouble searching news", err, rw, r)
		return ""
	}
	return encoding.Must(enc.Encode(l))
}
