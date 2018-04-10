package storageDriver

import "testing"

var Driver *mongoDriver

func TestNewMongoDriver(t *testing.T) {
	m, err := NewMongoDriver("mongodb://localhost:27017")
	if nil == err && nil == m {
		t.Fatal("both error and the meta cannot be null")
	}
	if nil != err {
		t.Fatal(err)
	}
	m, err = NewMongoDriver("mongodb://localhost:27017/test")
	if nil == m || nil != err {
		t.Fatal("cannot handle the path", err)
	}

	m, err = NewMongoDriver("localhost:27017/test")
	if nil != m || nil == err {
		t.Fatal("missing protocol is supposed to give an error", err)
	}
}

func TestDB(t *testing.T) {
	var m Meta = new(mongoDriver)
	if nil == m.DB("test") {
		t.Fatal("error cannot be nil while theres no session")
	}
	m, err := NewMongoDriver("mongodb://localhost:27017")
	if nil != err {
		t.Fatal(err)
	}
	err = m.DB("storageTest")
	if nil != err {
		t.Fatal(err)
	}
	if nil == m.(*mongoDriver).db {
		t.Fatal("db cannot be nil")
	}
}

func TestTable(t *testing.T) {
	var m Meta = new(mongoDriver)
	if nil == m.Table("test") {
		t.Fatal("error cannot be nil while theres no session and db")
	}
	m, err := NewMongoDriver("mongodb://localhost:27017")
	if nil != err {
		t.Fatal(err)
	}
	if nil == m.Table("test") {
		t.Fatal("without having a db this was supposed to give an error")
	}
	err = m.DB("storageTest")
	if nil != err {
		t.Fatal(err)
	}
	if nil != m.Table("test") {
		t.Fatal(err)
	}
	if nil == m.(*mongoDriver).col {
		t.Fatal("col cannot be nil")
	}
}

func TestClone(t *testing.T) {
	m, err := NewMongoDriver("mongodb://localhost:27017")
	if nil != err {
		t.Fatal(err)
	}
	var mClone = m.Clone()
	if mClone == m {
		t.Fatal("these two are supposed to be different addresses not the same")
	}
}

func TestDriver(t *testing.T) {
	m, err := NewMongoDriver("mongodb://localhost:27017")
	_, err = m.Driver()
	if nil == err {
		t.Fatal("there must be an error here")
	}
	if err := m.DB("test"); nil != err {
		t.Fatal(err)
	}
	_, err = m.Driver()
	if nil == err {
		t.Fatal("there must be an error here")
	}
	if err := m.Table("test"); nil != err {
		t.Fatal(err)
	}
	_, err = m.Driver()
	if nil != err {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	var d = getCleanDb()
	if err := d.Insert(Document{"name": "nikos"}); nil != err {
		t.Fatal(err)
	}
	var data = make(Document)
	err := d.(*mongoDriver).col.Find(Document{"name": "nikos"}).One(data)
	if nil != err {
		t.Fatal(err)
	}
}

func TestInsertMulti(t *testing.T) {
	var data = make([]Document, 100)
	for i := range data {
		data[i] = Document{"num": i, "key": "key"}
	}

	d := getCleanDb()
	err := d.InsertMulti(data)
	if nil != err {
		t.Fatal(err)
	}

	var allData = make([]Document, 0)
	err = d.(*mongoDriver).col.Find(Document{"key": "key"}).All(&allData)
	if nil != err {
		t.Fatal(err)
	}
	if len(allData) != 100 {
		t.Log(allData)
		t.Fatal("invalid data")
	}
}
func getCleanDb() StorageDriver {
	d, err := NewMongoDriver("mongodb://localhost:27017")
	if nil != err {
		panic(err)
	}
	if err := d.DB("testingDB"); nil != err {
		panic(err)
	}
	d.(*mongoDriver).db.DropDatabase()
	d.DB("testingDB")
	d.Table("testing")
	drv, err := d.Driver()
	if nil != err {
		panic(err)
	}
	return drv
}
