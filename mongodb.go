package mds

import (
	"errors"
	"fmt"
	mgo "gopkg.in/mgo.v2"
)

type (
	// MongoDB
	MongoDB struct {
		Use  bool
		Name string
		Type string
		//Url string
		DialInfo  *mgo.DialInfo
		Session   *mgo.Session
		Connected bool
	}

	// Mongodb#Collection
	Collection struct {
		*mgo.Collection
		Session *mgo.Session
	}

	// Options for MongoDB#Collection
	CollectionOption struct {
		Session *mgo.Session
		DbName  string
	}
)

// Return back session into pool
func (c *Collection) Close() {
	c.Session.Close()
}

// MongoDB to string
func (m *MongoDB) String() string {

	return fmt.Sprintf("name=%s, type=%s, connected=%t, addr=%s, database=%s, session=%p",
		m.Name,
		m.Type,
		m.Connected,
		m.DialInfo.Addrs,
		m.DialInfo.Database,
		m.Session,
	)
}

// Connecting to Mongodb
func (m *MongoDB) Connect() error {
	if m.Connected {
		return nil
	}

	session, err := mgo.DialWithInfo(m.DialInfo)

	if err != nil {
		msg := "Failed mongodb connect."
		Debug("%s", msg)
		return errors.New(msg)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	m.Session = session // Original session

	m.Connected = true

	return err
}

// Get a session (Singleton)
func (m *MongoDB) GetSession() (*mgo.Session, error) {
	if !m.Connected {
		return nil, errors.New("Do not establish a connection with MongoDB. Advance to the Connect() execution")
	}

	// singleton
	return m.Session, nil
}

// Get a new session
func (m *MongoDB) CopySession() (*mgo.Session, error) {
	s, err := m.GetSession()
	if err != nil {
		return nil, err
	}
	// Copy(New) session
	return s.Copy(), nil
}

// Get a collection
func (m *MongoDB) GetCollection(colname string, any ...interface{}) (*Collection, error) {

	var dbname string = ""
	var session *mgo.Session = nil
	var err error = nil

	if 0 < len(any) {
		option, ok := any[0].(*CollectionOption)

		if ok == false {
			return nil, errors.New("CollectionOption type assertion failed")
		}

		if option != nil {
			dbname = option.DbName
			session = option.Session
		}

	}

	if 1 < len(any) {
		msg := fmt.Sprintf("Many arguments. GetCollection(%s)", "1<any")
		Debug("%s", msg)
		return nil, errors.New(msg)
	}

	if session == nil { // Use original session
		session, err = m.GetSession()
		if err != nil {
			return nil, err
		}
	}

	// Get collection
	c := session.DB(dbname).C(colname)

	// wrap
	collection := &Collection{
		c,
		session,
	}

	return collection, nil

}
