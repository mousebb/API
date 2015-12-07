package middleware

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/curt-labs/API/models/brand"
)

type DataContext struct {
	BrandID     int
	WebsiteID   int
	APIKey      string
	CustomerID  int
	UserID      string
	Brands      []brand.Brand
	BrandString string
	BrandArray  []int
}

var (
	BuildContext = `select b.ID, b.name, b.code, ak.api_key, cu.id, c.customerID from ApiKey ak
					join CustomerUser cu on ak.user_id = cu.id
					join Customer c on cu.cust_ID = c.cust_id
					join ApiKeyToBrand akb on ak.id = akb.keyID
					join Brand b on akb.brandID = b.id
					join ApiKeyType as akt on ak.type_id = akt.id
					where ak.api_key = ? && (UPPER(akt.type) = ? || 1 = ?)`
)

func (ctx *APIContext) BuildDataContext(k, t string) error {
	stmt, err := ctx.DB.Prepare(BuildContext)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var rows *sql.Rows
	if t != "" {
		rows, err = stmt.Query(k, t, 0)
	} else {
		rows, err = stmt.Query(k, t, 1)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("%s", "authentication failed")
		}
		return err
	}
	defer rows.Close()

	var brandIDStr []string
	for rows.Next() {
		var bID, customerID *int
		var name, code, key, userID *string
		err = rows.Scan(
			&bID,
			&name,
			&code,
			&key,
			&userID,
			&customerID,
		)
		if err != nil || bID == nil || name == nil || code == nil || key == nil || userID == nil || customerID == nil {
			return fmt.Errorf("%s", "authentication failed")
		}

		if len(ctx.DataContext.Brands) == 0 {
			ctx.DataContext.APIKey = *key
			ctx.DataContext.CustomerID = *customerID
			ctx.DataContext.UserID = *userID
		}
		ctx.DataContext.Brands = append(ctx.DataContext.Brands, brand.Brand{
			ID:   *bID,
			Name: *name,
			Code: *code,
		})
		brandIDStr = append(brandIDStr, strconv.Itoa(*bID))
	}

	ctx.DataContext.brandArray()
	ctx.DataContext.brandString()

	return nil
}

func (dtx *DataContext) brandString() {
	var ids []string
	for _, b := range dtx.Brands {
		ids = append(ids, string(b.ID))
	}
	dtx.BrandString = strings.Join(ids, ",")
}

func (dtx *DataContext) brandArray() {
	dtx.BrandArray = []int{}
	for _, b := range dtx.Brands {
		dtx.BrandArray = append(dtx.BrandArray, b.ID)
	}
}
