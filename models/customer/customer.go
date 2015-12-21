package customer

import "database/sql"

func GetCustomerPrice(db *sql.DB, partID int) (float64, error) {
	return 0, nil
}

func GetCustomerCartReference(db *sql.DB, partID int) (int, error) {
	return 0, nil
}
