package mds

import (
	"errors"
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"
)

type (
	// MongoDB
	MongoDB struct {
		Use       bool
		Dn        string
		Type      string
		DialInfo  *mgo.DialInfo
		Session   *mgo.Session
		Connected bool
	}

	// Mongodb#Collection
	Collection struct {
		*mgo.Collection
		Session *mgo.Session `json:"-" bson:"-"`
	}

	// get options
	MongoDBOption struct {
		Session bool
		DbName  string
		ColName string
	}
)

// Return back session into pool
func (c *Collection) Close() {
	c.Session.Close()
}

// MongoDB to string
func (m *MongoDB) String() string {

	return fmt.Sprintf("dn=%s, type=%s, connected=%t, addr=%s, database=%s, session=%p",
		m.Dn,
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
func (m *MongoDB) GetSession(makeSession bool) (*mgo.Session, error) {
	if !m.Connected {
		return nil, errors.New("Do not establish a connection with MongoDB. Advance to the Connect() execution")
	}

	if makeSession {
		s, err := m.CopySession()
		return s, err
	}
	// singleton
	return m.Session, nil
}

// Get a new session
func (m *MongoDB) CopySession() (*mgo.Session, error) {
	s, err := m.GetSession(false)
	if err != nil {
		return nil, err
	}
	// Copy(New) session
	return s.Copy(), nil

}

// Get Database
func (m *MongoDB) GetDataBase(dbname string, makeSession bool) (*mgo.Database, error) {
	s, err := m.GetSession(makeSession)
	if err != nil {
		return nil, err
	}

	return s.DB(dbname), nil
}

// Get Collection
func (m *MongoDB) GetCollection(colname string, makeSession bool) (*Collection, error) {
	s, err := m.GetSession(makeSession)
	if err != nil {
		return nil, err
	}

	// Get collection
	c := s.DB(m.DialInfo.Database).C(colname)

	// wrap
	collection := &Collection{
		c,
		s,
	}

	return collection, nil
}

// Get Collection without any errors
// if exit == true, shutdown this application immediatry.
func (m *MongoDB) GetCollectionWithoutErr(colname string, makeSession bool, exit bool) *Collection {
	c, err := m.GetCollection(colname, makeSession)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if exit {
			panic(err.Error())
		}
		return nil
	}
	return c
}
