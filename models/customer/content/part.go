package custcontent

import (
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/api"
	"github.com/curt-labs/API/middleware"
	_ "github.com/go-sql-driver/mysql"
)

var (
	customerPartContent_Grouped = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
									ct.type,ct.allowHTML,ccb.partID
									from CustomerContent as cc
									join CustomerContentBridge as ccb on cc.id = ccb.contentID
									join ContentType as ct on cc.typeID = ct.cTypeID
									join Customer as c on cc.custID = c.cust_id
									join CustomerUser as cu on c.cust_id = cu.cust_ID
									join ApiKey as ak on cu.id = ak.user_id
									where api_key = ? and ccb.partID IN (?)
									group by cc.id, ccb.partID`

	allCustomerPartContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
								ct.type,ct.allowHTML,ccb.partID
								from CustomerContent as cc
								join CustomerContentBridge as ccb on cc.id = ccb.contentID
								join ContentType as ct on cc.typeID = ct.cTypeID
								join Customer as c on cc.custID = c.cust_id
								join CustomerUser as cu on c.cust_id = cu.cust_ID
								join ApiKey as ak on cu.id = ak.user_id
								where api_key = ? and ccb.partID > 0
								group by ccb.partID, cc.id
								order by ccb.partID`

	customerPartContent = `select cc.id, cc.text,cc.added,cc.modified,cc.deleted,
								ct.type,ct.allowHTML,ccb.partID
								from CustomerContent as cc
								join CustomerContentBridge as ccb on cc.id = ccb.contentID
								join ContentType as ct on cc.typeID = ct.cTypeID
								join Customer as c on cc.custID = c.cust_id
								join CustomerUser as cu on c.cust_id = cu.cust_ID
								join ApiKey as ak on cu.id = ak.user_id
								where api_key = ? and ccb.partID = ?
								group by cc.id`
)

type PartContent struct {
	PartId  int
	Content []CustomerContent
}

// Retrieves all part content for this customer
func GetAllPartContent(ctx *middleware.APIContext) (content []PartContent, err error) {

	stmt, err := ctx.DB.Prepare(allCustomerPartContent)
	if err != nil {
		return
	}

	res, err := stmt.Query(ctx.DataContext.APIKey)
	if err != nil {
		return
	}
	defer res.Close()

	rawContent := make(map[int][]CustomerContent, 0)
	var partId int
	var deleted *bool
	var added, modified *time.Time
	var ctype string

	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&ctype,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.ContentType.Type = "Part:" + ctype
		if added != nil {
			cc.Added = *added
		}
		if modified != nil {
			cc.Modified = *modified
		}
		if deleted != nil {
			cc.Hidden = *deleted
		}
		part_id := partId
		if part_id > 0 {
			rawContent[part_id] = append(rawContent[part_id], cc)
		}

	}

	for k, _ := range rawContent {
		pCon := PartContent{
			PartId:  k,
			Content: rawContent[k],
		}
		content = append(content, pCon)
	}

	return
}

// Retrieves specific part content for this customer
func GetPartContent(partID int, ctx *middleware.APIContext) (content []CustomerContent, err error) {
	content = make([]CustomerContent, 0) // initializer

	stmt, err := ctx.DB.Prepare(customerPartContent)
	if err != nil {
		return content, err
	}
	defer stmt.Close()

	res, err := stmt.Query(ctx.DataContext.APIKey, partID)
	if err != nil {
		return
	}
	defer res.Close()

	var partId int
	var deleted *bool
	var added, modified *time.Time
	var ctype string

	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&ctype,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.ContentType.Type = "Part:" + ctype
		if added != nil {
			cc.Added = *added
		}
		if modified != nil {
			cc.Modified = *modified
		}
		if deleted != nil {
			cc.Hidden = *deleted
		}
		content = append(content, cc)
	}
	defer res.Close()
	return
}

func GetGroupedPartContent(ids []string, ctx *middleware.APIContext) (content map[int][]CustomerContent, err error) {
	content = make(map[int][]CustomerContent, len(ids))

	for i := 0; i < len(ids); i++ {
		intId, err := strconv.Atoi(ids[i])
		if err == nil {
			content[intId] = make([]CustomerContent, 0)
		}
	}
	escaped_key := api_helpers.Escape(ctx.DataContext.APIKey)

	stmt, err := ctx.DB.Prepare(customerPartContent_Grouped)
	if err != nil {
		return
	}
	defer stmt.Close()

	var partId int
	var deleted *bool
	var added, modified *time.Time
	var ctype string

	res, err := stmt.Query(escaped_key, strings.Join(ids, ","))
	if err != nil {
		return
	}
	defer res.Close()

	for res.Next() {
		var cc CustomerContent
		err = res.Scan(
			&cc.Id,
			&cc.Text,
			&added,
			&modified,
			&deleted, //Not Used
			&ctype,
			&cc.ContentType.AllowHtml,
			&partId,
		)
		if err != nil {
			return content, err
		}
		cc.ContentType.Type = "Part:" + ctype
		if added != nil {
			cc.Added = *added
		}
		if modified != nil {
			cc.Modified = *modified
		}
		if deleted != nil {
			cc.Hidden = *deleted
		}
		content[partId] = append(content[partId], cc)
	}

	return
}
