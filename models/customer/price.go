package customer

import (
	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/middleware"

	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Price struct {
	ID        int
	CustID    int
	PartID    int
	Price     float64
	IsSale    int
	SaleStart time.Time
	SaleEnd   time.Time
}

type Prices []Price

type CustomerPrices struct {
	Customer Customer `json:"customer" xml:"customer"`
	Prices   Prices   `json:"prices" xml:"prices"`
}

var (
	getPrice             = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing WHERE cust_price_id = ?"
	getPrices            = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing"
	createPrice          = "INSERT INTO CustomerPricing (cust_id, partID, price, isSale, sale_start, sale_end) VALUES (?,?,?,?,?,?)"
	updatePrice          = "UPDATE CustomerPricing SET cust_id = ?, partID = ?, price = ?, isSale = ?, sale_start = ?, sale_end = ? WHERE cust_price_id = ?"
	deletePrice          = "DELETE FROM CustomerPricing WHERE cust_price_id = ?"
	getPricesByCustomer  = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing WHERE cust_id = (select cust_id from Customer where customerID = ?)"
	getPricesByPart      = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing WHERE partID = ?"
	getPricesBySaleRange = "SELECT cust_price_id, cust_id, partID, price, isSale, sale_start, sale_end FROM CustomerPricing WHERE sale_start >= ? AND sale_end <= ? AND cust_id = (select cust_id from Customer where customerID = ?)"
)

const (
	timeFormat        = "2006-01-02"
	allPricesRedisKey = "prices"
)

func (p *Price) Get(ctx *middleware.APIContext) error {
	redis_key := "price:" + strconv.Itoa(p.ID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &p)
		return err
	}

	stmt, err := ctx.DB.Prepare(getPrice)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(p.ID).Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)
	if err != nil {
		return err
	}

	go redis.Setex(redis_key, p, 86400)
	return nil
}

func GetAllPrices(ctx *middleware.APIContext) (Prices, error) {
	var ps Prices
	var err error
	data, err := redis.Get(allPricesRedisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ps)
		return ps, err
	}

	stmt, err := ctx.DB.Prepare(getPrices)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query()
	for res.Next() {
		var p Price
		err = res.Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)
		ps = append(ps, p)
	}
	defer res.Close()
	go redis.Setex(allPricesRedisKey, ps, 86400)
	return ps, nil
}

func (p *Price) Create(ctx *middleware.APIContext) error {

	tx, err := ctx.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(createPrice)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(p.CustID, p.PartID, p.Price, p.IsSale, p.SaleStart, p.SaleEnd)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := res.LastInsertId()
	p.ID = int(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	go redis.Delete(allPricesRedisKey)
	go redis.Setex("price:"+strconv.Itoa(p.ID), p, 86400)

	return nil
}
func (p *Price) Update(ctx *middleware.APIContext) error {

	tx, err := ctx.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(updatePrice)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.CustID, p.PartID, p.Price, p.IsSale, p.SaleStart, p.SaleEnd, p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	go redis.Setex("price:"+strconv.Itoa(p.ID), p, 86400)
	go redis.Delete(fmt.Sprintf("prices:part:%d", strconv.Itoa(p.PartID)))
	go redis.Delete(fmt.Sprintf("customers:prices:%d", strconv.Itoa(p.CustID)))
	return nil
}

func (p *Price) Delete(ctx *middleware.APIContext) error {

	tx, err := ctx.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(deletePrice)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	go redis.Delete("price:" + strconv.Itoa(p.ID))
	go redis.Delete(fmt.Sprintf("prices:part:%d", strconv.Itoa(p.PartID)))
	go redis.Delete(fmt.Sprintf("customers:prices:%d", strconv.Itoa(p.CustID)))

	return nil
}

func (c *Customer) GetPricesByCustomer(ctx *middleware.APIContext) (CustomerPrices, error) {
	var cps CustomerPrices
	redis_key := "customers:prices:" + strconv.Itoa(c.Id)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &cps)
		return cps, err
	}

	stmt, err := ctx.DB.Prepare(getPricesByCustomer)
	if err != nil {
		return cps, err
	}
	defer stmt.Close()

	res, err := stmt.Query(c.Id)
	for res.Next() {
		var p Price
		res.Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)

		cps.Prices = append(cps.Prices, p)
	}
	go redis.Setex(redis_key, cps, 86400)
	return cps, err
}

func GetPricesByPart(partID int, ctx *middleware.APIContext) (Prices, error) {
	var ps Prices
	redis_key := "prices:part:" + strconv.Itoa(partID)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ps)
		return ps, err
	}

	stmt, err := ctx.DB.Prepare(getPricesByPart)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()

	res, err := stmt.Query(partID)
	if err != nil {
		return ps, err
	}
	defer res.Close()

	for res.Next() {
		var p Price
		res.Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)
		ps = append(ps, p)
	}

	go redis.Setex(redis_key, ps, 86400)

	return ps, nil
}

func (c *Customer) GetPricesBySaleRange(startDate, endDate time.Time, ctx *middleware.APIContext) (Prices, error) {
	var err error
	var ps Prices

	stmt, err := ctx.DB.Prepare(getPricesBySaleRange)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()

	res, err := stmt.Query(startDate, endDate, c.Id)
	if err != nil {
		return ps, err
	}
	defer res.Close()

	for res.Next() {
		var p Price
		res.Scan(&p.ID, &p.CustID, &p.PartID, &p.Price, &p.IsSale, &p.SaleStart, &p.SaleEnd)
		ps = append(ps, p)
	}

	return ps, nil
}
