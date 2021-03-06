package apiKeyType

import (
	"database/sql"
	"github.com/curt-labs/API/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	getApiKeyType     = "SELECT id, type, date_added FROM ApiKeyType WHERE id = ? "
	getAllApiKeyTypes = "SELECT id, type, date_added FROM ApiKeyType "
	getKeyByDateType  = "SELECT id FROM ApiKeyType WHERE type = ?  AND date_added = ?"
	createApiKeyType  = "INSERT INTO ApiKeyType (id, type, date_added) VALUES (UUID(),?,?)"
	deleteApiKeyType  = "DELETE FROM ApiKeyType WHERE id = ? "
)

type ApiKeyType struct {
	ID        string    `json:"_id" xml:"id"`
	Type      string    `json:"type" xml:"type"`
	DateAdded time.Time `json:"dateAdded" xml:"dateAdded"`
}

type Scanner interface {
	Scan(...interface{}) error
}

const (
	timeFormat = "2006-01-02 03:04:05"
)

func (a *ApiKeyType) Get() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getApiKeyType)
	if err != nil {
		return
	}
	defer stmt.Close()
	res := stmt.QueryRow(a.ID)
	a, err = ScanKey(res)

	return
}

func GetAllApiKeyTypes() (as []ApiKeyType, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(getAllApiKeyTypes)
	if err != nil {
		return
	}
	defer stmt.Close()
	res, err := stmt.Query() //returns *sql.Rows
	if err != nil {
		return
	}

	for res.Next() {
		a, err := ScanKey(res)
		if err != nil {
			return as, err
		}
		as = append(as, *a)
	}
	defer res.Close()
	return as, err
}

func ScanKey(s Scanner) (*ApiKeyType, error) {
	a := &ApiKeyType{}
	err := s.Scan(&a.ID, &a.Type, &a.DateAdded)
	return a, err
}

func (a *ApiKeyType) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare(createApiKeyType)
	if err != nil {
		return
	}
	defer stmt.Close()
	added := time.Now().Format(timeFormat)
	_, err = stmt.Exec(a.Type, added)
	if err != nil {
		return
	}

	stmt, err = db.Prepare(getKeyByDateType)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(a.Type, added).Scan(&a.ID)
	if err != nil {
		return err
	}
	return
}

func (a *ApiKeyType) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteApiKeyType)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.ID)
	if err != nil {
		return err
	}
	return
}
