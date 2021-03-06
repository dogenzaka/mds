package mds

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Person struct {
		Name  string
		Phone string
	}
	PersonModel struct {
		*Collection
		Name  string        `json:"name"`
		Phone string        `json:"phone"`
		ID    bson.ObjectId `json:"_id" bson:"_id"`
	}
)

func (p *PersonModel) One() {
	defer p.Close()
	name := "Ale"
	p.Find(bson.M{"name": name}).One(&p)

}

func NewPersonalModel(c *Collection) *PersonModel {
	model := &PersonModel{
		&Collection{c.Collection, c.Session},
		"",
		"",
		bson.NewObjectId(),
	}

	return model
}

var DATABASE string = "Mds_Test"
var DATABASE_1 string = "Mds_Test1"

var COLLECTION string = "people"

func TestMongoDB(t *testing.T) {
	ds := &MongoDB{
		Use:  true,
		Dn:   "GoTest",
		Type: MONGODB,
		DialInfo: &mgo.DialInfo{
			Addrs:    []string{"localhost"},
			Database: DATABASE,
		},
	}

	Convey("MongoDB operation", t, func() {

		Convey("Connect()", func() {
			// default session
			session, err := ds.GetSession(false)
			So(session, ShouldEqual, nil)
			So(err, ShouldNotEqual, nil)

			err = ds.Connect()
			So(err, ShouldEqual, nil)
			So(ds.Connected, ShouldEqual, true)

			err = ds.Connect()
			So(err, ShouldNotEqual, true)

		})

		Convey("GetDatabase() (default session)", func() {

			db, err := ds.GetDataBase("", false)
			So(db.Name, ShouldEqual, DATABASE)
			So(err, ShouldEqual, nil)

			db, err = ds.GetDataBase(DATABASE_1, false)
			So(db.Name, ShouldEqual, DATABASE_1)
			So(err, ShouldEqual, nil)

		})

		Convey("GetDatabase() (make session)", func() {

			db, err := ds.GetDataBase("", true)
			So(db.Name, ShouldEqual, DATABASE)
			So(err, ShouldEqual, nil)
			defer db.Session.Close()

			db, err = ds.GetDataBase(DATABASE_1, true)
			So(db.Name, ShouldEqual, DATABASE_1)
			So(err, ShouldEqual, nil)
			defer db.Session.Close()

		})

		Convey("GetCollection (default session)", func() {
			// Original Session
			col, err := ds.GetCollection(COLLECTION, false)
			So(col.Session, ShouldEqual, ds.Session)
			So(err, ShouldEqual, nil)
		})

		Convey("GetCollection (make session)", func() {
			// Original Session
			col, err := ds.GetCollection(COLLECTION, true)
			So(col.Session, ShouldNotEqual, ds.Session)
			So(err, ShouldEqual, nil)
			defer col.Session.Close()

		})

		Convey("GetCollectionWithoutErr (default session, not exit)", func() {
			// Original Session
			col := ds.GetCollectionWithoutErr(COLLECTION, false, false)
			So(col.Session, ShouldEqual, ds.Session)
		})

		Convey("GetCollectionWithoutErr (make session, not exit)", func() {
			// Original Session
			col := ds.GetCollectionWithoutErr(COLLECTION, true, false)
			So(col.Session, ShouldNotEqual, ds.Session)
			defer col.Session.Close()
		})

		Convey("GetCollectionWithoutErr (default session, exit)", func() {
			// Original Session
			col := ds.GetCollectionWithoutErr(COLLECTION, false, true)
			So(col.Session, ShouldEqual, ds.Session)
		})

		Convey("GetCollectionWithoutErr (make session, exit)", func() {
			// Original Session
			col := ds.GetCollectionWithoutErr(COLLECTION, true, true)
			So(col.Session, ShouldNotEqual, ds.Session)
			defer col.Session.Close()
		})

		Convey("Query", func() {

			name := "Ale"
			phone := "+55 53 8116 9639"

			var err error
			c, _ := ds.GetCollection(COLLECTION, false)
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

		Convey("Model", func() {
			c, _ := ds.GetCollection(COLLECTION, true)
			model := NewPersonalModel(c)

			model.One()
			So(model.Phone, ShouldEqual, "+55 53 8116 9639")

			_, err := json.Marshal(model)
			So(err, ShouldBeNil)
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
