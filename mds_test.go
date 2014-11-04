package mds

import (
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

var DnName string = "First"
var Type string = "MongoDB"

func pre() []map[string]interface{} {

	datastores := []map[string]interface{}{
		map[string]interface{}{
			"Use":  true,
			"Dn":   DnName,
			"Type": Type,
			"DialInfo": map[string]interface{}{
				"Addrs":    []string{"localhost:27017"},
				"Database": "MDS_TEST",
			},
		},
		map[string]interface{}{
			"Use":  false,
			"Dn":   "Second",
			"Type": Type,
		},
	}

	return datastores
}

func preError() []map[string]interface{} {

	datastores := []map[string]interface{}{
		map[string]interface{}{
			"Use":  true,
			"Dn":   DnName,
			"Type": Type,
			"DialInfo": &mgo.DialInfo{
				Addrs:    []string{"localhost:27017"},
				Database: "MDS_TEST",
			},
		},
		map[string]interface{}{
			"Use":  false,
			"Name": "Second",
			"Type": Type,
		},
	}

	return datastores
}

func TestMDS(t *testing.T) {

	Convey("mds operation", t, func() {

		Convey("String", func() {
			So(String(), ShouldNotEqual, "")
		})

		Convey("Setup", func() {
			// Error
			err := Setup(preError(), false)
			So(err, ShouldNotEqual, nil)

			// Success
			err = Setup(pre(), true)
			So(err, ShouldEqual, nil)

			// Error
			err = Setup(pre(), false)
			So(err, ShouldNotEqual, nil)
		})

		Convey("Get", func() {
			mds := Get()
			So(mds.Setuped, ShouldBeTrue)
		})

		Convey("GetDatabase", func() {
			ds, err := GetDataStore(DnName)

			So(err, ShouldEqual, nil)
			So(ds, ShouldNotEqual, nil)

			ds, err = GetDataStore("Error")

			So(err, ShouldNotEqual, nil)
			So(ds, ShouldEqual, nil)

		})

		Convey("GetDatabaseMongoDB", func() {
			ds, err := GetDataStoreMongoDB(DnName)

			So(err, ShouldEqual, nil)
			So(ds, ShouldNotEqual, true)

			ds, err = GetDataStoreMongoDB("Error")

			So(err, ShouldNotEqual, nil)
			So(ds, ShouldEqual, nil)
		})

		Convey("Debug", func() {
			DEBUG = true
			Debug("[Test] %s", "hoge", "foo", "bar")
		})
	})
}
