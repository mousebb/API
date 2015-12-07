package products

import (
	"strconv"
	"strings"

	"github.com/curt-labs/API/middleware"
)

func (l *Lookup) GetSubmodels(ctx *middleware.APIContext) error {
	stmtBeginning := `
		select distinct s.SubmodelName from vcdb_Vehicle as v
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && bv.YearID = ? && ma.MakeName = ? && mo.ModelName = ? `
	stmtEnd := ` order by s.SubmodelName`
	brandStmt := " && p.brandID in ("

	for _, b := range l.Brands {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt + stmtEnd

	stmt, err := ctx.DB.Prepare(wholeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make, l.Vehicle.Base.Model)
	if err != nil {
		return err
	}

	l.Submodels = make([]string, 0)
	for res.Next() {
		var m string
		err = res.Scan(&m)
		if err == nil {
			l.Submodels = append(l.Submodels, m)
		}
	}
	defer res.Close()

	l.Pagination = Pagination{
		TotalItems:    len(l.Submodels),
		ReturnedCount: len(l.Submodels),
		Page:          1,
		PerPage:       len(l.Submodels),
		TotalPages:    1,
	}

	return nil
}
