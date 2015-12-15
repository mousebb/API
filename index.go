package main

import (
	"flag"
	"net/http"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/router"

	"log"
	"time"
)

var (
	listenAddr = flag.String("http", ":8080", "http listen address")
)

/**
 * All GET routes require either public or private api keys to be passed in.
 *
 * All POST routes require private api keys to be passed in.
 */
func main() {
	flag.Parse()

	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	r := router.New()

	// m.Group("/customer", func(r martini.Router) {
	// 	r.Get("", customer_ctlr.GetCustomer)
	// 	r.Post("", customer_ctlr.GetCustomer)
	//
	// 	r.Post("/auth", customer_ctlr.AuthenticateUser)
	// 	r.Get("/auth", customer_ctlr.KeyedUserAuthentication)
	// 	r.Post("/user/changePassword", customer_ctlr.ChangePassword)
	// 	r.Post("/user", customer_ctlr.GetUser)
	// 	r.Post("/user/register", customer_ctlr.RegisterUser)
	// 	r.Post("/user/resetPassword", customer_ctlr.ResetPassword)
	// 	r.Delete("/deleteKey", customer_ctlr.DeleteUserApiKey)
	// 	r.Post("/generateKey/user/:id/key/:type", customer_ctlr.GenerateApiKey)
	// 	r.Get("/user/:id", customer_ctlr.GetUserById)
	// 	r.Post("/user/:id", customer_ctlr.UpdateCustomerUser)
	// 	r.Delete("/user/:id", customer_ctlr.DeleteCustomerUser)
	// 	r.Any("/users", customer_ctlr.GetUsers)
	//
	// 	r.Delete("/allUsersByCustomerID/:id", customer_ctlr.DeleteCustomerUsersByCustomerID) //Takes CustomerID (UUID)---danger!
	//
	// 	r.Put("/location/json", customer_ctlr.SaveLocationJson)
	// 	r.Put("/location/json/:id", customer_ctlr.SaveLocationJson)
	// 	r.Post("/location", customer_ctlr.SaveLocation)
	// 	r.Get("/location/:id", customer_ctlr.GetLocation)
	// 	r.Put("/location/:id", customer_ctlr.SaveLocation)
	// 	r.Delete("/location/:id", customer_ctlr.DeleteLocation)
	//
	// 	r.Get("/locations", customer_ctlr.GetLocations)
	// 	r.Post("/locations", customer_ctlr.GetLocations)
	//
	// 	r.Get("/price/:id", customer_ctlr.GetCustomerPrice)           //{part id}
	// 	r.Get("/cartRef/:id", customer_ctlr.GetCustomerCartReference) //{part id}
	//
	// 	// Customer CMS endpoints
	// 	// All Customer Contents
	// 	r.Get("/cms", customer_ctlr.GetAllContent)
	// 	// Content Types
	// 	r.Get("/cms/content_types", customer_ctlr.GetAllContentTypes)
	//
	// 	// Customer Part Content
	// 	r.Get("/cms/part", customer_ctlr.AllPartContent)
	// 	r.Get("/cms/part/:id", customer_ctlr.UniquePartContent)
	// 	r.Put("/cms/part/:id", customer_ctlr.UpdatePartContent) //partId
	// 	r.Post("/cms/part/:id", customer_ctlr.CreatePartContent)
	// 	r.Delete("/cms/part/:id", customer_ctlr.DeletePartContent)
	//
	// 	// Customer Category Content
	// 	r.Get("/cms/category", customer_ctlr.AllCategoryContent)
	// 	r.Get("/cms/category/:id", customer_ctlr.UniqueCategoryContent)
	// 	r.Post("/cms/category/:id", customer_ctlr.UpdateCategoryContent) //categoryId
	// 	r.Delete("/cms/category/:id", customer_ctlr.DeleteCategoryContent)
	//
	// 	// Customer Content By Content Id
	// 	r.Get("/cms/:id", customer_ctlr.GetContentById)
	// 	r.Get("/cms/:id/revisions", customer_ctlr.GetContentRevisionsById)
	//
	// 	//Customer prices
	// 	r.Get("/prices/part/:id", customer_ctlr.GetPricesByPart)         //{id}; id refers to partId
	// 	r.Post("/prices/sale", customer_ctlr.GetSales)                   //{start}{end}{id} -all required params; id refers to customerId
	// 	r.Get("/prices/:id", customer_ctlr.GetPrice)                     //{id}; id refers to {id} refers to customerPriceId
	// 	r.Get("/prices", customer_ctlr.GetAllPrices)                     //returns all {sort=field&direction=dir}
	// 	r.Put("/prices/:id", customer_ctlr.CreateUpdatePrice)            //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
	// 	r.Post("/prices", customer_ctlr.CreateUpdatePrice)               //updates when an id is present; otherwise, creates; {id} refers to customerPriceId
	// 	r.Delete("/prices/:id", customer_ctlr.DeletePrice)               //{id} refers to customerPriceId
	// 	r.Get("/pricesByCustomer/:id", customer_ctlr.GetPriceByCustomer) //{id} refers to customerId; returns CustomerPrices
	//
	// 	r.Post("/:id", customer_ctlr.SaveCustomer)
	// 	r.Delete("/:id", customer_ctlr.DeleteCustomer)
	// 	r.Put("", customer_ctlr.SaveCustomer)
	// })

	// m.Group("/faqs", func(r martini.Router) {
	// 	r.Get("", faq_controller.GetAll)          //get all faqs; takes optional sort param {sort=true} to sort by question
	// 	r.Get("/search", faq_controller.Search)   //takes {question, answer, page, results} - all parameters are optional
	// 	r.Get("/(:id)", faq_controller.Get)       //get by id {id}
	// })

	// m.Group("/geography", func(r martini.Router) {
	// 	r.Get("/states", geography.GetAllStates)
	// 	r.Get("/countries", geography.GetAllCountries)
	// 	r.Get("/countrystates", geography.GetAllCountriesAndStates)
	// })

	// m.Group("/news", func(r martini.Router) {
	// 	r.Get("", news_controller.GetAll)           //get all news; takes optional sort param {sort=title||lead||content||startDate||endDate||active||slug} to sort by question
	// 	r.Get("/titles", news_controller.GetTitles) //get titles!{page, results} - all parameters are optional
	// 	r.Get("/leads", news_controller.GetLeads)   //get leads!{page, results} - all parameters are optional
	// 	r.Get("/search", news_controller.Search)    //takes {title, lead, content, publishStart, publishEnd, active, slug, page, results, page, results} - all parameters are optional
	// 	r.Get("/:id", news_controller.Get)          //get by id {id}
	// })

	// m.Group("/lp", func(r martini.Router) {
	// 	r.Get("/:id", landingPage.Get)
	// })

	// m.Group("/showcase", func(r martini.Router) {
	// 	r.Get("", showcase.GetAllShowcases)
	// 	r.Get("/:id", showcase.GetShowcase)
	// 	r.Post("", showcase.Save)
	// 	// r.Put("/:id", showcase.Save)
	// 	// r.Delete("/:id", showcase.Delete)
	// })

	// m.Group("/testimonials", func(r martini.Router) {
	// 	r.Get("", testimonials.GetAllTestimonials)
	// 	r.Get("/:id", testimonials.GetTestimonial)
	// 	r.Post("", testimonials.Save)
	// 	r.Put("/:id", testimonials.Save)
	// 	r.Delete("/:id", testimonials.Delete)
	// })


	// // ARIES Year/Make/Model/Style
	// m.Post("/vehicle", vehicle.Query)
	// m.Post("/findVehicle", vehicle.GetVehicle)
	// m.Post("/vehicle/inquire", vehicle.Inquire)
	// m.Get("/vehicle/mongo/cols", vehicle.Collections)
	// m.Post("/vehicle/mongo/apps", vehicle.ByCategory)
	// m.Post("/vehicle/mongo/allCollections/category", vehicle.AllCollectionsLookupCategory)
	// m.Post("/vehicle/mongo/allCollections", vehicle.AllCollectionsLookup)
	// m.Post("/vehicle/mongo", vehicle.Lookup)
	// m.Post("/vehicle/mongo/import", vehicle.ImportCsv)
	// m.Get("/vehicle/mongo/all/:collection", vehicle.GetAllCollectionApplications)
	// m.Put("/vehicle/mongo/:collection", vehicle.UpdateApplication)
	// m.Delete("/vehicle/mongo/:collection", vehicle.DeleteApplication)
	// m.Post("/vehicle/mongo/:collection", vehicle.CreateApplication)
	//
	// // CURT Year/Make/Model/Style
	// m.Post("/vehicle/curt", vehicle.CurtLookup)
	//
	// m.Group("/videos", func(r martini.Router) {
	// 	r.Get("/distinct", videos_ctlr.DistinctVideos) //old "videos" table - curtmfg?
	// 	r.Get("/channel/type", videos_ctlr.GetAllChannelTypes)
	// 	r.Get("/channel/type/:id", videos_ctlr.GetChannelType)
	// 	r.Get("/channel", videos_ctlr.GetAllChannels)
	// 	r.Get("/channels", videos_ctlr.GetAllChannels)
	// 	r.Get("/channel/:id", videos_ctlr.GetChannel)
	// 	r.Get("/cdn/type", videos_ctlr.GetAllCdnTypes)
	// 	r.Get("/cdn/type/:id", videos_ctlr.GetCdnType)
	// 	r.Get("/cdn", videos_ctlr.GetAllCdns)
	// 	r.Get("/cdn/:id", videos_ctlr.GetCdn)
	// 	r.Get("/type", videos_ctlr.GetAllVideoTypes)
	// 	r.Get("/type/:id", videos_ctlr.GetVideoType)
	// 	r.Get("", videos_ctlr.GetAllVideos)
	// 	r.Get("/details/:id", videos_ctlr.GetVideoDetails)
	// 	r.Get("/:id", videos_ctlr.Get)
	// })
	//
	// m.Group("/vin", func(r martini.Router) {
	// 	//option 1 - two calls - ultimately returns parts
	// 	r.Get("/configs/:vin", vinLookup.GetConfigs)                    //returns vehicles - user must call vin/vehicle with vehicleID to get parts
	// 	r.Get("/vehicleID/:vehicleID", vinLookup.GetPartsFromVehicleID) //returns an array of parts
	//
	// 	//option 2 - one call - returns vehicles with parts
	// 	r.Get("/:vin", vinLookup.GetParts) //returns vehicles + configs with associates parts -or- an array of parts if only one vehicle config matches
	// })

	// m.Get("/status", func(ctx *fasthttp.RequestCtx) {
	// 	ctx.Response.SetStatusCode(200)
	// 	ctx.Response.SetBody([]byte("running"))
	// })
	//
	// m.Get("/", func(ctx *fasthttp.RequestCtx) {
	// 	ctx.Redirect("http://labs.curtmfg.com", fasthttp.StatusFound)
	// })

	srv := &http.Server{
		Addr:         *listenAddr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on 127.0.0.1%s\n", *listenAddr)
	log.Fatal(srv.ListenAndServe())
}
