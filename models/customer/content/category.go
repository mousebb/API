package custcontent

import (
	"time"

	"github.com/curt-labs/API/middleware"
)

type CategoryContent struct {
	CategoryId int
	Content    []CustomerContent
}

var (
	allCustomerCategoryContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
									ct.type,ct.allowHTML,ccb.catID
									from CustomerContent as cc
									join CustomerContentBridge as ccb on cc.id = ccb.contentID
									join ContentType as ct on cc.typeID = ct.cTypeID
									join Customer as c on cc.custID = c.cust_id
									join CustomerUser as cu on c.cust_id = cu.cust_ID
									join ApiKey as ak on cu.id = ak.user_id
									where api_key = ? and ccb.catID > 0
									group by ccb.catID, cc.id
									order by ccb.catID`
	customerCategoryContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
									ct.cTypeID, ct.type,ct.allowHTML,ccb.catID
									from CustomerContent as cc
									join CustomerContentBridge as ccb on cc.id = ccb.contentID
									join ContentType as ct on cc.typeID = ct.cTypeID
									join Customer as c on cc.custID = c.cust_id
									join CustomerUser as cu on c.cust_id = cu.cust_ID
									join ApiKey as ak on cu.id = ak.user_id
									where api_key = ? and ccb.catID = ?
									group by cc.id`
)

// Retrieves all category content for this customer
func GetAllCategoryContent(ctx *middleware.APIContext) (content []CategoryContent, err error) {

	stmt, err := ctx.DB.Prepare(allCustomerCategoryContent)
	if err != nil {
		return content, err
	}

	res, err := stmt.Query(ctx.DataContext.APIKey)
	if err != nil {
		return
	}
	defer res.Close()

	var catID int
	var added, modified *time.Time
	var deleted *bool
	rawContent := make(map[int][]CustomerContent, 0)
	for res.Next() {
		var c CustomerContent
		err = res.Scan(
			&c.Id,
			&c.Text,
			&added,
			&modified,
			&deleted,
			&c.ContentType.Type,
			&c.ContentType.AllowHtml,
			&catID,
		)
		if err != nil {
			return content, err
		}

		if catID > 0 {
			rawContent[catID] = append(rawContent[catID], c)
		}
		if added != nil {
			c.Added = *added
		}
		if modified != nil {
			c.Modified = *modified
		}
		if deleted != nil {
			c.Hidden = *deleted
		}
	}

	for k, _ := range rawContent {
		catCon := CategoryContent{
			CategoryId: k,
			Content:    rawContent[k],
		}
		content = append(content, catCon)
	}

	return
}

// Retrieves specific category content for this customer
func GetCategoryContent(catID int, ctx *middleware.APIContext) (content []CustomerContent, err error) {
	content = make([]CustomerContent, 0) // initializer

	stmt, err := ctx.DB.Prepare(customerCategoryContent)
	if err != nil {
		return
	}
	res, err := stmt.Query(ctx.DataContext.APIKey, catID)
	if err != nil {
		return
	}
	defer res.Close()

	var deleted *bool
	var added, modified *time.Time

	for res.Next() {
		var c CustomerContent
		err = res.Scan(
			&c.Id,
			&c.Text,
			&added,
			&modified,
			&deleted,
			&c.ContentType.Id,
			&c.ContentType.Type,
			&c.ContentType.AllowHtml,
			&catID,
		)
		if err != nil {
			return content, err
		}
		if added != nil {
			c.Added = *added
		}
		if modified != nil {
			c.Modified = *modified
		}
		if deleted != nil {
			c.Hidden = *deleted
		}
		content = append(content, c)
	}

	return
}
