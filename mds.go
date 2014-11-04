package mds

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"strings"
)

const (
	// mds version
	VERSION = "0.0.1"

	// type: MongoDB
	MONGODB = "MongoDB"

	// type: Redis
	//REDIS = "Redis"
)

type (
	// Global Mds
	Mds struct {
		DataStores map[string]interface{}
		Version    string
		Setuped    bool
	}
)

// mds to string
func String() string {
	return fmt.Sprintf("Version: %s, DataStores Count:%d", mds.Version, len(mds.DataStores))
}

// Debug flag
var DEBUG = false

// Debug output
func Debug(f string, msgs ...string) {
	if DEBUG {
		fmt.Printf(""+f+"\n", strings.Join(([]string)(msgs), " "))
	}
}

// Global Mds instance
var mds *Mds = &Mds{
	Version:    VERSION,
	DataStores: make(map[string]interface{}),
	Setuped:    false,
}

// Add a datastore
func AddDataStore(dn string, value interface{}) interface{} {
	mds.DataStores[dn] = value
	Debug("Add the datastore. dn=%s\n", dn)
	return value
}

// Get a datastore
func GetDataStore(dn string) (interface{}, error) {
	ds := mds.DataStores[dn]

	if ds == nil {
		return nil, errors.New("Datastore not found")
	}
	return ds, nil

}

// Get the Mongodb datastore
func GetDataStoreMongoDB(dn string) (*MongoDB, error) {
	ds, err := GetDataStore(dn)
	if err != nil {
		return nil, err
	}

	if ret, ok := ds.(*MongoDB); ok {
		return ret, nil
	}

	return nil, errors.New("Internal data store type is invalid")
}

// mds setup (Once)
// connected: セットアップと同時にコネクティングする
func Setup(dss []map[string]interface{}, autoconnected bool) error {

	if mds.Setuped {
		return errors.New("Already setup performed")
	}

	for _, ds := range dss {

		if !ds["Use"].(bool) {
			Debug("Skip the datastore. dn=%s\n", ds["Dn"].(string))
			continue
		}

		switch ds["Type"].(string) {
		case MONGODB:
			mongodb := &MongoDB{}

			err := mapstructure.Decode(ds, mongodb)

			if err != nil {
				return err
			}

			if autoconnected == true {
				Debug("auto-connecting dn=%s\n", ds["Dn"].(string))
				err := mongodb.Connect()
				if err != nil {
					return err
				}
			}

			AddDataStore(mongodb.Dn, mongodb)
		}
	}

	mds.Setuped = true

	return nil
}
