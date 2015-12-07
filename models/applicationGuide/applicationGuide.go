package applicationGuide

import (
	"database/sql"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"
	"github.com/curt-labs/API/models/site"
)

type ApplicationGuide struct {
	ID       int               `json:"id,omitempty" xml:"id,omitempty"`
	Url      string            `json:"url,omitempty" xml:"url,omitempty"`
	Website  site.Website      `json:"website,omitempty" xml:"website,omitempty"`
	FileType string            `json:"fileType,omitempty" xml:"fileType,omitempty"`
	Category products.Category `json:"category,omitempty" xml:"category,omitempty"`
	Icon     string            `json:"icon,omitempty" xml:"icon,omitempty"`
}

const (
	fields = ` ag.url, ag.websiteID, ag.fileType, ag.catID, ag.icon `
)

var (
	createApplicationGuide = `insert into ApplicationGuides (url, websiteID, fileType, catID, icon, brandID) values (?,?,?,?,?,?)`
	deleteApplicationGuide = `delete from ApplicationGuides where ID = ?`
	getApplicationGuide    = `select ag.ID, ` + fields + `, c.catTitle from ApplicationGuides as ag
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

func (ag *ApplicationGuide) Create(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(createApplicationGuide)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(ag.Url, ag.Website.ID, ag.FileType, ag.Category.CategoryID, ag.Icon, ctx.DataContext.BrandID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	ag.ID = int(id)
	return
}
func (ag *ApplicationGuide) Delete(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(deleteApplicationGuide)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ag.ID)

	return err
}

func populateApplicationGuide(row *sql.Row, ch chan ApplicationGuide) {
	var ag ApplicationGuide
	var catID *int
	var icon []byte
	var catName *string
	err := row.Scan(
		&ag.ID,
		&ag.Url,
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
			&ag.Url,
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
