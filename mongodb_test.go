package mds

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Person struct {
		Name  string
		Phone string
	}
)

func TestMongoDB(t *testing.T) {
	ds := &MongoDB{
		Use:  true,
		Name: "GoTest",
		Type: MONGODB,
		DialInfo: &mgo.DialInfo{
			Addrs:    []string{"localhost"},
			Database: "Mds_Test",
		},
	}

	Convey("MongoDB operation", t, func() {

		Convey("Connect()", func() {
			session, err := ds.GetSession()
			So(session, ShouldEqual, nil)
			So(err, ShouldNotEqual, nil)

			err = ds.Connect()
			So(err, ShouldEqual, nil)
			So(ds.Connected, ShouldEqual, true)

			err = ds.Connect()
			So(err, ShouldNotEqual, true)

		})

		Convey("GetCollection (default session)", func() {

			// Original Session
			col_default, err := ds.GetCollection("people")

			So(col_default.Session, ShouldEqual, ds.Session)
			So(err, ShouldEqual, nil)
		})

		Convey("GetCollection (copy session)", func() {
			// New Session
			s, err := ds.CopySession()
			So(err, ShouldEqual, nil)

			var c *Collection

			// error: may args
			c, err = ds.GetCollection("people", &CollectionOption{
				Session: s,
			}, &CollectionOption{})

			So(c, ShouldEqual, nil)
			So(err, ShouldNotEqual, nil)

			// success
			c, err = ds.GetCollection("people", &CollectionOption{
				Session: s,
			})

			So(c.Session, ShouldNotEqual, ds.Session)

			c.Close()
			So(c.Session, ShouldNotEqual, true)
		})

		Convey("Query", func() {

			name := "Ale"
			phone := "+55 53 8116 9639"

			var err error
			c, _ := ds.GetCollection("people")
			err = c.Insert(&Person{
				Name:  name,
				Phone: phone,
			},
				&Person{"Cla", "+55 53 8402 8510"},
			)

			So(err, ShouldEqual, nil)

			result := Person{}
			err = c.Find(bson.M{"name": name}).One(&result)

			So(err, ShouldEqual, nil)
			So(result.Name, ShouldEqual, name)
			So(result.Phone, ShouldEqual, phone)

		})
		Convey("MongoDB.String()", func() {
			So(ds.String(), ShouldNotEqual, "")
		})

		Convey("Close original session", func() {
			ds.Session.Close()
			So(ds.Session, ShouldNotEqual, true)
		})

	})
}
