package cartIntegration

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

type CartIntegration struct {
	ID         int `json:"id,omitempty" xml:"id,omitempty"`
	PartID     int `json:"partId,omitempty" xml:"partId,omitempty"`
	CustPartID int `json:"custPartId,omitempty" xml:"custPartId,omitempty"`
	CustID     int `json:"custId,omitempty" xml:"custId,omitempty"`
}

var (
	getCIsByPartID = `select ci.referenceID, ci.partID, ci.custPartID, ci.custID from CartIntegration as ci
						Join CustomerToBrand as cub on cub.cust_id = ci.custID
						Join ApiKeyToBrand as akb on akb.brandID = cub.brandID
						Join ApiKey as ak on akb.keyID = ak.id
						where ci.partID = ? && (ak.api_key = ? && (cub.brandID = ? OR 0=?))`
	getCIsByCustID = `select ci.referenceID, ci.partID, ci.custPartID, ci.custID from CartIntegration as ci
						Join CustomerToBrand as cub on cub.cust_id = ci.custID
						Join ApiKeyToBrand as akb on akb.brandID = cub.brandID
						Join ApiKey as ak on akb.keyID = ak.id
						where ci.custID = ? && (ak.api_key = ? && (cub.brandID = ? OR 0=?))`
	getCI = `select ci.referenceID, ci.partID, ci.custPartID, ci.custID from CartIntegration as ci
						Join CustomerToBrand as cub on cub.cust_id = ci.custID
						Join ApiKeyToBrand as akb on akb.brandID = cub.brandID
						Join ApiKey as ak on akb.keyID = ak.id
						where ci.custID = ? && ci.partID = ? && (ak.api_key = ? && (cub.brandID = ? OR 0=?))`
	insertCI = `insert into CartIntegration (partID, custPartID, custID) values (?,?,?)`
	updateCI = `update CartIntegration set partID = ?, custPartID = ?, custID = ? where referenceID = ?`
	deleteCI = `delete from CartIntegration where referenceID = ?`
)

func GetAllCartIntegrations(dtx *apicontext.DataContext) (cis []CartIntegration, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cis, err
	}
	defer db.Close()
	// reason for this is because a join on api key table was wayyyyyy too slow (too many records).
	stmtBeginning := `select ci.referenceID, ci.partID, ci.custPartID, ci.custID from CartIntegration as ci
						Join CustomerToBrand as cub on cub.cust_id = ci.custID `
	brandStmt := "where cub.brandID in ("
	for _, b := range dtx.BrandArray {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt
	stmt, err := db.Prepare(wholeStmt)

	if err != nil {
		return cis, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	var c CartIntegration
	for res.Next() {
		err = res.Scan(&c.ID, &c.PartID, &c.CustPartID, &c.CustID)
		if err != nil {
			return cis, err
		}
		cis = append(cis, c)
	}
	defer res.Close()
	return cis, err
}

func GetCartIntegrationsByPart(ci CartIntegration, dtx *apicontext.DataContext) (cis []CartIntegration, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cis, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCIsByPartID)
	if err != nil {
		return cis, err
	}
	defer stmt.Close()
	res, err := stmt.Query(ci.PartID, dtx.APIKey, dtx.BrandID, dtx.BrandID)
	var c CartIntegration
	for res.Next() {
		err = res.Scan(&c.ID, &c.PartID, &c.CustPartID, &c.CustID)
		if err != nil {
			return cis, err
		}
		cis = append(cis, c)
	}
	defer res.Close()
	return cis, err
}

func GetCartIntegrationsByCustomer(ci CartIntegration, dtx *apicontext.DataContext) (cis []CartIntegration, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cis, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCIsByCustID)
	if err != nil {
		return cis, err
	}
	defer stmt.Close()
	res, err := stmt.Query(ci.CustID, dtx.APIKey, dtx.BrandID, dtx.BrandID)
	if err != nil {
		return cis, err
	}

	var c CartIntegration
	for res.Next() {
		err = res.Scan(&c.ID, &c.PartID, &c.CustPartID, &c.CustID)
		if err != nil {
			return cis, err
		}
		cis = append(cis, c)
	}
	defer res.Close()
	return cis, err
}

func (c *CartIntegration) Get(dtx *apicontext.DataContext) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCI)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.CustID, c.PartID, dtx.APIKey, dtx.BrandID, dtx.BrandID).Scan(&c.ID, &c.PartID, &c.CustPartID, &c.CustID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (c *CartIntegration) Create(dtx *apicontext.DataContext) (err error) {
	if err := c.Get(dtx); err == nil && c.ID > 0 {
		return c.Update()
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertCI)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(c.PartID, c.CustPartID, c.CustID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = int(id)
	return nil
}

func (c *CartIntegration) Update() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(updateCI)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.PartID, c.CustPartID, c.CustID, c.ID)
	if err != nil {
		return err
	}
	return err
}

func (c *CartIntegration) Delete() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteCI)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID)
	if err != nil {
		return err
	}
	return err
}
