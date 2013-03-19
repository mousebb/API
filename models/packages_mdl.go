package models

import (
	"../helpers/database"
)

var (
	partPackageStmt = `select pp.height as height, pp.length as length, pp.width as width, pp.weight as weight, pp.quantity as quantity,
				um_dim.code as dimensionUnit, um_dim.name as dimensionUnitLabel, um_wt.code as weightUnit, um_wt.name as weightUnitLabel,
				um_pkg.code as packageUnit, um_pkg.name as packageUnitLabel
				from PartPackage as pp
				join UnitOfMeasure as um_dim on pp.dimensionUOM = um_dim.ID
				join UnitOfMeasure as um_wt on pp.weightUOM = um_wt.ID
				join UnitOfMeasure as um_pkg on pp.packageUOM = um_pkg.ID
				where pp.partID = %d`
)		

type Package struct {
	Height, Width, Length, Quantity   float64
	Weight                            float64
	DimensionUnit, DimensionUnitLabel string
	WeightUnit, WeightUnitLabel       string
	PackageUnit, PackageUnitLabel     string
}

func (part *Part) GetPartPackaging() error {
	db := database.Db

	rows, res, err := db.Query(partPackageStmt, part.PartId)
	if database.MysqlError(err) {
		return err
	}

	height := res.Map("height")
	length := res.Map("length")
	width := res.Map("width")
	weight := res.Map("weight")
	qty := res.Map("quantity")
	dimUnit := res.Map("dimensionUnit")
	dimUnitLabel := res.Map("dimensionUnitLabel")
	weightUnit := res.Map("weightUnit")
	weightUnitLabel := res.Map("weightUnitLabel")
	pkgUnit := res.Map("packageUnit")
	pkgUnitLabel := res.Map("packageUnitLabel")

	var pkgs []Package
	for _, row := range rows {
		p := Package {
			Height: row.Float(height), 
			Width: row.Float(width), 
			Length: row.Float(length), 
			Quantity: row.Float(qty),
			Weight: row.Float(weight),
			DimensionUnit: row.Str(dimUnit), 
			DimensionUnitLabel: row.Str(dimUnitLabel),
			WeightUnit: row.Str(weightUnit), 
			WeightUnitLabel: row.Str(weightUnitLabel),
			PackageUnit: row.Str(pkgUnit), 
			PackageUnitLabel: row.Str(pkgUnitLabel),
		}
		pkgs = append(pkgs, p)
	}

	part.Packages = pkgs
	return nil
}