package customerContent

import (
	"database/sql"
	"time"

	"github.com/curt-labs/API/models/customer"
)

type Content struct {
	Id          int
	Text        string
	Added       time.Time
	Modified    time.Time
	ContentType ContentType
	Hidden      bool
	Customer    *customer.Customer
	User        *customer.User
	Revisions   []ContentRevision
}

type ContentType struct {
	Id        int
	Type      string
	AllowHtml bool
}

type ContentRevision struct {
	Id             int
	User           customer.User
	Customer       customer.Customer
	OldText        string
	NewText        string
	Date           time.Time
	ChangeType     string
	ContentId      int
	OldContentType ContentType
	NewContentType ContentType
}

// Retrieves specific part content for this customer
func GetPartContent(db *sql.DB, partID int) (content []Content, err error) {
	content = make([]Content, 0) // initializer

	return content, err
}
