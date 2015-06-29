package vehicle

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/products"
	"strconv"

	"net/http"
)

func Collections(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	cols, err := products.GetAriesVehicleCollections()
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(cols))
}

func Lookup(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	var v products.NoSqlVehicle
	var collection string //e.g. interior/exterior

	//Get collection
	collection = r.FormValue("collection")
	delete(r.Form, "collection")

	// Get vehicle year
	v.Year = r.FormValue("year")
	delete(r.Form, "year")

	// Get vehicle make
	v.Make = r.FormValue("make")
	delete(r.Form, "make")

	// Get vehicle model
	v.Model = r.FormValue("model")
	delete(r.Form, "model")

	// Get vehicle submodel
	v.Style = r.FormValue("style")
	delete(r.Form, "style")

	l, err := products.FindVehicles(v, collection, dtx)
	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(l))
}

func ByCategory(w http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	collection := r.FormValue("collection")
	page, _ := strconv.Atoi(r.FormValue("page"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	var offset int
	if page == 0 {
		offset = 0
	} else if page == 1 {
		offset = 101
	} else {
		offset = page*limit + 1
	}

	res, err := products.FindApplications(collection, offset, limit)
	if err != nil {
		apierror.GenerateError("Trouble finding vehicles.", err, w, r)
		return ""
	}

	return encoding.Must(enc.Encode(res))
}
