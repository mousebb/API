package customer

import (
	"database/sql"
	"fmt"
)

var (
	getCustIDFromCustomerID = `select cust_id from Customer where customerID = ? limit 1`
)

func GetCustomerPrice(db *sql.DB, partID int) (float64, error) {
	return 0, nil
}

func GetCustomerCartReference(db *sql.DB, partID int) (int, error) {
	return 0, nil
}

func numberToID(db *sql.DB, customerID int) (int, error) {

	stmt, err := db.Prepare(getCustIDFromCustomerID)
	if err != nil {
		return 0, err
	}

	var custID *int
	err = stmt.QueryRow(customerID).Scan(&custID)
	if err != nil || custID == nil || *custID == 0 {
		return 0, fmt.Errorf("failed to retrieve ID %s", err.Error())
	}

	return *custID, nil
}
