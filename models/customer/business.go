package customer

import (
	"github.com/curt-labs/API/helpers/sortutil"
	"github.com/curt-labs/API/middleware"
)

var (
	getBusinessClassesStmt = `select b.BusinessClassID, b.name, b.sort, b.showOnWebsite from BusinessClass as b
		join ApiKeyToBrand as atb on atb.brandID = b.brandID
		join ApiKey as a on a.id = atb.keyID
		where a.api_key = ? && (atb.brandID = ? or 0 = ?) && b.showOnWebsite = 1
		group by b.name
		order by b.sort`
	createBusinessClass = `insert into BusinessClass (name, sort, showOnWebsite, brandID) values (?,?,?,?)`
	deleteBusinessClass = `delete from BusinessClass where BusinessClassID = ?`
)

type BusinessClasses []BusinessClass
type BusinessClass struct {
	ID            int    `json:"id" xml:"id"`
	Name          string `json:"name" xml:"name"`
	Sort          int    `json:"sort" xml:"sort"`
	ShowOnWebsite bool   `json:"show" xml:"show"`
	BrandID       int    `json:"brandID omitempty" xml:"brandID omitempty"`
}

func GetAllBusinessClasses(ctx *middleware.APIContext) (classes BusinessClasses, err error) {

	stmt, err := ctx.DB.Prepare(getBusinessClassesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	if err != nil {
		return
	}
	var bc BusinessClass
	for rows.Next() {
		bc = BusinessClass{}
		err = rows.Scan(
			&bc.ID,
			&bc.Name,
			&bc.Sort,
			&bc.ShowOnWebsite,
		)
		if err != nil {
			return
		}
		classes = append(classes, bc)
	}
	defer rows.Close()

	sortutil.AscByField(classes, "Sort")
	return
}

func (b *BusinessClass) Create(ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(createBusinessClass)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(b.Name, b.Sort, b.ShowOnWebsite, b.BrandID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = int(id)
	return err
}

func (b *BusinessClass) Delete(ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(deleteBusinessClass)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(b.ID)
	if err != nil {
		return err
	}
	return err
}
