package storageDriver

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"gopkg.in/mgo.v2"
)

var defaultSession *mgo.Session

type mongoDriver struct {
	session *mgo.Session
	db      *mgo.Database
	col     *mgo.Collection
}

func (d *mongoDriver) DB(name string) error {
	if nil == d.session {
		return fmt.Errorf("no session was set")
	}
	db := d.session.DB(name)
	d.db = db
	return nil
}
func (d *mongoDriver) Table(name string) error {
	if nil == d.db || nil == d.session {
		return fmt.Errorf("no session or db was set")
	}
	d.col = d.db.C(name)
	return nil
}
func (d *mongoDriver) Clone() Meta {
	var a = *d
	return &a
}
func (d *mongoDriver) Driver() (StorageDriver, error) {
	if d.db == nil || d.col == nil || d.session == nil {
		return nil, fmt.Errorf("db ,session or col cannot be nil")
	}
	return d, nil
}
func (d *mongoDriver) Save(Query, Doc Document) error {
	return fmt.Errorf("not implemented")
}
func (d *mongoDriver) Get(query Document) ([]Document, error) {
	return nil, fmt.Errorf("not implemented")
}
func (d *mongoDriver) GetOne(query Document) (Document, error) {
	return nil, fmt.Errorf("not implemented")
}
func (d *mongoDriver) Custom(query interface{}) ([]Document, error) {
	return nil, fmt.Errorf("not implemented")
}
func (d *mongoDriver) Update(query, updateFields Document) error {
	return fmt.Errorf("not implemented")
}
func (d *mongoDriver) UpdateMulti(query, updateFields Document) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
func (d *mongoDriver) Insert(Doc Document) error {
	return fmt.Errorf("not implemented")
}
func (d *mongoDriver) InsertMulti(docs []Document) error {
	return fmt.Errorf("not implemented")
}
func (d *mongoDriver) InsertMultiNoFail(docs []Document, ErrorOut ...io.Writer) []error {
	return nil
}
func (d *mongoDriver) Remove(query Document) error {
	return fmt.Errorf("not implemented")
}
func NewMongoDriver(addr string) (Meta, error) {
	adrs, err := url.Parse(addr)
	if nil != err {
		return nil, err
	}
	if adrs.Host == "" {
		return nil, fmt.Errorf("invalid host")
	}
	var fullAddr = adrs.Hostname() + ":" + adrs.Port()
	session, err := mgo.DialWithTimeout(fullAddr, time.Second*5)
	if nil != err {
		return nil, err
	}
	return &mongoDriver{
		session: session,
	}, nil
}
