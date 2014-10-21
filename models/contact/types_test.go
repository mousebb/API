package contact

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTypes(t *testing.T) {
	var err error
	var ct ContactType

	Convey("Testing Gets", t, func() {
		Convey("Testing GetAll()", func() {
			types, err := GetAllContactTypes()
			So(len(types), ShouldBeGreaterThanOrEqualTo, 0)
			So(err, ShouldBeNil)
		})

		Convey("Testing Get()", func() {
			Convey("ContactType with ID of 0", func() {
				ct = ContactType{}
				err = ct.Get()

				So(ct.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
			})

			Convey("ContactType with non-zero ID", func() {
				ct = ContactType{ID: 1}
				err = ct.Get()

				So(ct.ID, ShouldNotEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Testing Add/Update/Delete", t, func() {
		ct = ContactType{
			Name: "TEST",
		}
		Convey("Add Missing Values", func() {
			Convey("Missing Name", func() {
				ct.Name = ""
				err = ct.Add()
				So(ct.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				ct.Name = "TEST"
			})
		})
		Convey("Update Missing Values", func() {
			Convey("Missing Name", func() {
				ct.Name = ""
				err = ct.Update()
				So(ct.ID, ShouldEqual, 0)
				So(err, ShouldNotBeNil)
				ct.Name = "TEST"
			})
		})

		Convey("Add Valid ContactType", func() {
			err = ct.Add()
			So(ct.ID, ShouldNotEqual, 0)
			So(err, ShouldBeNil)

			Convey("Update Valid ContactType", func() {
				ct.Name = "TESTER"
				err = ct.Update()
				So(ct.Name, ShouldEqual, "TESTER")
				So(err, ShouldBeNil)

				Convey("Test getReceiversByType", func() {
					var cr ContactReceiver
					cr.LastName = "testLName"
					cr.Email = "testEmail@test.com"
					err = cr.Add()
					So(err, ShouldBeNil)

					err = cr.Get()

					Convey("Test Join", func() {
						err = cr.CreateTypeJoin(ct)
						t.Log(err)
						crs, err := ct.GetReceivers()
						if err != sql.ErrNoRows {
							So(err, ShouldBeNil)
							So(len(crs), ShouldBeGreaterThanOrEqualTo, 1)
						}
						//cleanup
						cr.DeleteTypeJoin(ct)
						So(err, ShouldBeNil)
						cr.Delete()
						So(err, ShouldBeNil)

					})

				})

				Convey("Delete Valid ContactType", func() {
					err = ct.Delete()
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("Delete Empty ContactType", func() {
			ctype := ContactType{}
			err = ctype.Delete()
			So(err, ShouldNotBeNil)
		})
	})
}
