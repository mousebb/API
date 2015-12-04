package customer

import (
	"github.com/curt-labs/API/middleware"
	_ "github.com/go-sql-driver/mysql"
)

var (
	createCustomerBrand     = `insert into CustomerToBrand (cust_id,brandID) values(?,?)`
	deleteCustomerBrand     = `delete from CustomerToBrand where cust_id = ? and brandID = ?`
	deleteAllCustomerBrands = `delete from CustomerToBrand where cust_id = ?`
)

func (c *Customer) CreateCustomerBrand(brandID int, ctx *middleware.APIContext) error {
	var err error

	stmt, err := ctx.DB.Prepare(createCustomerBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id, brandID)

	return err
}

func (c *Customer) DeleteCustomerBrand(brandID int, ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(deleteCustomerBrand)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Id, brandID)

	return err
}

func (c *Customer) DeleteAllCustomerBrands(ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(deleteAllCustomerBrands)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Id)

	return err
}
