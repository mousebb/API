package customer

import (
	"net/url"
	"time"

	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/geography"
	"gopkg.in/mgo.v2/bson"
)

// Customer Holds everything that we store about a customer account.
type Customer struct {
	Identifier     bson.ObjectId `bson:"_id" json:"-" xml:"-"`
	ID             int           `bson:"id" json:"id" xml:"id,attr"`
	CustomerNumber int           `bson:"customerNumber" json:"customerNumber" xml:"customerNumber,attr"`
	Name           string        `bson:"name" json:"name" xml:"name"`
	EmailAddress   string        `bson:"email" json:"email" xml:"email"`
	Address        Address       `bson:"address" json:"address" xml:"address"`
	Phone          string        `bson:"phone" json:"phone" xml:"contact>phone"`
	Fax            string        `bson:"fax" json:"fax" xml:"contact>fax"`
	ContactPerson  string        `bson:"contactPerson" json:"contactPerson" xml:"contact>contactPerson"`

	Parent         *Customer `bson:"parentAccount" json:"parentAccount" xml:"parentAccount"`
	Website        *url.URL  `bson:"website" json:"website" xml:"online>website"`
	Elocal         *url.URL  `bson:"elocal" json:"elocal" xml:"online>elocal"`
	SearchEndpoint *url.URL  `bson:"searchURL" json:"searchURL" xml:"online>searchURL"`
	Logo           *url.URL  `bson:"logo" json:"logo" xml:"online>logo"`

	Types    []Type    `bson:"types" json:"type" xml:"type"`
	Tiers    []Tier    `bson:"tiers" json:"tier" xml:"tier"`
	MapIcons []MapIcon `bson:"mapIcons" json:"mapIcons" xml:"mapIcons"`

	Brands   []brand.Brand       `bson:"brands" json:"brands" xml:"brands"`
	Mapics   MapicsCode          `bson:"mapics" json:"mapics" xml:"mapics"`
	SalesRep SalesRepresentative `bson:"salesRep" json:"salesRep" xml:"salesRep"`

	Users     []User     `bson:"users" json:"users" xml:"users"`
	Locations []Location `bson:"locations" json:"locations" xml:"locations"`

	Accounts []ComnetAccount `bson:"accounts" json:"accounts" xml:"accounts"`

	DummyAccount bool `bson:"dummyAccount" json:"dummyAccount" xml:"dummyAccount,attr"`
}

// Coordinates Geograhical spatial data
type Coordinates struct {
	Latitude  float64 `bson:"latitude" json:"latitude" xml:"latitude"`
	Longitude float64 `bson:"longtiude" json:"longitude" xml:"longitude"`
}

// MapicsCode Defines the designation of this customer in Mapics.
type MapicsCode struct {
	Code        string `bson:"code" json:"code" xml:"code,attr"`
	Description string `bson:"description" json:"description" xml:"description,attr"`
}

// SalesRepresentative The sales rep that is affiliated with a given customer.
type SalesRepresentative struct {
	ID   int    `bson:"id" json:"id" xml:"id,attr"`
	Code string `bson:"code" json:"code" xml:"code,attr"`
	Name string `bson:"name" json:"name" xml:"name,attr"`
}

// Address Describes the geopgraphic location by mailing address and
// assoicated spatial data.
type Address struct {
	StreetAddress  string          `bson:"streetAddress" json:"streetAddress" xml:"streetAddress"`
	StreetAddress2 string          `bson:"streetAddress2" json:"streetAddress2" xml:"streetAddress2"`
	City           string          `bson:"city" json:"city" xml:"city"`
	State          geography.State `bson:"state" json:"state" xml:"state"`
	PostalCode     string          `bson:"postalCode" json:"postalCode" xml:"postalCode"`
	Coordinates    Coordinates     `bson:"coordinates" json:"coordinates" xml:"coordinates"`
}

// Warehouse Describes the geopgraphical location and contact information for
// a warehouse.
type Warehouse struct {
	ID            int     `bson:"id" json:"-" xml:"-"`
	Name          string  `bson:"name" json:"name" xml:"name"`
	Code          string  `bson:"code" json:"code" xml:"code"`
	Address       Address `bson:"address" json:"address" xml:"address"`
	TollFreePhone string  `bson:"tollFreePhone" json:"tollFreePhone" xml:"tollFreePhone"`
	Fax           string  `bson:"fax" json:"fax" xml:"fax"`
	LocalPhone    string  `bson:"localPhone" json:"localPhone" xml:"localPhone"`
	Manager       string  `bson:"manager" json:"manager" xml:"manager"`
}

// ComnetAccount Describes the Comnet Integrations for a User.
type ComnetAccount struct {
	Credentials   *ComnetCredential `bson:"credentials,omitempty" json:"credentials,omitempty" xml:"credentials,omitempty"`
	AccountNumber string            `bson:"accountNumber" json:"accountNumber" xml:"accountNumber"`
	FreightLimit  float64           `bson:"freightLimit" json:"freightLimit" xml:"freightLimit"`
	Warehouse     Warehouse         `bson:"warehouse" json:"warehouse" xml:"warehouse"`
	Type          ComnetAccountType `bson:"type" json:"type" xml:"type"`
	Status        string            `bson:"status" json:"status,omitempty" xml:"status,omitempty"`
}

// ComnetAccountType Describes which ComNET service to integrate
// with. (ARIES Interior, ARIES Exterior, CURT)
type ComnetAccountType struct {
	ID    int      `bson:"id" json:"-" xml:"-"`
	Title string   `bson:"title" json:"title" xml:"title"`
	URL   *url.URL `bson:"url" json:"url" xml:"url"`
}

// ComnetCredential Object to use when authenticating with
// a ComNET account.
type ComnetCredential struct {
	Username string `bson:"username" json:"username" xml:"username"`
	Password string `bson:"password" json:"password" xml:"password"`
}

// APIKey Used to authenticate with CURT Web Services.
type APIKey struct {
	Key       string        `bson:"key" json:"key" xml:"key,attr"`
	Type      APIKeyType    `bson:"type" json:"type" xml:"type,attr"`
	DateAdded time.Time     `bson:"dateAdded" json:"dateAdded" xml:"dateAdded,attr"`
	Brands    []brand.Brand `bson:"brands" json:"brands" xml:"brands"`
}

// APIKeyType Describes the type of API key. (Public, Private, etc)
type APIKeyType struct {
	Type      string    `bson:"type" json:"type" xml:"type"`
	DateAdded time.Time `bson:"dateAdded" json:"dateAdded" xml:"dateAdded"`
}

// Location Physical location for a Customer
type Location struct {
	ID              int     `bson:"id" json:"-" xml:"-"`
	Name            string  `bson:"name" json:"name" xml:"name,attr"`
	Address         Address `bson:"address" json:"address" xml:"address"`
	Email           string  `bson:"email" json:"email" xml:"email,attr"`
	Phone           string  `bson:"phone" json:"phone" xml:"phone,attr"`
	Fax             string  `bson:"fax" json:"fax" xml:"fax,attr"`
	ContactPerson   string  `bson:"contactPerson" json:"contactPerson" xml:"contactPerson,attr"`
	PrimaryLocation bool    `bson:"primaryLocation" json:"primaryLocation" xml:"primaryLocation,attr"`
	ShippingDefault bool    `bson:"shippingDefault" json:"shippingDefault" xml:"shippingDefault,attr"`
}

// MapIcon Image icons to use for Google Maps
type MapIcon struct {
	Icon   *url.URL    `bson:"icon" json:"icon" xml:"icon"`
	Shadow *url.URL    `bson:"shadow" json:"shadow" xml:"shadow"`
	Brand  brand.Brand `bson:"brand" json:"brand" xml:"brand"`
}

// Tier Declares what level of customer (Silver, Gold, Platinum).
type Tier struct {
	Tier  string      `bson:"tier" json:"tier" xml:"tier,attr"`
	Sort  int         `bson:"sort" json:"sort" xml:"sort,attr"`
	Brand brand.Brand `bson:"brand" json:"brand" xml:"brand"`
}

// Type Declares whether a customer is an Online, Retail, or Installer seller.
type Type struct {
	Type   string      `bson:"type" json:"type" xml:"type,attr"`
	Online bool        `bson:"online" json:"online" xml:"online,attr"`
	Show   bool        `bson:"show" json:"show" xml:"show,attr"`
	Label  string      `bson:"label" json:"label" xml:"label,attr"`
	Brand  brand.Brand `bson:"brand" json:"brand" xml:"brand"`
}

// User An authenticated user that is affiliated with a single Customer.
type User struct {
	ID   string `bson:"id" json:"id" xml:"id,attr"`
	Name string `bson:"name" json:"name" xml:"name,attr"`

	// Redundant if nested within a Customer
	// but valuable if retrieved on it's own.
	CustomerNumber int       `bson:"customerNumber" json:"customerNumber" xml:"customerNumber,attr"`
	Email          string    `bson:"email" json:"email" xml:"email,attr"`
	Password       string    `bson:"password" json:"-" xml:"-"`
	DateAdded      time.Time `bson:"dateAdded" json:"dateAdded" xml:"dateAdded,attr"`

	Location  Location `bson:"location" json:"location" xml:"location"`
	SuperUser bool     `bson:"superUser" json:"superUser" xml:"superUser,attr"`

	Keys []APIKey `bson:"keys" json:"keys" xml:"keys"`

	ComnetAccounts []ComnetAccount `bson:"comnetAccounts" json:"comnetAccounts" xml:"comnetAccounts"`
}

type Scanner interface {
	Scan(...interface{}) error
}

type StateRegion struct {
	Id           int          `json:"id,omitempty" xml:"id,omitempty"`
	Name         string       `json:"name,omitempty" xml:"name,omitempty"`
	Abbreviation string       `json:"abbreviation,omitempty" xml:"abbreviation,omitempty"`
	Count        int          `json:"count,omitempty" xml:"count,omitempty"`
	Polygons     []MapPolygon `json:"polygon,omitempty" xml:"polygon,omitempty"`
}

type MapPolygon struct {
	Id          int           `json:"id,omitempty" xml:"id,omitempty"`
	Coordinates []Coordinates `json:"coordinates,omitempty" xml:"coordinates,omitempty"`
}
