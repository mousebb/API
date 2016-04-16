package apiKeyType

import (
	"database/sql"
	"time"

	"github.com/curt-labs/API/helpers/database"
)

var (
	// GetAllTypes Database query to retrieve all APIKeyType
	GetAllTypes = "SELECT id, type, date_added FROM ApiKeyType order by type"
)

// KeyType Declares the type reference for an API key (Public, Private, Authentication).
type KeyType struct {
	ID        string    `bson:"-" json:"-" xml:"-"`
	Type      string    `bson:"type" json:"type" xml:"type"`
	DateAdded time.Time `bson:"dateAdded" json:"dateAdded" xml:"dateAdded"`
}

const (
	timeFormat = "2006-01-02 03:04:05"
)

// GetAllKeyTypes Returns a list of all the API key types in the database.
func GetAllKeyTypes(tx *sql.Tx) (as []KeyType, err error) {

	stmt, err := tx.Prepare(GetAllTypes)
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err := stmt.Query() //returns *sql.Rows
	if err != nil {
		return
	}

	var a *KeyType
	for res.Next() {
		a, err = scan(res)
		if err != nil {
			return as, err
		}
		as = append(as, *a)
	}
	defer res.Close()

	return as, err
}

func scan(s database.Scanner) (*KeyType, error) {
	a := &KeyType{}
	err := s.Scan(&a.ID, &a.Type, &a.DateAdded)
	return a, err
}
