package customer

import (
	"encoding/json"
	"strconv"

	"github.com/curt-labs/API/helpers/conversions"
	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/geography"
)

type CustomerLocation struct {
	Id              int             `json:"id,omitempty" xml:"id,omitempty"`
	Name            string          `json:"name,omitempty" xml:"name,omitempty"`
	Email           string          `json:"email,omitempty" xml:"email,omitempty"`
	Address         string          `json:"address,omitempty" xml:"address,omitempty"`
	City            string          `json:"city,omitempty" xml:"city,omitempty"`
	PostalCode      string          `json:"postalCode,omitempty" xml:"postalCode,omitempty"`
	State           geography.State `json:"state,omitempty" xml:"state,omitempty"`
	Phone           string          `json:"phone,omitempty" xml:"phone,omitempty"`
	Fax             string          `json:"fax,omitempty" xml:"fax,omitempty"`
	Coordinates     Coordinates     `json:"coords,omitempty" xml:"coords,omitempty"`
	CustomerId      int             `json:"customerId,omitempty" xml:"customerId,omitempty"`
	ContactPerson   string          `json:"contactPerson,omitempty" xml:"contactPerson,omitempty"`
	IsPrimary       bool            `json:"isPrimary,omitempty" xml:"isPrimary,omitempty"`
	ShippingDefault bool            `json:"shippingDefault,omitempty" xml:"shippingDefault,omitempty"`
}

var (
	getLocation  = "SELECT locationID, name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault FROM CustomerLocations WHERE locationID= ? "
	getLocations = `SELECT cl.locationID, cl.name, cl.address, cl.city, cl.stateID, cl.email,cl.phone, cl.fax, cl.latitude, cl.longitude, cl.cust_id, cl.contact_person, cl.isprimary, cl.postalCode, cl.ShippingDefault
			FROM CustomerLocations as cl
			join CustomerToBrand as ctb on ctb.cust_id = cl.cust_id
			join ApiKeyToBrand as akb on akb.brandID = ctb.brandID
			join ApiKey as ak on ak.id = akb.keyID
			where ak.api_key = ? && (ctb.BrandID = ? or 0 = ?)`
	createLocation = "INSERT INTO CustomerLocations (name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	updateLocation = "UPDATE CustomerLocations SET name = ?, address = ?,  city = ?,  stateID = ?, email = ?,  phone = ?,  fax = ?,  latitude = ?,  longitude = ?,  cust_id = ?, contact_person = ?,  isprimary = ?, postalCode = ?, ShippingDefault = ? WHERE locationID = ?"
	deleteLocation = "DELETE FROM CustomerLocations WHERE locationID = ? "
)

func (l *CustomerLocation) Get(ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(getLocation)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var name, address, city, email, phone, fax, contactPerson, postal []byte
	err = stmt.QueryRow(l.Id).Scan(
		&l.Id,
		&name,
		&address,
		&city,
		&l.State.Id,
		&email,
		&phone,
		&fax,
		&l.Coordinates.Latitude,
		&l.Coordinates.Longitude,
		&l.CustomerId,
		&contactPerson,
		&l.IsPrimary,
		&postal,
		&l.ShippingDefault,
	)
	if err != nil {
		return err
	}
	l.Name, err = conversions.ByteToString(name)
	l.Address, err = conversions.ByteToString(address)
	l.City, err = conversions.ByteToString(city)
	l.Email, err = conversions.ByteToString(email)
	l.Phone, err = conversions.ByteToString(phone)
	l.Fax, err = conversions.ByteToString(fax)
	l.ContactPerson, err = conversions.ByteToString(contactPerson)
	l.PostalCode, err = conversions.ByteToString(postal)
	if err != nil {
		return err
	}

	return err
}

func GetAllLocations(ctx *middleware.APIContext) (CustomerLocations, error) {
	var ls CustomerLocations
	var err error
	redis_key := "customers:locations:" + ctx.DataContext.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ls)
		return ls, err
	}

	stmt, err := ctx.DB.Prepare(getLocations)
	if err != nil {
		return ls, err
	}
	defer stmt.Close()

	res, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	if err != nil {
		return ls, err
	}
	defer res.Close()

	var name, address, city, email, phone, fax, contactPerson, postal []byte
	for res.Next() {
		var l CustomerLocation
		err = res.Scan(
			&l.Id,
			&name,
			&address,
			&city,
			&l.State.Id,
			&email,
			&phone,
			&fax,
			&l.Coordinates.Latitude,
			&l.Coordinates.Longitude,
			&l.CustomerId,
			&contactPerson,
			&l.IsPrimary,
			&postal,
			&l.ShippingDefault,
		)
		if err != nil {
			return ls, err
		}
		l.Name, err = conversions.ByteToString(name)
		l.Address, err = conversions.ByteToString(address)
		l.City, err = conversions.ByteToString(city)
		l.Email, err = conversions.ByteToString(email)
		l.Phone, err = conversions.ByteToString(phone)
		l.Fax, err = conversions.ByteToString(fax)
		l.ContactPerson, err = conversions.ByteToString(contactPerson)
		l.PostalCode, err = conversions.ByteToString(postal)
		if err != nil {
			return ls, err
		}
		ls = append(ls, l)
	}

	go redis.Setex(redis_key, ls, 86400)

	return ls, err
}

func (l *CustomerLocation) Create(ctx *middleware.APIContext) error {

	tx, err := ctx.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(createLocation)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(
		l.Name,
		l.Address,
		l.City,
		l.State.Id,
		l.Email,
		l.Phone,
		l.Fax,
		l.Coordinates.Latitude,
		l.Coordinates.Longitude,
		l.CustomerId,
		l.ContactPerson,
		l.IsPrimary,
		l.PostalCode,
		l.ShippingDefault,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	l.Id = int(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	go redis.Delete("customers:locations:" + ctx.DataContext.BrandString)
	go redis.Delete("customerLocations:" + strconv.Itoa(l.CustomerId))
	return err
}

func (l *CustomerLocation) Update(ctx *middleware.APIContext) error {

	tx, err := ctx.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateLocation)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		l.Name,
		l.Address,
		l.City,
		l.State.Id,
		l.Email,
		l.Phone,
		l.Fax,
		l.Coordinates.Latitude,
		l.Coordinates.Longitude,
		l.CustomerId,
		l.ContactPerson,
		l.IsPrimary,
		l.PostalCode,
		l.ShippingDefault,
		l.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	go redis.Delete("customers:locations:" + ctx.DataContext.BrandString)
	go redis.Delete("customerLocations:" + strconv.Itoa(l.CustomerId))
	return err
}

func (l *CustomerLocation) Delete(ctx *middleware.APIContext) error {

	tx, err := ctx.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deleteLocation)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(l.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	go redis.Delete("customers:locations:" + ctx.DataContext.BrandString)
	go redis.Delete("customerLocations:" + strconv.Itoa(l.CustomerId))
	return err
}
