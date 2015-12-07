package products

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/curt-labs/API/middleware"
)

type Input struct {
	Year  string
	Make  string
	Model string
	Style string
	Part  string
}

type Application struct {
	Year  string `bson:"year"`
	Make  string `bson:"make"`
	Model string `bson:"model"`
	Style string `bson:"style"`
	Parts []int  `bson:"parts"`
}

var (
	VehicleApplications map[string]Application
	PartConversion      map[string]int
)

func Import(ctx *middleware.APIContext, f multipart.File, collectionName string) ([]error, []error, error) {
	var err error
	var conversionErrs []error
	var insertErrs []error
	VehicleApplications = make(map[string]Application)
	PartConversion = make(map[string]int)

	es, err := CaptureCsv(f)
	if err != nil {
		return conversionErrs, insertErrs, err
	}

	for _, e := range es {
		if cerr := ConvertToApplication(ctx, e); cerr != nil {
			conversionErrs = append(conversionErrs, cerr)
			continue
		}
	}

	ClearCollection(ctx, collectionName)

	for _, app := range VehicleApplications {
		if ierr := IntoDB(ctx, app, collectionName); ierr != nil {
			insertErrs = append(insertErrs, ierr)
			continue
		}
	}

	return conversionErrs, insertErrs, err
}

//Csv to Struct
func CaptureCsv(f multipart.File) ([]Input, error) {
	var e Input
	var es []Input

	reader := csv.NewReader(f)

	lines, err := reader.ReadAll()
	if err != nil {
		return es, err
	}

	for _, line := range lines {
		if len(line) < 5 {
			continue
		}
		e = Input{
			Make:  strings.ToLower(strings.TrimSpace(line[0])),
			Model: strings.ToLower(strings.TrimSpace(line[1])),
			Style: strings.ToLower(strings.TrimSpace(line[2])),
			Part:  strings.TrimSpace(line[3]),
			Year:  strings.ToLower(strings.TrimSpace(line[4])),
		}

		es = append(es, e)
	}
	return es, nil
}

//Convert Input ot Applications array
func ConvertToApplication(ctx *middleware.APIContext, e Input) error {
	var partID int

	if partID = PartConversion[e.Part]; partID == 0 {

		stmt, err := ctx.DB.Prepare("select partID from Part where oldPartNumber = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()

		if err := stmt.QueryRow(e.Part).Scan(&partID); err != nil || partID == 0 {
			return fmt.Errorf("invalid part: %s", e.Part)
		}

		PartConversion[e.Part] = partID
	}

	tmp := Application{
		Parts: []int{partID},
		Year:  e.Year,
		Make:  e.Make,
		Model: e.Model,
		Style: e.Style,
	}

	idx := VehicleApplications[tmp.string()]
	if idx.Year == "" {
		VehicleApplications[tmp.string()] = tmp
		return nil
	}

	idx.Parts = append(idx.Parts, partID)
	VehicleApplications[tmp.string()] = idx

	return nil
}

//Dump into mongo
func IntoDB(ctx *middleware.APIContext, app Application, collectionName string) error {
	return ctx.AriesSession.DB(ctx.AriesMongoDatabase).C(collectionName).Insert(app)
}

//Drop collection specified
func ClearCollection(ctx *middleware.APIContext, name string) error {
	return ctx.AriesSession.DB(ctx.AriesMongoDatabase).C(name).DropCollection()
}

//ToString
func (a *Application) string() string {
	return fmt.Sprintf("%s%s%s%s", a.Year, a.Make, a.Model, a.Style)
}
