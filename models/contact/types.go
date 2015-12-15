package contact

import (
	"errors"
	"github.com/curt-labs/API/middleware"
)

var (
	getAllContactTypesStmt = `select ct.contactTypeID, ct.name, ct.showOnWebsite, ct.brandID from ContactType as ct
		join ApiKeyToBrand as akb on akb.brandID = ct.brandID
		join ApiKey as ak on ak.id = akb.keyID
		where ak.api_key = ? && (ct.brandID = ? or 0 = ?)`
	getContactTypeStmt    = `select contactTypeID, name, showOnWebsite from ContactType where contactTypeID = ?`
	getReceiverByType     = `select cr.contactReceiverID, cr.first_name, cr.last_name, cr.email from ContactReceiver_ContactType as crct
								left join ContactReceiver as cr on crct.contactReceiverID = cr.contactReceiverID
								where crct.contactTypeID = ?`
	getTypeNameFromId = `select name from ContactType where contactTypeID = ?`
)

type ContactType struct {
	ID            int    `json:"id" xml:"id"`
	Name          string `json:"name" xml: "name"`
	ShowOnWebsite bool   `json:"show" xml:"show"`
	BrandID       int    `json:"brandId" xml:"brandId"`
}

func GetAllContactTypes(ctx *middleware.APIContext) (types []ContactType, err error) {

	stmt, err := ctx.DB.Prepare(getAllContactTypesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ct ContactType
		err = rows.Scan(
			&ct.ID,
			&ct.Name,
			&ct.ShowOnWebsite,
			&ct.BrandID,
		)
		if err != nil {
			return
		}
		types = append(types, ct)
	}
	return
}

func (ct *ContactType) Get(ctx *middleware.APIContext) error {
	if ct.ID == 0 {
		return errors.New("invalid ContactType identifier")
	}

		stmt, err := ctx.DB.Prepare(getContactTypeStmt)
		if err != nil {
			return err
		}
		defer stmt.Close()

		return stmt.QueryRow(ct.ID).Scan(
			&ct.ID,
			&ct.Name,
			&ct.ShowOnWebsite,
		)
}

func GetContactTypeNameFromId(ctx *middleware.APIContext, id int) (string, error) {
	var err error
	var name string

	stmt, err := ctx.DB.Prepare(getTypeNameFromId)
	if err != nil {
		return name, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&name)

	return name, err
}

func (ct *ContactType) GetReceivers(ctx *middleware.APIContext) (crs []ContactReceiver, err error) {

	stmt, err := ctx.DB.Prepare(getReceiverByType)
	if err != nil {
		return crs, err
	}
	defer stmt.Close()

	var cr ContactReceiver
	res, err := stmt.Query(ct.ID)
	if err != nil {
		return crs, err
	}

	for res.Next() {
		err = res.Scan(
			&cr.ID,
			&cr.FirstName,
			&cr.LastName,
			&cr.Email,
		)
		if err != nil {
			return crs, err
		}
		crs = append(crs, cr)
	}
	defer res.Close()

	return crs, err
}
