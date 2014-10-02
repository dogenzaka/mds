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
func (e *Mds) String() string {
	return fmt.Sprintf("Version: %s, DataStores Count:%d", e.Version, e.DataStores)
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
func (e *Mds) AddDataStore(name string, value interface{}) interface{} {
	e.DataStores[name] = value
	Debug("Add the datastore. name=%s\n", name)
	return value
}

// Get a datastore
func (e *Mds) GetDataStore(name string) (interface{}, error) {
	ds := e.DataStores[name]

	if ds == nil {
		return nil, errors.New("Datastore not found")
	}
	return ds, nil

}

// Get the Mongodb datastore
func (e *Mds) GetDataStoreMongoDB(name string) (*MongoDB, error) {
	ds, err := e.GetDataStore(name)
	if err != nil {
		return nil, err
	}

	if ret, ok := ds.(*MongoDB); ok {
		return ret, nil
	}

	return nil, errors.New("Internal data store type is invalid")
}

// mds setup (Once)
func Setup(dss []map[string]interface{}) error {
	if mds.Setuped {
		return errors.New("Already setup performed")
	}

	for _, ds := range dss {

		if !ds["Use"].(bool) {
			Debug("Skip the datastore. name=%s\n", ds["Name"].(string))
			continue
		}

		switch ds["Type"].(string) {
		case MONGODB:
			mongodb := &MongoDB{}

			err := mapstructure.Decode(ds, mongodb)
			if err != nil {
				return err
			}
			mds.AddDataStore(mongodb.Name, mongodb)
		}
	}

	mds.Setuped = true

	return nil
}
