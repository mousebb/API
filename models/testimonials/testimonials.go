package testimonials

import (
	"github.com/curt-labs/API/middleware"

	"database/sql"
	"errors"
	"time"
)

const (
	testimonialFields = ` t.testimonialID, t.rating, t.title, t.testimonial, t.dateAdded, t.approved, t.active, t.first_name, t.last_name, t.location, t.brandID `
)

var (
	getAllTestimonialsStmt = `select ` + testimonialFields + ` from Testimonial as t
																	Join ApiKeyToBrand as akb on akb.brandID = t.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (t.brandID = ? OR 0=?)) && t.active = 1 && t.approved = 1 order by t.dateAdded desc`
	getTestimonialsByPageStmt = `select ` + testimonialFields + ` from Testimonial as t
																	Join ApiKeyToBrand as akb on akb.brandID = t.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (t.brandID = ? OR 0=?)) && t.active = 1 && t.approved = 1 order by t.dateAdded desc limit ?,?`
	getRandomTestimonalsStmt = `select ` + testimonialFields + ` from Testimonial as t
																	Join ApiKeyToBrand as akb on akb.brandID = t.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (t.brandID = ? OR 0=?)) && t.active = 1 && t.approved = 1 order by Rand() limit ?`
	getTestimonialStmt = `select ` + testimonialFields + ` from Testimonial as t
																	Join ApiKeyToBrand as akb on akb.brandID = t.brandID
																	Join ApiKey as ak on akb.keyID = ak.id
																	where (ak.api_key = ? && (t.brandID = ? OR 0=?)) && t.testimonialID = ?`
	createTestimonial = `insert into Testimonial (rating, title, testimonial, dateAdded, approved, active, first_name, last_name, location, brandID) values (?,?,?,?,?,?,?,?,?,?)`
	updateTestimonial = `update Testimonial set rating = ?, title = ?, testimonial = ?, approved = ?, active = ?, first_name = ?, last_name = ?, location = ?, brandID = ? where testimonialID = ?`
	deleteTestimonial = `delete from Testimonial where testimonialID = ?`
)

type Testimonials []Testimonial
type Testimonial struct {
	ID        int       `json:"id,omitempty" xml:"id,omitempty"`
	Rating    float64   `json:"rating,omitempty" xml:"rating,omitempty"`
	Title     string    `json:"title,omitempty" xml:"title,omitempty"`
	Content   string    `json:"content,omitempty" xml:"content,omitempty"`
	DateAdded time.Time `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	Approved  bool      `json:"approved,omitempty" xml:"approved,omitempty"`
	Active    bool      `json:"active,omitempty" xml:"active,omitempty"`
	FirstName string    `json:"firstName,omitempty" xml:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty" xml:"lastName,omitempty"`
	Location  string    `json:"location,omitempty" xml:"location,omitempty"`
	BrandID   int       `json:"brandId,omitempty" xml:"brandId,omitempty"`
}

func GetAllTestimonials(ctx *middleware.APIContext, page int, count int, randomize bool) (tests Testimonials, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if page == 0 && count == 0 {
		stmt, err = ctx.DB.Prepare(getAllTestimonialsStmt)
		if err != nil {
			return
		}
		defer stmt.Close()
		rows, err = stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	} else if randomize {
		stmt, err = ctx.DB.Prepare(getRandomTestimonalsStmt)
		if err != nil {
			return
		}
		defer stmt.Close()
		rows, err = stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID, count)
	} else {
		stmt, err = ctx.DB.Prepare(getTestimonialsByPageStmt)
		if err != nil {
			return
		}
		defer stmt.Close()
		rows, err = stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID, page, count)
	}

	if err != nil {
		return
	}

	for rows.Next() {
		var t Testimonial
		err = rows.Scan(
			&t.ID,
			&t.Rating,
			&t.Title,
			&t.Content,
			&t.DateAdded,
			&t.Approved,
			&t.Active,
			&t.FirstName,
			&t.LastName,
			&t.Location,
			&t.BrandID,
		)
		if err != nil {
			return
		}

		tests = append(tests, t)
	}
	defer rows.Close()
	return
}

func (t *Testimonial) Get(ctx *middleware.APIContext) error {
	if t.ID == 0 {
		return errors.New("Invalid testimonial ID")
	}

	stmt, err := ctx.DB.Prepare(getTestimonialStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID, t.ID).Scan(
		&t.ID,
		&t.Rating,
		&t.Title,
		&t.Content,
		&t.DateAdded,
		&t.Approved,
		&t.Active,
		&t.FirstName,
		&t.LastName,
		&t.Location,
		&t.BrandID,
	)

	return err
}

func (t *Testimonial) Create(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(createTestimonial)
	if err != nil {
		return err
	}
	defer stmt.Close()
	t.DateAdded = time.Now()

	res, err := stmt.Exec(t.Rating, t.Title, t.Content, t.DateAdded, t.Approved, t.Active, t.FirstName, t.LastName, t.Location, t.BrandID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	t.ID = int(id)
	return nil
}

func (t *Testimonial) Update(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(updateTestimonial)
	if err != nil {
		return err
	}
	defer stmt.Close()

	t.DateAdded = time.Now()
	_, err = stmt.Exec(t.Rating, t.Title, t.Content, t.Approved, t.Active, t.FirstName, t.LastName, t.Location, t.BrandID, t.ID)

	return err
}

func (t *Testimonial) Delete(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(deleteTestimonial)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.ID)

	return err
}
