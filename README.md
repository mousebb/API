
CURT API v3
=========
---------

[!codeship status](https://codeship.com/projects/41c48f00-68a6-0134-316b-625c6894da4e/status?branch=master)
[![wercker status](https://app.wercker.com/status/f8512c3d160ff9b1d198ea38b2e9568b/s "wercker status")](https://app.wercker.com/project/bykey/f8512c3d160ff9b1d198ea38b2e9568b)

> The 3rd iteration of the CURT Manufacturing API. This is developed using MongoDB
as a storage service that is being distributed from our master data source in
real time. All CURT assets are developed using the endpoints available in this series
of endpoints.

  - Concurrent MySQL access using [Goroutines](http://golang.org/doc/effective_go.html#concurrency)
  - JSON rendering powered by [encoding/json](http://golang.org/pkg/encoding/json/)
  - XML rendering powered by [encoding/xml](http://golang.org/pkg/encoding/xml/)
  - MongoDB Persistence using [mgo](https://labix.org/mgo)
  - Multi-Brand Support
  - ACES Compliant vehicle lookup with product groups
  - Custom lookup to support each brands vehicle requirements
  - Tailored product content
  - Category Hierarchy with products matched all the way down


Philoshopy
-

> This version if the API is meant to focus on data quantity while maintaining, if not improving performance, by leveraging concurrency. We would like the client to have the ability to make fewer requests to the API Server and be provided with a larger amount of data in the response.

Deployment
-

Deployment will be done using the master branch on Github. Continuous Delivery/Integration
will go through [Jenkins](https://jenkins.io/). The application will run in a
cluster of Docker containers that is orchestrated by a Kubernetes arbiter.

Contributors
-
* Alex Ninneman
    * [Github](http://github.com/ninnemana)
* David Vaini
    * [Github](https://github.com/DavidVaini)
* John Shenk
    * [Github](https://github.com/stinkyfingers)
* Broc Seigneurie
    * [Github](https://github.com/baseigneurie)

License
-

MIT

*Free-ish software? Open, we'll call it open software. Oh and :beers:*
