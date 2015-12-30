package landingPage

import (
	"net/url"
	"time"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	_ "github.com/go-sql-driver/mysql"

	"github.com/russross/blackfriday"
)

type LandingPage struct {
	Id, WebsiteID                               int
	NewWindow                                   bool
	Name, PageContent, LinkClasses              string
	HtmlContent                                 string
	StartDate, EndDate                          time.Time
	ConversionID, ConversionLabel, MenuPosition string
	Url                                         *url.URL
	LandingPageDatas                            []LandingPageData
	LandingPageImages                           []LandingPageImage
}

type LandingPageData struct {
	Id, LandingPageID  int
	DataKey, DataValue string
}

type LandingPageImage struct {
	Id, LandingPageID, Sort int
	Url                     *url.URL
}

var (
	GetLandingPageByID = `select lp.id, lp.name, lp.startDate, lp.endDate, lp.url, lp.pageContent, lp.linkClasses, lp.conversionID, lp.conversionLabel, lp.newWindow, lp.menuPosition, lp.websiteID from LandingPage as lp
							Join WebsiteToBrand as wub on wub.WebsiteID = lp.websiteID
							Join ApiKeyToBrand as akb on akb.brandID = wub.brandID
							Join ApiKey as ak on akb.keyID = ak.id
							where lp.id = ? && lp.startDate <= NOW() && lp.endDate >= NOW() && (ak.api_key = ? && (wub.brandID = ? OR 0=?))
							limit 1`
	GetLandingPageImagesStmt = `SELECT lpi.id, lpi.landingPageID, lpi.url, lpi.sort from LandingPageImages as lpi
									WHERE lpi.landingPageID = ?
									ORDER BY lpi.sort asc`
	GetLandingPageDatasStmt = `SELECT lpd.id, lpd.landingPageID, lpd.dataKey, lpd.dataValue from LandingPageData as lpd
									WHERE lpd.landingPageID = ?`
)

func (lp *LandingPage) Get(ctx *middleware.APIContext) (err error) {
	stmt, err := ctx.DB.Prepare(GetLandingPageByID)
	if err != nil {
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(lp.Id, ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)

	err = lp.pageScan(ctx, row)
	if err != nil {
		return
	}
	lp.HtmlContent = string(blackfriday.MarkdownBasic([]byte(lp.PageContent)))
	return
}

func (lp *LandingPage) pageScan(ctx *middleware.APIContext, s database.Scanner) error {
	var pageContent, linkClasses, conversionID, conversionLabel, urlstr *string

	err := s.Scan(
		&lp.Id,
		&lp.Name,
		&lp.StartDate,
		&lp.EndDate,
		&urlstr,
		&pageContent,
		&linkClasses,
		&conversionID,
		&conversionLabel,
		&lp.NewWindow,
		&lp.MenuPosition,
		&lp.WebsiteID,
	)
	if err != nil {
		return err
	}

	if pageContent != nil {
		lp.PageContent = *pageContent
	}
	if linkClasses != nil {
		lp.LinkClasses = *linkClasses
	}
	if conversionID != nil {
		lp.ConversionID = *conversionID
	}
	if conversionLabel != nil {
		lp.ConversionLabel = *conversionLabel
	}
	if urlstr != nil {
		lp.Url, _ = url.Parse(*urlstr)
	}
	lp.LandingPageDatas, _ = data(ctx, lp.Id)
	lp.LandingPageImages, _ = images(ctx, lp.Id)

	return nil
}

func data(ctx *middleware.APIContext, LandingPageID int) (datas []LandingPageData, err error) {
	stmt, err := ctx.DB.Prepare(GetLandingPageDatasStmt)
	if err != nil {
		return datas, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(LandingPageID)
	if err != nil {
		return datas, err
	}

	for rows.Next() {
		var d LandingPageData
		d, err = dataScan(rows)
		if err != nil {
			return datas, err
		}
		datas = append(datas, d)
	}
	defer rows.Close()
	return datas, nil
}

func images(ctx *middleware.APIContext, LandingPageID int) (images []LandingPageImage, err error) {

	stmt, err := ctx.DB.Prepare(GetLandingPageImagesStmt)
	if err != nil {
		return images, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(LandingPageID)
	if err != nil {
		return images, err
	}

	for rows.Next() {
		var img LandingPageImage
		img, err = imageScan(rows)
		if err != nil {
			return images, err
		}
		images = append(images, img)
	}
	defer rows.Close()
	return images, nil
}

func imageScan(s database.Scanner) (LandingPageImage, error) {
	var img LandingPageImage
	var urlstr *string
	err := s.Scan(
		&img.Id,
		&img.LandingPageID,
		&urlstr,
		&img.Sort,
	)
	if err != nil {
		return img, err
	}
	if urlstr != nil {
		img.Url, _ = url.Parse(*urlstr)
	}
	return img, nil
}

func dataScan(s database.Scanner) (LandingPageData, error) {
	var d LandingPageData
	err := s.Scan(
		&d.Id,
		&d.LandingPageID,
		&d.DataKey,
		&d.DataValue,
	)
	if err != nil {
		return d, err
	}
	return d, nil
}
