
CURT Go API
=========
---------

> The new version of the CURT API used the [GoEngine Boilerplate](http://github.com/ninnemana/goengine-boilerplate)
for being Content-Type agnostic to XML and JSON. Some of the best features a listed below:

  - Concurrent MySQL access using [Goroutines](http://golang.org/doc/effective_go.html#concurrency)
  - JSON rendering powered by [encoding/json](http://golang.org/pkg/encoding/json/)
  - XML rendering powered by [encoding/xml](http://golang.org/pkg/encoding/xml/)
  - MySQL Persistence using [mymysql](https://github.com/ziutek/mymysql)
  - ACES Compliant vehicle lookup with product groups
 

--------
Endpoints
---------
---------

#### Vehicle

---

*Get Years*

    GET - http://api.curtmfg.com/v3/vehicle?key=[public api key]

*Get Makes*

    GET - http://api.curtmfg.com/v3/vehicle/2012?key=[public api key]

*Get Models*

    GET - http://api.curtmfg.com/v3/vehicle/2012/Audi?key=[public api key]

*Get SubModels*

    GET - http://api.curtmfg.com/v3/vehicle/2012/Audi/A5?key=[public api key]

*Get Dynamic Configuration Option*

    GET - http://api.curtmfg.com/v3/vehicle/2012/Audi/A5/Cabriolet?key=[public api key]

*Get Next Dynamic Configuration Option*

    GET - http://api.curtmfg.com/v3/vehicle/2012/Audi/A5/Cabriolet/Coupe?key=[public api key]

#### Parts

---

*Get Part by Part #

    GET - http://api.curtmfg.com/v3/part/110003?key=[public api key]

*Reverse Lookup by Part #

    GET - http://api.curtmfg.com/v3/part/110003/vehicles?key=[public api key]

----

#### Categories

---

*Get Category By Category Title

    GET - http://api.curtmfg.com/v3/category/Hitches?key=[public api key]

*Get Category By Category Id

    GET - http://api.curtmfg.com/v3/category/1?key=[public api key]

*Get Top Level Categories

    GET - http://api.curtmfg.com/v3/category?key=[public api key]

*Get Sub-Categories By Category Id

    GET - http://api.curtmfg.com/v3/category/1/subs?key=[public api key]

*Get Sub-Categories By Category Title

    GET - http://api.curtmfg.com/v3/category/Hitches/subs?key=[public api key]

*Get Category Parts By Category Id

    GET - http://api.curtmfg.com/v3/category/3/parts?key=[public api key]
> Keep in mind that the Get Category Parts endpoint implements paging. Below are example endpoints to help demonstrate implementing the pager.

*Get Category Parts By Category Id with Paging

    GET - http://api.curtmfg.com/v3/category/3/parts/2/20?key=[public api key]
> In the above example 2 references the second "page" and 20 in the total count returned.

*Get Category Parts By Category Title

    GET - http://api.curtmfg.com/v3/category/Class I Trailer Hitches/parts?key=[public api key]
> Keep in mind that the Get Category Parts endpoint implements paging. Below are example endpoints to help demonstrate implementing the pager.

*Get Category Parts By Category Title with Paging

    GET - http://api.curtmfg.com/v3/category/Class I Trailer Hitches/parts/2/20?key=[public api key]
> In the above example 2 references the second "page" and 20 in the total count returned.

Philoshopy
-

> This version if the API is meant to focus on data quantity while maintaining, if not improving performance, by leveraging concurrency. We would like the client to have the ability to make fewer requests to the API Server and be provided with a larger amount of data in the response.

Deployment
-

We will be using Capistrano for deployment, you can see the capistrano configuration in the [deploy.rb](https://github.com/curt-labs/GoAPI/blob/master/config/deploy.rb).

Capistrano is configured to pull from github, so you will need to commit and push to the master branch for changes to be pulled.

Once everything has been updated on Github you can run the following command to deploy to curt-api-server1.cloudapp.net and curt-api-server2.cloudapp.net

    cap deploy

Contributors
-
* Alex Ninneman
    * [Github](http://github.com/ninnemana)
    * [Twitter](https://twitter.com/ninnemana)
* Jessica Janiuk
    * [Github](http://github.com/janiukjf)
    * [Twitter](http://twitter.com/janiukjf)



License
-

MIT

*Free Software, Fuck Yeah!*
  
    