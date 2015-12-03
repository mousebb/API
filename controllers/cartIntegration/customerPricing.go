package pricingCtlr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/cartIntegration"
	"github.com/curt-labs/API/models/customer"
)

func setDB(r *http.Request) error {
	brandID, err := strconv.Atoi(r.URL.Query().Get("brandID"))
	if err != nil {
		return err
	}
	cartIntegration.Brand_ID = brandID
	return nil
}

func setCustomerID(r *http.Request) error {
	c := customer.Customer{}
	err := c.GetCustomerIdFromKey(r.URL.Query().Get("key"))
	if err != nil {
		return err
	}
	cartIntegration.Customer_ID = c.Id
	return nil
}

// GetPricing Requires APIKEY and brandID in header
func GetPricing(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}
	err = setCustomerID(r)
	if err != nil {
		return nil, err
	}

	return cartIntegration.GetCustomerPrices()
}

// GetPricingPaged Requires APIKEY and brandID in header
// Requires count and page in params
func GetPricingPaged(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}
	err = setCustomerID(r)
	if err != nil {
		return nil, err
	}

	page, err := strconv.Atoi(ctx.Params.ByName("page"))
	if page < 1 || err != nil {
		return nil, err
	}

	count, err := strconv.Atoi(ctx.Params.ByName("count"))
	if count < 1 || err != nil {
		return nil, err
	}

	return cartIntegration.GetPricingPaged(page, count)
}

//GetPricingCount Returns int
func GetPricingCount(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}
	err = setCustomerID(r)
	if err != nil {
		return nil, err
	}

	return cartIntegration.GetPricingCount()
}

// GetPartPricesByPartID Returns Mfr Prices for a part
func GetPartPricesByPartID(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}
	partID, err := strconv.Atoi(ctx.Params.ByName("part"))
	if partID < 1 || err != nil {
		return nil, err
	}

	return cartIntegration.GetPartPricesByPartID(partID)
}

// GetAllPartPrices Returns Mfr Prices
func GetAllPartPrices(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	err := setDB(r)
	if err != nil {
		return nil, err
	}

	return cartIntegration.GetPartPrices()
}

// CreatePrice ...
func CreatePrice(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var price cartIntegration.CustomerPrice
	err = json.Unmarshal(body, &price)
	if err != nil {
		return nil, err
	}
	err = setCustomerID(r)
	if err != nil {
		return nil, err
	}
	price.CustID = cartIntegration.Customer_ID
	err = validatePrice(price)
	if err != nil {
		return nil, err
	}
	err = price.Create()
	if err != nil {
		return nil, err
	}
	err = price.InsertCartIntegration()

	return price, err
}

// UpdatePrice ...
func UpdatePrice(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}

	var price cartIntegration.CustomerPrice
	err = json.NewDecoder(r.Body).Decode(&price)
	if err != nil {
		return nil, err
	}

	err = setCustomerID(r)
	if err != nil {
		return nil, err
	}

	price.CustID = cartIntegration.Customer_ID
	err = validatePrice(price)
	if err != nil {
		return nil, err
	}
	err = price.Update()
	if err != nil {
		return nil, err
	}

	err = price.UpdateCartIntegration()

	return price, err
}

// ResetAllToMap Set all of a customer's prices to MAP
func ResetAllToMap(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}
	err = setCustomerID(r)
	if err != nil {
		return nil, err
	}
	custPrices, err := cartIntegration.GetCustomerPrices()
	if err != nil {
		return nil, err
	}

	//create map of MAP prices
	prices, err := cartIntegration.GetMAPPartPrices()
	if err != nil {
		return nil, err
	}
	priceMap := make(map[int]cartIntegration.Price)
	for _, p := range prices {
		priceMap[p.PartID] = p
	}

	//set to MAP
	for i := range custPrices {
		custPrices[i].Price = priceMap[custPrices[i].PartID].Price
		if custPrices[i].CustID == 0 {
			custPrices[i].CustID = cartIntegration.Customer_ID
		}
		if custPrices[i].ID == 0 {
			err = custPrices[i].Create()
		} else {
			err = custPrices[i].Update()
		}
		if err != nil {
			return nil, err
		}
	}

	return custPrices, nil
}

// Global Sets all of a customer's prices to a percentage of the price type specified in params
func Global(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}
	err = setCustomerID(r)
	if err != nil {
		return nil, err
	}
	priceType := ctx.Params.ByName("type")
	percent, err := strconv.ParseFloat(ctx.Params.ByName("percentage"), 64)
	if err != nil {
		return nil, err
	}
	percent = percent / 100

	//create partPriceMap
	prices, err := cartIntegration.GetPartPrices()
	if err != nil {
		return nil, err
	}

	priceMap := make(map[string]float64)
	for _, p := range prices {
		key := strconv.Itoa(p.PartID) + p.Type
		priceMap[key] = p.Price
	}

	//get CustPrices
	custPrices, err := cartIntegration.GetCustomerPrices()
	if err != nil {
		return nil, err
	}

	//set to percentage
	for i := range custPrices {
		if custPrices[i].CustID == 0 {
			custPrices[i].CustID = cartIntegration.Customer_ID
		}
		custPrices[i].Price = priceMap[strconv.Itoa(custPrices[i].PartID)+priceType] * percent
		if custPrices[i].ID == 0 {
			err = custPrices[i].Create()
		} else {
			err = custPrices[i].Update()

		}
		if err != nil {
			return nil, err
		}
	}

	return custPrices, nil
}

// GetAllPriceTypes Get those price types
func GetAllPriceTypes(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = setDB(r)
	if err != nil {
		return nil, err
	}

	return cartIntegration.GetAllPriceTypes()
}

//Utility
func validatePrice(p cartIntegration.CustomerPrice) error {
	if p.CustID < 1 {
		return errors.New("Customer ID cannot be less than 1")
	}
	if p.PartID < 1 {
		return errors.New("Part ID cannot be less than 1")
	}
	if p.IsSale == 1 {
		if p.SaleStart.Before(time.Now()) {
			return errors.New("The starting date is required and cannot be set to a date prior to today.")
		}

		if p.SaleStart.After(*p.SaleEnd) {
			return errors.New("The sale starting date cannot be set to a date after the sale ending date.")
		}

		if p.SaleEnd.Before(time.Now()) || p.SaleEnd.Before(*p.SaleStart) {
			return errors.New("The ending date is required and cannot be set to a date prior to today or the sale starting date.")
		}
	}
	return nil
}
