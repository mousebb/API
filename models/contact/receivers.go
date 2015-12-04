package contact

import (
	"errors"

	"github.com/curt-labs/API/helpers/email"
	"github.com/curt-labs/API/middleware"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllContactReceiversStmt              = `select contactReceiverID, first_name, last_name, email from ContactReceiver`
	getContactReceiverStmt                  = `select contactReceiverID, first_name, last_name, email from ContactReceiver where contactReceiverID = ?`
	addContactReceiverStmt                  = `insert into ContactReceiver(first_name, last_name, email) values (?,?,?)`
	updateContactReceiverStmt               = `update ContactReceiver set first_name = ?, last_name = ?, email = ? where contactReceiverID = ?`
	deleteContactReceiverStmt               = `delete from ContactReceiver where contactReceiverID = ?`
	createReceiverContactTypeJoin           = `insert into ContactReceiver_ContactType (ContactReceiverID, ContactTypeID) values (?,?)`
	deleteReceiverContactTypeJoin           = `delete from ContactReceiver_ContactType where ContactReceiverID = ? and  ContactTypeID = ?`
	deleteReceiverContactTypeJoinByReceiver = `delete from ContactReceiver_ContactType where ContactReceiverID = ?`
	getContactTypesByReceiver               = `select crct.contactTypeID, ct.name, ct.showOnWebsite, ct.brandID from ContactReceiver_ContactType as crct
												left join ContactType as ct on crct.ContactTypeID = ct.contactTypeID where crct.contactReceiverID = ?`
)

type ContactReceivers []ContactReceiver
type ContactReceiver struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	ContactTypes ContactTypes
}

func GetAllContactReceivers(ctx *middleware.APIContext) (receivers ContactReceivers, err error) {

	stmt, err := ctx.DB.Prepare(getAllContactReceiversStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var cr ContactReceiver
		err = rows.Scan(
			&cr.ID,
			&cr.FirstName,
			&cr.LastName,
			&cr.Email,
		)
		if err != nil {
			return
		}
		err = cr.GetContactTypes(ctx)
		if err != nil {
			return
		}
		receivers = append(receivers, cr)
	}
	defer rows.Close()

	return
}

func (cr *ContactReceiver) Get(ctx *middleware.APIContext) error {
	if cr.ID == 0 {
		return errors.New("invalid receiver identifier")
	}

	stmt, err := ctx.DB.Prepare(getContactReceiverStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(cr.ID).Scan(
		&cr.ID,
		&cr.FirstName,
		&cr.LastName,
		&cr.Email,
	)

	return cr.GetContactTypes(ctx)
}

func (cr *ContactReceiver) Add(ctx *middleware.APIContext) error {
	if !email.IsEmail(cr.Email) {
		return errors.New("Empty or invalid email address.")
	}

	stmt, err := ctx.DB.Prepare(addContactReceiverStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(cr.FirstName, cr.LastName, cr.Email)
	if err != nil {
		return err
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		cr.ID = int(id)
	}
	//add contact types
	if len(cr.ContactTypes) > 0 {
		for _, ct := range cr.ContactTypes {
			err = cr.CreateTypeJoin(ct, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cr *ContactReceiver) Update(ctx *middleware.APIContext) error {
	if cr.ID == 0 {
		return errors.New("Invalid ContactReceiver ID")
	}
	if !email.IsEmail(cr.Email) {
		return errors.New("Empty or invalid email address.")
	}

	stmt, err := ctx.DB.Prepare(updateContactReceiverStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cr.FirstName, cr.LastName, cr.Email, cr.ID)

	//update type joins
	if len(cr.ContactTypes) == 0 {
		return err
	}

	err = cr.DeleteTypeJoinByReceiver(ctx)
	if err != nil {
		return err
	}
	for _, ct := range cr.ContactTypes {
		err = cr.CreateTypeJoin(ct, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cr *ContactReceiver) Delete(ctx *middleware.APIContext) error {
	if cr.ID == 0 {
		return errors.New("invalid reciever identifier")
	}

	stmt, err := ctx.DB.Prepare(deleteContactReceiverStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cr.ID)
	if err != nil {
		return err
	}

	//delete receiver-type join
	return cr.DeleteTypeJoinByReceiver(ctx)
}

//get a contact receiver's contact types
func (cr *ContactReceiver) GetContactTypes(ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(getContactTypesByReceiver)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var ct ContactType
	res, err := stmt.Query(cr.ID)
	if err != nil {
		return err
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&ct.ID, &ct.Name, &ct.ShowOnWebsite, &ct.BrandID)
		if err != nil {
			return err
		}
		cr.ContactTypes = append(cr.ContactTypes, ct)
	}

	return nil
}

func (cr *ContactReceiver) CreateTypeJoin(ct ContactType, ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(createReceiverContactTypeJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(cr.ID, ct.ID)

	return err
}

//delete joins for a receiver-type pair
func (cr *ContactReceiver) DeleteTypeJoin(ct ContactType, ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(deleteReceiverContactTypeJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cr.ID, ct.ID)

	return err
}

//delete all type-receiver joins for a receiver
func (cr *ContactReceiver) DeleteTypeJoinByReceiver(ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(deleteReceiverContactTypeJoinByReceiver)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cr.ID)

	return err
}
