package mds

import (
	. "github.com/smartystreets/goconvey/convey"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

var DSName string = "First"

func pre() []map[string]interface{} {

	datastores := []map[string]interface{}{
		map[string]interface{}{
			"Use":  true,
			"Name": DSName,
			"Type": "MongoDB",
			"DialInfo": map[string]interface{}{
				"Addrs":    []string{"localhost:27017"},
				"Database": "MDS_TEST",
			},
		},
		map[string]interface{}{
			"Use":  false,
			"Name": "Second",
			"Type": "MongoDB",
		},
	}

	return datastores
}

func preError() []map[string]interface{} {

	datastores := []map[string]interface{}{
		map[string]interface{}{
			"Use":  true,
			"Name": DSName,
			"Type": "MongoDB",
			"DialInfo": &mgo.DialInfo{
				Addrs:    []string{"localhost:27017"},
				Database: "MDS_TEST",
			},
		},
		map[string]interface{}{
			"Use":  false,
			"Name": "Second",
			"Type": "MongoDB",
		},
	}

	return datastores
}

func TestMDS(t *testing.T) {

	Convey("mds operation", t, func() {

		Convey("String", func() {
			So(mds.String(), ShouldNotEqual, "")
		})

		Convey("Setup", func() {
			// Error
			err := Setup(preError())
			So(err, ShouldNotEqual, nil)

			// Success
			err = Setup(pre())
			So(err, ShouldEqual, nil)

			// Error
			err = Setup(pre())
			So(err, ShouldNotEqual, nil)
		})

		Convey("GetDatabase", func() {
			ds, err := mds.GetDataStore(DSName)

			So(err, ShouldEqual, nil)
			So(ds, ShouldNotEqual, nil)

			ds, err = mds.GetDataStore("Error")

			So(err, ShouldNotEqual, nil)
			So(ds, ShouldEqual, nil)

		})

		Convey("GetDatabaseMongoDB", func() {
			ds, err := mds.GetDataStoreMongoDB(DSName)

			So(err, ShouldEqual, nil)
			So(ds, ShouldNotEqual, true)

			ds, err = mds.GetDataStoreMongoDB("Error")

			So(err, ShouldNotEqual, nil)
			So(ds, ShouldEqual, nil)
		})

		Convey("Debug", func() {
			DEBUG = true
			Debug("[Test] %s", "hoge", "foo", "bar")
		})
	})
}
