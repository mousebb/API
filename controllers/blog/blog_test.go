package blog_controller

import (
	"github.com/curt-labs/API/helpers/apicontextmock"
	"github.com/curt-labs/API/helpers/pagination"
	"github.com/curt-labs/API/helpers/testThatHttp"
	"github.com/curt-labs/API/models/blog"
	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestBlog(t *testing.T) {
	var b blog_model.Blog
	var bc blog_model.BlogCategory
	var err error
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}
	Convey("Testing Blog", t, func() {
		//test create blog cats
		form := url.Values{"name": {"test cat"}, "slug": {"a slug here"}}
		v := form.Encode()
		body := strings.NewReader(v)
		testThatHttp.Request("post", "/blogs/categories", "", "", CreateBlogCategory, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bc)
		So(err, ShouldBeNil)
		So(bc, ShouldHaveSameTypeAs, blog_model.BlogCategory{})

		//test create blog
		form = url.Values{"title": {"test"}, "slug": {"a slug"}, "texts": {"some text here"}, "categoryID": {strconv.Itoa(bc.ID)}}
		v = form.Encode()
		body = strings.NewReader(v)
		testThatHttp.RequestWithDtx("post", "/blogs", "", "", CreateBlog, body, "application/x-www-form-urlencoded", dtx)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &b)
		So(err, ShouldBeNil)
		So(b, ShouldHaveSameTypeAs, blog_model.Blog{})

		//test get blogs
		testThatHttp.Request("get", "/blog", "", "", GetAll, nil, "application/x-www-form-urlencoded")
		var bs blog_model.Blogs
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bs)
		So(len(bs), ShouldBeGreaterThanOrEqualTo, 0)
		So(err, ShouldBeNil)

		//test get blog
		testThatHttp.Request("get", "/blog/", ":id", strconv.Itoa(b.ID), GetBlog, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &b)
		So(b, ShouldHaveSameTypeAs, blog_model.Blog{})
		So(err, ShouldBeNil)
		So(b.Title, ShouldEqual, "test")

		//test get blog cats
		testThatHttp.RequestWithDtx("get", "/blog/categories", "", "", GetAllCategories, nil, "application/x-www-form-urlencoded", dtx)
		var bcs blog_model.BlogCategories
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bcs)
		So(len(bcs), ShouldBeGreaterThanOrEqualTo, 0)
		So(err, ShouldBeNil)

		//test get blog cat
		testThatHttp.RequestWithDtx("get", "/blog/category/", ":id", strconv.Itoa(bc.ID), GetBlogCategory, nil, "application/x-www-form-urlencoded", dtx)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bc)
		So(bc, ShouldHaveSameTypeAs, blog_model.BlogCategory{})
		So(err, ShouldBeNil)

		//test search
		testThatHttp.RequestWithDtx("get", "/blog/search/", "", "?title="+b.Title, Search, nil, "application/x-www-form-urlencoded", dtx)
		var l pagination.Objects
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &l)
		So(len(l.Objects), ShouldBeGreaterThanOrEqualTo, 0)
		So(err, ShouldBeNil)

		//test update blog
		form = url.Values{"name": {"test cat"}, "slug": {"a slug here"}}
		v = form.Encode()
		body = strings.NewReader(v)
		testThatHttp.Request("put", "/blogs/", ":id", strconv.Itoa(b.ID), UpdateBlog, body, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &b)
		So(err, ShouldBeNil)
		So(b, ShouldHaveSameTypeAs, blog_model.Blog{})

		//test delete blog cat
		testThatHttp.Request("delete", "/blog/categories/", ":id", strconv.Itoa(bc.ID), DeleteBlogCategory, nil, "application/x-www-form-urlencoded")
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &bc)
		So(err, ShouldBeNil)

		//test delete blog
		testThatHttp.RequestWithDtx("delete", "/blog/", ":id", strconv.Itoa(b.ID), DeleteBlog, nil, "application/x-www-form-urlencoded", dtx)
		So(testThatHttp.Response.Code, ShouldEqual, 200)
		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &b)
		So(err, ShouldBeNil)

	})
	_ = apicontextmock.DeMock(dtx)
}

func BenchmarkBlog(b *testing.B) {
	testThatHttp.RequestBenchmark(b.N, "GET", "/blog/1", nil, GetBlog)
	testThatHttp.RequestBenchmark(b.N, "GET", "/blog", nil, GetAll)
	testThatHttp.RequestBenchmark(b.N, "GET", "/blog/categories", nil, GetAllCategories)
	testThatHttp.RequestBenchmark(b.N, "GET", "/blog/category/1", nil, GetBlogCategory)
}
