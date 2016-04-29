package customerCtlr

// import (
// 	"encoding/json"
// 	"net/url"
// 	"strconv"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/curt-labs/httprunner"
// 	"github.com/curt-labs/API/helpers/testThatHttp"
// 	"github.com/curt-labs/API/models/customer"
// 	. "github.com/smartystreets/goconvey/convey"
// )

// func TestCustomerPrice(t *testing.T) {
// 	var err error
// 	var p customer.Price
// 	var ps customer.Prices
// 	var c customer.Customer
// 	c.Name = "Dog Bountyhunter"
// 	c.Create()

// 	Convey("Testing customer/Price", t, func() {
// 		//test create customer price
// 		form := url.Values{"custID": {strconv.Itoa(c.Id)}, "partID": {"11000"}, "price": {"123456"}}
// 		v := form.Encode()
// 		body := strings.NewReader(v)
// 		thyme := time.Now()
// 		testThatHttp.Request("post", "/customer/prices", "", "", CreateUpdatePrice, body, "application/x-www-form-urlencoded")
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
// 		So(err, ShouldBeNil)
// 		So(p, ShouldHaveSameTypeAs, customer.Price{})
// 		So(p.ID, ShouldBeGreaterThan, 0)

// 		//test update customer price
// 		form = url.Values{"isSale": {"true"}, "saleStart": {"01/01/2001"}, "saleEnd": {"01/01/2015"}}
// 		v = form.Encode()
// 		body = strings.NewReader(v)
// 		thyme = time.Now()
// 		testThatHttp.Request("post", "/customer/prices/", ":id", strconv.Itoa(p.ID), CreateUpdatePrice, body, "application/x-www-form-urlencoded")
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
// 		So(err, ShouldBeNil)
// 		So(p, ShouldHaveSameTypeAs, customer.Price{})
// 		So(p.IsSale, ShouldEqual, 1)
// 		start, _ := time.Parse(inputTimeFormat, "01/01/2001")
// 		So(p.SaleStart, ShouldResemble, start)

// 		//test get customer price
// 		thyme = time.Now()
// 		testThatHttp.Request("get", "/customer/prices/", ":id", strconv.Itoa(p.ID), GetPrice, nil, "")
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
// 		So(err, ShouldBeNil)
// 		So(p, ShouldHaveSameTypeAs, customer.Price{})

// 		//test get all customer price
// 		thyme = time.Now()
// 		testThatHttp.Request("get", "/customer/prices", "", "", GetAllPrices, nil, "")
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()*8) //Long
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ps)
// 		So(err, ShouldBeNil)
// 		So(ps, ShouldHaveSameTypeAs, customer.Prices{})

// 		//test get customer price by part
// 		thyme = time.Now()
// 		testThatHttp.Request("get", "/customer/prices/part/", ":id", strconv.Itoa(p.ID), GetPricesByPart, nil, "")
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ps)
// 		So(err, ShouldBeNil)
// 		So(ps, ShouldHaveSameTypeAs, customer.Prices{})

// 		//test get customer price by customer
// 		thyme = time.Now()
// 		testThatHttp.Request("get", "/customer/pricesByCustomer/", ":id", strconv.Itoa(c.Id), GetPriceByCustomer, nil, "")
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
// 		So(err, ShouldBeNil)
// 		So(p, ShouldHaveSameTypeAs, customer.Price{})

// 		//test get sales
// 		form = url.Values{"id": {strconv.Itoa(c.Id)}, "start": {"01/01/2000"}, "end": {"01/01/2016"}}
// 		v = form.Encode()
// 		body = strings.NewReader(v)
// 		thyme = time.Now()
// 		testThatHttp.Request("post", "/customer/prices/sale", "", "", GetSales, body, "application/x-www-form-urlencoded")
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &ps)
// 		So(err, ShouldBeNil)
// 		So(ps, ShouldHaveSameTypeAs, customer.Prices{})

// 		//test delete customer price
// 		thyme = time.Now()
// 		testThatHttp.Request("delete", "/customer/prices/", ":id", strconv.Itoa(p.ID), DeletePrice, nil, "")
// 		So(time.Since(thyme).Nanoseconds(), ShouldBeLessThan, time.Second.Nanoseconds()/2)
// 		So(testThatHttp.Response.Code, ShouldEqual, 200)
// 		err = json.Unmarshal(testThatHttp.Response.Body.Bytes(), &p)
// 		So(err, ShouldBeNil)
// 		So(p, ShouldHaveSameTypeAs, customer.Price{})
// 	})
// 	//teardown
// 	c.Delete()
// }

// func BenchmarkCRUDCustomerPrice(b *testing.B) {
// 	var p customer.Price
// 	var c customer.Customer
// 	c.Name = "Axl Rose"
// 	c.Create()
// 	qs := make(url.Values, 0)

// 	Convey("CustomerPrice", b, func() {
// 		form := url.Values{"custID": {strconv.Itoa(c.Id)}, "partID": {"11000"}, "price": {"123456"}}
// 		//create
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "POST",
// 			Route:              "/customer/prices",
// 			ParameterizedRoute: "/customer/prices",
// 			Handler:            CreateUpdatePrice,
// 			QueryString:        &qs,
// 			JsonBody:           p,
// 			FormBody:           form,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//get
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "GET",
// 			Route:              "/customer/prices",
// 			ParameterizedRoute: "/customer/prices/" + strconv.Itoa(p.ID),
// 			Handler:            GetPrice,
// 			QueryString:        &qs,
// 			JsonBody:           p,
// 			FormBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//get all
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "GET",
// 			Route:              "/customer/prices",
// 			ParameterizedRoute: "/customer/prices",
// 			Handler:            GetAllPrices,
// 			QueryString:        &qs,
// 			JsonBody:           p,
// 			FormBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//get by part
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "GET",
// 			Route:              "/customer/prices/part",
// 			ParameterizedRoute: "/customer/prices/part/" + strconv.Itoa(p.ID),
// 			Handler:            GetPricesByPart,
// 			QueryString:        &qs,
// 			JsonBody:           p,
// 			FormBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//get by
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "GET",
// 			Route:              "/customer/pricesByCustomer",
// 			ParameterizedRoute: "/customer/pricesByCustomer/" + strconv.Itoa(c.Id),
// 			Handler:            GetPriceByCustomer,
// 			QueryString:        &qs,
// 			JsonBody:           p,
// 			FormBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()

// 		//delete
// 		(&httprunner.BenchmarkOptions{
// 			Method:             "DELETE",
// 			Route:              "/customer/prices",
// 			ParameterizedRoute: "/customer/prices/" + strconv.Itoa(p.ID),
// 			Handler:            DeleteLocation,
// 			QueryString:        &qs,
// 			JsonBody:           p,
// 			FormBody:           nil,
// 			Runs:               b.N,
// 		}).RequestBenchmark()
// 	})
// 	//teardown
// 	c.Delete()
// }
