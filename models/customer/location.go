package customer

import (
	"database/sql"
	"fmt"
	"strings"
)

var (
	getLocationByAddress = `select locationID from CustomerLocations cl
                            join Customer c on cl.cust_id = c.cust_id
                            where lower(cl.address) = lower(?) && c.customerID = ?
                            limit 1`
	insertLocation = `insert into CustomerLocations(name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault)
                        values(?, ?, ?, ?, ?, ?, ?, ?, ?, (
                        	select cust_id from Customer where customerID = ? limit 1
                        ), ?, ?, ?, ?)`
)

// Exists Uses the Address.StreetAddress and Address.StreetAddress2
// to comapre against the `address` column on `CustomerLocations`
// to determine if a location already
func (l *Location) Exists(db *sql.DB, customerID int) (bool, error) {

	stmt, err := db.Prepare(getLocationByAddress)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	qry := strings.TrimSpace(fmt.Sprintf("%s %s", l.Address.StreetAddress, l.Address.StreetAddress2))

	var locID *int
	err = stmt.QueryRow(qry, customerID).Scan(&locID)
	if err != nil || locID == nil || *locID == 0 {
		if err == sql.ErrNoRows {
			err = nil
		}
		return false, err
	}

	l.ID = *locID

	return true, nil
}

func (l *Location) insert(tx *sql.Tx, customerID int) error {

	stmt, err := tx.Prepare(insertLocation)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(
		l.Name,
		l.Address,
		l.Address.City,
		l.Address.State.Id,
		l.Email,
		l.Phone,
		l.Fax,
		l.Address.Coordinates.Latitude,
		l.Address.Coordinates.Longitude,
		customerID,
		l.ContactPerson,
		l.PrimaryLocation,
		l.Address.PostalCode,
		l.ShippingDefault,
	)
	if err != nil {
		return err
	}

	locationID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	l.ID = int(locationID)

	return nil
}
