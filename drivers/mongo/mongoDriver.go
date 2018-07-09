package mongo

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/nikosEfthias/storageDriver"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var defaultSession *mgo.Session

type crs struct {
	or    []interface{}
	and   bson.M
	queue []func()
	q     *mgo.Query
}

type mongoDriver struct {
	session *mgo.Session
	db      *mgo.Database
	col     *mgo.Collection
	cursor  *crs
}

var _driver = &mongoDriver{}

func (d *mongoDriver) Connect(addr string) (storageDriver.Meta, error) {
	adrs, err := url.Parse(addr)
	if nil != err {
		return nil, err
	}
	if adrs.Host == "" {
		return nil, fmt.Errorf("invalid host")
	}
	var fullAddr = adrs.Hostname() + ":" + adrs.Port()
	if nil != adrs.User {
		usr := adrs.User
		fullAddr = usr.String() + "@" + fullAddr
	}

	session, err := mgo.DialWithTimeout(fullAddr, time.Second*5)
	if nil != err {
		return nil, err
	}
	d.session = session
	return d, nil
}
func init() {
	storageDriver.Register("mongodb", _driver)
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
func (d *mongoDriver) Clone() storageDriver.Meta {
	var a = *d
	return &a
}
func (d *mongoDriver) Driver() (storageDriver.Driver, error) {
	if d.db == nil || d.col == nil || d.session == nil {
		return nil, fmt.Errorf("db ,session or col cannot be nil")
	}
	return d, nil
}
func (d *mongoDriver) AggregateMongo(doc []storageDriver.Document) ([]storageDriver.Document, error) {
	var dc = make([]storageDriver.Document, 0)
	err := d.col.Pipe(doc).All(&dc)
	return dc, err
}
func (d *mongoDriver) Lt(Doc storageDriver.Document) storageDriver.Document {
	var newDoc = make(storageDriver.Document)
	for k, v := range Doc {
		newDoc[k] = storageDriver.Document{"$lt": v}
	}
	return newDoc
}
func (d *mongoDriver) Lte(Doc storageDriver.Document) storageDriver.Document {
	var newDoc = make(storageDriver.Document)
	for k, v := range Doc {
		newDoc[k] = storageDriver.Document{"$lte": v}
	}
	return newDoc
}
func (d *mongoDriver) Gte(Doc storageDriver.Document) storageDriver.Document {
	var newDoc = make(storageDriver.Document)
	for k, v := range Doc {
		newDoc[k] = storageDriver.Document{"$gte": v}
	}
	return newDoc
}
func (d *mongoDriver) Gt(Doc storageDriver.Document) storageDriver.Document {
	var newDoc = make(storageDriver.Document)
	for k, v := range Doc {
		newDoc[k] = storageDriver.Document{"$gt": v}
	}
	return newDoc
}
func (d *mongoDriver) In(key string, values []interface{}) storageDriver.Document {
	return storageDriver.Document{key: storageDriver.Document{"$in": values}}
}
func (d *mongoDriver) Between(key string, values [2]interface{}) storageDriver.Document {
	return storageDriver.Document{key: storageDriver.Document{"$gte": values[0], "$lte": values[1]}}
}
func (d *mongoDriver) Not(Doc storageDriver.Document) storageDriver.Document {
	var newDoc = make(storageDriver.Document)
	for k, v := range Doc {
		newDoc[k] = storageDriver.Document{"$ne": v}
	}
	return newDoc
}
func (d *mongoDriver) Regex(key, value string) storageDriver.Document {
	return storageDriver.Document{key: storageDriver.Document{"$regex": value}}
}
func (d *mongoDriver) Cursor() storageDriver.Cursor {
	d.cursor = new(crs)
	d.cursor.and = bson.M{}
	d.cursor.or = make([]interface{}, 0)
	d.cursor.queue = make([]func(), 0)
	return d
}
func (d *mongoDriver) And(Doc storageDriver.Document) storageDriver.Cursor {
	for k, v := range Doc {
		d.cursor.and[k] = v
	}
	return d
}

func (d *mongoDriver) Or(Doc []interface{}) storageDriver.Cursor {
	d.cursor.or = append(d.cursor.or, Doc...)
	return d
}

func (d *mongoDriver) Select(fields ...string) storageDriver.Cursor {
	fieldMap := bson.M{"_id": 0}
	for _, field := range fields {
		fieldMap[field] = 1
	}
	fn := func() {
		d.cursor.q.Select(fieldMap)
	}
	d.cursor.queue = append(d.cursor.queue, fn)
	return d
}

func (d *mongoDriver) Sort(Doc ...string) storageDriver.Cursor {
	fn := func() {
		d.cursor.q.Sort(Doc...)
	}
	d.cursor.queue = append(d.cursor.queue, fn)
	return d
}

func (d *mongoDriver) Limit(num int) storageDriver.Cursor {
	fn := func() {
		d.cursor.q.Limit(num)
	}
	d.cursor.queue = append(d.cursor.queue, fn)
	return d
}

func (d *mongoDriver) Skip(num int) storageDriver.Cursor {
	fn := func() {
		d.cursor.q.Skip(num)
	}
	d.cursor.queue = append(d.cursor.queue, fn)
	return d
}

func (d *mongoDriver) One(Doc interface{}) error {
	q := getQuery(d.cursor.and, d.cursor.or)
	d.cursor.q = d.col.Find(q)
	for _, fn := range d.cursor.queue {
		fn()
	}
	return d.cursor.q.One(Doc)
}

func (d *mongoDriver) Count(num *int) error {
	q := getQuery(d.cursor.and, d.cursor.or)
	d.cursor.q = d.col.Find(q)
	for _, fn := range d.cursor.queue {
		fn()
	}
	_num, err := d.cursor.q.Count()
	*num = _num
	return err
}

func (d *mongoDriver) All(Doc interface{}) error {
	q := getQuery(d.cursor.and, d.cursor.or)
	d.cursor.q = d.col.Find(q)
	for _, fn := range d.cursor.queue {
		fn()
	}
	return d.cursor.q.All(Doc)
}

func (d *mongoDriver) Distinct(key string, result interface{}) error {
	q := getQuery(d.cursor.and, d.cursor.or)
	d.cursor.q = d.col.Find(q)
	for _, fn := range d.cursor.queue {
		fn()
	}
	return d.cursor.q.Distinct(key, result)
}

func (d *mongoDriver) Save(Query, Doc storageDriver.Document) error {
	_, err := d.col.Upsert(Query, storageDriver.Document{"$set": Doc})
	return err
}
func (d *mongoDriver) Get(query storageDriver.Document) ([]storageDriver.Document, error) {
	var docs = make([]storageDriver.Document, 0)
	err := d.col.Find(query).All(&docs)
	return docs, err
}
func (d *mongoDriver) GetOne(query storageDriver.Document) (storageDriver.Document, error) {
	var doc = make(storageDriver.Document)
	err := d.col.Find(query).One(&doc)
	return doc, err
}
func (d *mongoDriver) Custom(query interface{}) ([]storageDriver.Document, error) {
	return nil, fmt.Errorf("not implemented")
}
func (d *mongoDriver) Update(query, updateFields storageDriver.Document) error {
	return d.col.Update(query, storageDriver.Document{"$set": updateFields})
}
func (d *mongoDriver) UpdateMulti(query, updateFields storageDriver.Document) (int, error) {
	info, err := d.col.UpdateAll(query, storageDriver.Document{"$set": updateFields})
	if nil != info {
		return info.Updated, err
	}
	return 0, err
}
func (d *mongoDriver) Insert(Doc storageDriver.Document) error {
	return d.col.Insert(Doc)
}
func (d *mongoDriver) InsertMulti(docs []storageDriver.Document) error {
	var dcs = make([]interface{}, len(docs))
	for i := range docs {
		dcs[i] = docs[i]
	}
	return d.col.Insert(dcs...)
}
func (d *mongoDriver) InsertMultiNoFail(docs []storageDriver.Document, ErrorOut ...io.Writer) []error {
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
func (d *mongoDriver) Remove(query storageDriver.Document) error {
	return d.col.Remove(query)
}
func getQuery(and storageDriver.Document, or []interface{}) storageDriver.Document {
	if len(or) > 0 {
		and["$or"] = or
	}
	return and
}
