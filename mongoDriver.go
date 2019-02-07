package storageDriver

import (
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var defaultSession *mgo.Session

type crs struct {
	sync.Mutex
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
func (d *mongoDriver) AggregateMongo(doc []Document) ([]Document, error) {
	var dc = make([]Document, 0)
	err := d.col.Pipe(doc).All(&dc)
	return dc, err
}
func (d *mongoDriver) Lt(Doc Document) Document {
	var newDoc = make(Document)
	for k, v := range Doc {
		newDoc[k] = Document{"$lt": v}
	}
	return newDoc
}
func (d *mongoDriver) Lte(Doc Document) Document {
	var newDoc = make(Document)
	for k, v := range Doc {
		newDoc[k] = Document{"$lte": v}
	}
	return newDoc
}
func (d *mongoDriver) Gte(Doc Document) Document {
	var newDoc = make(Document)
	for k, v := range Doc {
		newDoc[k] = Document{"$gte": v}
	}
	return newDoc
}
func (d *mongoDriver) Gt(Doc Document) Document {
	var newDoc = make(Document)
	for k, v := range Doc {
		newDoc[k] = Document{"$gt": v}
	}
	return newDoc
}
func (d *mongoDriver) In(key string, values []interface{}) Document {
	return Document{key: Document{"$in": values}}
}
func (d *mongoDriver) Between(key string, values [2]interface{}) Document {
	return Document{key: Document{"$gte": values[0], "$lte": values[1]}}
}
func (d *mongoDriver) Not(Doc Document) Document {
	var newDoc = make(Document)
	for k, v := range Doc {
		newDoc[k] = Document{"$ne": v}
	}
	return newDoc
}
func (d *mongoDriver) Regex(key, value string, options string) Document {
	return Document{key: Document{"$regex": value, "$options": options}}
}
func (d *mongoDriver) Cursor() Cursor {
	d.cursor = new(crs)
	d.cursor.and = bson.M{}
	d.cursor.or = make([]interface{}, 0)
	d.cursor.queue = make([]func(), 0)
	return d
}
func (d *mongoDriver) And(Doc Document) Cursor {
	d.cursor.Lock()
	defer d.cursor.Unlock()
	for k, v := range Doc {
		d.cursor.and[k] = v
	}
	return d
}

func (d *mongoDriver) Or(Doc []interface{}) Cursor {
	d.cursor.or = append(d.cursor.or, Doc...)
	return d
}

func (d *mongoDriver) Select(fields ...string) Cursor {
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

func (d *mongoDriver) Sort(Doc ...string) Cursor {
	fn := func() {
		d.cursor.q.Sort(Doc...)
	}
	d.cursor.queue = append(d.cursor.queue, fn)
	return d
}

func (d *mongoDriver) Limit(num int) Cursor {
	fn := func() {
		d.cursor.q.Limit(num)
	}
	d.cursor.queue = append(d.cursor.queue, fn)
	return d
}

func (d *mongoDriver) Skip(num int) Cursor {
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

func (d *mongoDriver) Save(Query, Doc Document) error {
	_, err := d.col.Upsert(Query, Document{"$set": Doc})
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

func (d *mongoDriver) RemoveAll(query Document) error {
	_, err := d.col.RemoveAll(query)
	return err
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
	if nil != adrs.User {
		usr := adrs.User
		fullAddr = usr.String() + "@" + fullAddr
	}

	session, err := mgo.DialWithTimeout(fullAddr, time.Second*5)
	if nil != err {
		return nil, err
	}
	return &mongoDriver{
		session: session,
	}, nil
}
func getQuery(and Document, or []interface{}) Document {

	if len(or) > 0 {
		and["$or"] = or
	}
	return and
}
