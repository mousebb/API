package brand

import (
	"database/sql"
	"errors"
	"net/url"
)

var (
	brandFields           = `ID, name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID`
	getAllBrandsStmt      = `select ` + brandFields + ` from Brand`
	getBrandStmt          = `select ` + brandFields + ` from Brand where ID = ?`
	insertBrandStmt       = `insert into Brand(name, code, logo, logoAlt, formalName, longName, primaryColor, autocareID) values (?,?,?,?,?,?,?,?)`
	updateBrandStmt       = `update Brand set name = ?, code = ?, logo = ?, logoAlt = ?, formalName = ?, longName = ?, primaryColor = ?, autocareID = ? where ID = ?`
	deleteBrandStmt       = `delete from Brand where ID = ?`
	getCustomerUserBrands = `select b.ID, b.name, b.code, b.logo, b.logoAlt, b.formalName, b.longName, b.primaryColor, b.autocareID
								from Brand as b
								join CustomerToBrand as ctb on ctb.BrandID = b.ID
								join Customer as c on c.cust_id = ctb.cust_id
								where c.cust_id = ?
								group by b.ID`
	getAllWebsitesStmt = `select w.ID, w.description, w.url, wb.brandID from Website as w
							join WebsiteToBrand as wb on w.ID = wb.WebsiteID
							order by wb.brandID, w.ID`
	getBrandWebsitesStmt = `select w.ID, w.description, w.url, wb.brandID from Website as w
							join WebsiteToBrand as wb on w.ID = wb.WebsiteID
							where wb.brandID = ?
							order by w.ID`
)

type Brand struct {
	ID            int       `json:"id" xml:"id,attr"`
	Name          string    `json:"name" xml:"name,attr"`
	Code          string    `json:"code" xml:"code,attr"`
	Logo          *url.URL  `json:"logo" xml:"logo,attr"`
	LogoAlternate *url.URL  `json:"logo_alternate" xml:"logo_alternate,attr"`
	FormalName    string    `json:"formal_name" xml:"formal_name,attr"`
	LongName      string    `json:"long_name" xml:"long_name,attr"`
	PrimaryColor  string    `json:"primary_color" xml:"primary_color,attr"`
	AutocareID    string    `json:"autocareId" xml:"autocareId,attr"`
	Websites      []Website `json:"websites" xml:"websites"`
}

type Website struct {
	ID          int      `json:"id" xml:"id,attr"`
	Description string   `json:"description" xml:"description"`
	URL         *url.URL `json:"url" xml:"url"`
	BrandID     int      `json:"brand_id" xml:"brand_id"`
}

type Scanner interface {
	Scan(...interface{}) error
}

func ScanBrand(res Scanner) (Brand, error) {
	var logo, logoAlt, formal, long, primary, autocare *string
	var b Brand
	err := res.Scan(&b.ID, &b.Name, &b.Code, &logo, &logoAlt, &formal, &long, &primary, &autocare)
	if err != nil {
		return b, err
	}
	if logo != nil {
		b.Logo, err = url.Parse(*logo)
		if err != nil {
			return b, err
		}
	}
	if logoAlt != nil {
		b.LogoAlternate, err = url.Parse(*logoAlt)
		if err != nil {
			return b, err
		}
	}
	if formal != nil {
		b.FormalName = *formal
	}
	if long != nil {
		b.LongName = *long
	}
	if primary != nil {
		b.PrimaryColor = *primary
	}
	if autocare != nil {
		b.AutocareID = *autocare
	}
	return b, err
}

func GetAllBrands(db *sql.DB) (brands []Brand, err error) {

	stmt, err := db.Prepare(getAllBrandsStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var b Brand
		b, err = ScanBrand(rows)
		if err != nil {
			return
		}
		brands = append(brands, b)
	}
	defer rows.Close()

	return
}

func (b *Brand) Get(db *sql.DB) error {
	if b.ID == 0 {
		return errors.New("Invalid Brand ID")
	}

	stmt, err := db.Prepare(getBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res := stmt.QueryRow(b.ID)
	*b, err = ScanBrand(res)
	if err != nil {
		return err
	}
	return nil
}

func getWebsites(brandID int, db *sql.DB) ([]Website, error) {
	sites := make([]Website, 0)

	var err error
	var rows *sql.Rows

	if brandID > 0 {
		stmt, err := db.Prepare(getBrandWebsitesStmt)
		if err != nil {
			return sites, err
		}
		defer stmt.Close()

		rows, err = stmt.Query(brandID)
	} else {
		stmt, err := db.Prepare(getAllWebsitesStmt)
		if err != nil {
			return sites, err
		}
		defer stmt.Close()

		rows, err = stmt.Query()
	}

	if err != nil {
		return sites, err
	}

	for rows.Next() {
		var s Website
		var u *string
		err = rows.Scan(&s.ID, &s.Description, &u, &s.BrandID)
		if err != nil || u == nil {
			continue
		}

		s.URL, err = url.Parse(*u)
		if err != nil {
			continue
		}

		sites = append(sites, s)
	}

	return sites, nil
}

func GetUserBrands(id int, db *sql.DB) ([]Brand, error) {
	brands := make([]Brand, 0)
	var err error

	stmt, err := db.Prepare(getCustomerUserBrands)
	if err != nil {
		return brands, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return brands, err
	}

	sites, err := getWebsites(0, db)
	if err != nil {
		return brands, err
	}

	indexedSites := make(map[int][]Website, 0)
	for _, site := range sites {
		if _, ok := indexedSites[site.BrandID]; !ok {
			indexedSites[site.BrandID] = make([]Website, 0)
		}

		indexedSites[site.BrandID] = append(indexedSites[site.BrandID], site)
	}

	for rows.Next() {
		var b Brand
		b, err = ScanBrand(rows)
		if err != nil {
			return brands, err
		}

		b.Websites = indexedSites[b.ID]
		brands = append(brands, b)
	}
	return brands, nil
}
