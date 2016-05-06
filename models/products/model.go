package products

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/curt-labs/API/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func (l *Lookup) GetModels(ctx *middleware.APIContext) error {
	stmtBeginning := `select distinct mo.ModelName from vcdb_Model as mo
		join BaseVehicle as bv on mo.ID = bv.ModelID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && bv.YearID = ? && ma.MakeName = ? `

	stmtEnd := ` order by mo.ModelName`
	brandStmt := " && p.brandID in ("

	for _, b := range l.Brands {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt + stmtEnd

	if ctx == nil || ctx.DB == nil {
		return fmt.Errorf("%s", "failed to establish database connection")
	}

	stmt, err := ctx.DB.Prepare(wholeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make)
	if err != nil {
		return err
	}

	l.Models = make([]string, 0)
	for res.Next() {
		var m string
		err = res.Scan(&m)
		if err == nil {
			l.Models = append(l.Models, m)
		}
	}
	defer res.Close()

	l.Pagination = Pagination{
		TotalItems:    len(l.Models),
		ReturnedCount: len(l.Models),
		Page:          1,
		PerPage:       len(l.Models),
		TotalPages:    1,
	}

	return nil
}
