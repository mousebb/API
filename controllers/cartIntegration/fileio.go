package pricingCtlr

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/cartIntegration"
)

//TODO - extremely untested

// Upload ...
func Upload(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	if fileHeader != nil {
		contentType := fileHeader.Header.Get("Content-Type")

		if contentType != "text/comma-separated-values" && contentType != "text/csv" &&
			contentType != "application/csv" && contentType != "application/excel" &&
			contentType != "application/vnd.ms-excel" && contentType != "application/vnd.msexcel" {
			err = errors.New("The file you tried uploading was not a valid CSV file. Please try again using a valid CSV file.")
			return nil, err
		}
	}

	err = cartIntegration.UploadFile(file, ctx)
	if err != nil {
		return nil, err
	}

	return true, nil
}

// Download ...
func Download(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) {

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)

	customerPrices, err := cartIntegration.GetCustomerPrices(ctx)
	if err != nil {
		apierror.GenerateError(err.Error(), err, rw, r, http.StatusInternalServerError)
		return
	}

	//Price map
	prices, err := cartIntegration.GetPartPrices(ctx)
	if err != nil {
		apierror.GenerateError(err.Error(), err, rw, r, http.StatusInternalServerError)
		return
	}
	priceMap := make(map[string]float64)
	for _, p := range prices {
		priceMap[strconv.Itoa(p.PartID)+":"+p.Type] = p.Price
	}

	//Write
	wr.Write([]string{
		"CURT Part ID",
		"Customer Part ID",
		"Sale Price",
		"Sale Start Date",
		"Sale End Date",
		"Map Price",
		"List Price"})

	for _, price := range customerPrices {
		mapPrice := ""
		listPrice := ""

		mapPrice = strconv.FormatFloat(priceMap[strconv.Itoa(price.PartID)+":Map"], 'f', 2, 64)
		listPrice = strconv.FormatFloat(priceMap[strconv.Itoa(price.PartID)+":List"], 'f', 2, 64)

		//stringify dates
		var start, end string
		if price.SaleStart != nil && !price.SaleStart.IsZero() {
			start = price.SaleStart.Format(cartIntegration.DATE_FORMAT)
		}
		if price.SaleEnd != nil && !price.SaleStart.IsZero() {
			end = price.SaleEnd.Format(cartIntegration.DATE_FORMAT)
		}
		// log.Print(start, end)
		wr.Write([]string{
			strconv.Itoa(price.PartID),
			strconv.Itoa(price.CustomerPartID), //TODO - get CartIntegration at the same time
			strconv.FormatFloat(price.Price, 'f', 2, 64),
			start,
			end,
			mapPrice,
			listPrice,
		})

	}

	wr.Flush()
	rw.Header().Set("Content-Type", "text/csv")
	rw.Header().Set("Content-Disposition", "attachment;filename=data.csv")
	rw.Write(b.Bytes())

	return
}
