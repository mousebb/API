package vehicle

import (
	"fmt"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"

	// "log"
	"net/http"
	"strings"
)

type ErrorResp struct {
	ConversionErrs []error `json:"conversion_errors" xml:"conversion_errors"`
	InsertErrs     []error `json:"insert_errors" xml:"insert_errors"`
}

//requires the "Consolidated App Guides" that MJ produces in Excel
//intended to be a short term solution until Aries-Curt data merge is complete
//powers the Godzilla application

//Import a Csv
//Fields are expected to be: Part (oldpartnumber), Make, Model, Style, Year - 5 columns total
func ImportCsv(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	contentTypeHeader := r.Header.Get("Content-Type")
	contentTypeArr := strings.Split(contentTypeHeader, ";")
	if len(contentTypeArr) < 1 {
		return nil, fmt.Errorf("%s", "Content-Type is not multipart/form-data")
	}
	contentType := contentTypeArr[0]
	if contentType != "multipart/form-data" {
		return nil, fmt.Errorf("%s", "Content-Type is not multipart/form-data")
	}
	file, header, err := r.FormFile("file")

	if err != nil {
		return nil, fmt.Errorf("%s", "Error getting file")
	}
	defer file.Close()

	collectionName := header.Filename

	conversionErrs, insertErrs, err := products.Import(file, collectionName)
	if err != nil {
		return nil, err
	}

	errResp := ErrorResp{
		ConversionErrs: conversionErrs,
		InsertErrs:     insertErrs,
	}

	return errResp, nil
}
