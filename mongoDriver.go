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
	_, err := d.col.Upsert(Query, Doc)
	return err
}
func (d *mongoDriver) Get(query Document) ([]Document, error) {
	var docs = make([]Document, 0)
	err := d.col.Find(query).All(&docs)
	return docs, err
}
func (d *mongoDriver) GetOne(query Document) (Document, error) {
	var doc = make(Document)
	err := d.col.Find(query).One(&doc)
	return doc, err
}
func (d *mongoDriver) Custom(query interface{}) ([]Document, error) {
	return nil, fmt.Errorf("not implemented")
}
func (d *mongoDriver) Update(query, updateFields Document) error {
	return d.col.Update(query, Document{"$set": updateFields})
}
func (d *mongoDriver) UpdateMulti(query, updateFields Document) (int, error) {
	info, err := d.col.UpdateAll(query, Document{"$set": updateFields})
	if nil != info {
		return info.Updated, err
	}
	return 0, err
}
func (d *mongoDriver) Insert(Doc Document) error {
	return d.col.Insert(Doc)
}
func (d *mongoDriver) InsertMulti(docs []Document) error {
	var dcs = make([]interface{}, len(docs))
	for i := range docs {
		dcs[i] = docs[i]
	}
	return d.col.Insert(dcs...)
}
func (d *mongoDriver) InsertMultiNoFail(docs []Document, ErrorOut ...io.Writer) []error {
	var errs = make([]error, 0)
	for _, doc := range docs {
		if err := d.col.Insert(doc); nil != err {
			errs = append(errs, err)
			if len(ErrorOut) > 0 {
				ErrorOut[0].Write([]byte(err.Error()))
			}
		}
	}
	return errs
}
func (d *mongoDriver) Remove(query Document) error {
	return d.col.Remove(query)
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
