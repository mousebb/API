package applicationGuide

import (
	"database/sql"
	"fmt"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"
)

type ApplicationGuide struct {
	ID       int               `json:"id,omitempty" xml:"id,omitempty"`
	URL      string            `json:"url,omitempty" xml:"url,omitempty"`
	Website  Website           `json:"website,omitempty" xml:"website,omitempty"`
	FileType string            `json:"fileType,omitempty" xml:"fileType,omitempty"`
	Category products.Category `json:"category,omitempty" xml:"category,omitempty"`
	Icon     string            `json:"icon,omitempty" xml:"icon,omitempty"`
}

type Website struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	Url         string `json:"url,omitempty" xml:"url,omitempty"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
	BrandIDs    []int  `json:"brandId,omitempty" xml:brandId,omitempty"`
}

const (
	fields = ` ag.url, ag.websiteID, ag.fileType, ag.catID, ag.icon `
)

var (
	getApplicationGuide = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag
										left join Categories as c on c.catID = ag.catID
										Join ApiKeyToBrand as akb on akb.brandID = ag.brandID
										Join ApiKey as ak on akb.keyID = ak.id
										where (ak.api_key = ? && (ag.brandID = ? OR 0=?)) && ag.ID = ? `
	getApplicationGuides = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag
										left join Categories as c on c.catID = ag.catID
										Join ApiKeyToBrand as akb on akb.brandID = ag.brandID
										Join ApiKey as ak on akb.keyID = ak.id
										where ak.api_key = ? && (ag.brandID = ? OR 0=?)
										`
	getApplicationGuidesBySite = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag
										left join Categories as c on c.catID = ag.catID
										Join ApiKeyToBrand as akb on akb.brandID = ag.brandID
										Join ApiKey as ak on akb.keyID = ak.id
										where (ak.api_key = ? && (ag.brandID = ? OR 0=?)) && websiteID = ?`
)

// Get Returns the infomration for the given ApplicationGuide
func (ag *ApplicationGuide) Get(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(getApplicationGuide)
	if err != nil {
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID, ag.ID)

	ch := make(chan ApplicationGuide)
	go populateApplicationGuide(row, ch)
	*ag = <-ch

	if ag.ID == 0 {
		return fmt.Errorf("failed to retrieve application guide")
	}
	return
}

func (ag *ApplicationGuide) GetBySite(ctx *middleware.APIContext) (ags []ApplicationGuide, err error) {

	stmt, err := ctx.DB.Prepare(getApplicationGuidesBySite)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID, ag.Website.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return ags, err
	}

	ch := make(chan []ApplicationGuide)
	go populateApplicationGuides(rows, ch)
	ags = <-ch
	return
}

func populateApplicationGuide(row database.Scanner, ch chan ApplicationGuide) {
	var ag ApplicationGuide
	var catID *int
	var icon []byte
	var catName *string
	err := row.Scan(
		&ag.ID,
		&ag.URL,
		&ag.Website.ID,
		&ag.FileType,
		&catID,
		&icon,
		&catName,
	)
	if err != nil {
		ch <- ag
	}
	if catID != nil {
		ag.Category.CategoryID = *catID
	}
	if catName != nil {
		ag.Category.Title = *catName
	}
	if icon != nil {
		ag.Icon = string(icon[:])
	}
	ch <- ag
	return
}

func populateApplicationGuides(rows *sql.Rows, ch chan []ApplicationGuide) {
	var ag ApplicationGuide
	var ags []ApplicationGuide
	var catID *int
	var icon []byte
	var catName *string
	for rows.Next() {
		err := rows.Scan(
			&ag.ID,
			&ag.URL,
			&ag.Website.ID,
			&ag.FileType,
			&catID,
			&icon,
			&catName,
		)
		if err != nil {
			ch <- ags
		}
		if catID != nil {
			ag.Category.CategoryID = *catID
		}
		if catName != nil {
			ag.Category.Title = *catName
		}
		if icon != nil {
			ag.Icon = string(icon[:])
		}
		ags = append(ags, ag)
	}
	ch <- ags
	return
}
