package storageDriver

import (
	"testing"
)

var d = new(mapDriver)

func Test_Insert(t *testing.T) {
	d.store = make(map[string]map[string][]Document)
	d.Insert(Document{
		"test": "test",
	})
	if len(d.store) != 1 {
		t.Fail()
	}
}
func Test_Get(t *testing.T) {
	d.store = make(map[string]map[string][]Document)
	for i := 0; i < 500; i++ {
		d.Insert(Document{"num": i})
	}
	docs, err := d.Get(Document{
		"num": 5,
	})
	if nil != err {
		t.Fatal(err)
	}
	if len(docs) != 1 {
		t.Fatal("wrong docs num")
	}
	if docs[0]["num"] != 5 {
		t.Fatal("wrong docs")
	}
	if docs, _ := d.Get(Document{}); len(docs) != 500 {
		t.Fatal("somethign weird is happening")
	}
	docs, err = d.Get(Document{"num": 550})
	if nil == err {
		t.Fatal("error cannot be nil here")
	}
}
func Test_GetOne(t *testing.T) {
	d.store = make(map[string]map[string][]Document)
	for i := 0; i < 500; i++ {
		d.Insert(Document{"num": i})
	}
	doc, err := d.GetOne(Document{"num": 14})
	if nil != err {
		t.Fatal(err)
	}
	if doc["num"] != 14 {
		t.Fatal("somethign weird is happening")
	}
	doc, err = d.GetOne(Document{"num": 550})
	if nil == err {
		t.Fatal("error cannot be nil here")
	}
}
func Test_Custom(t *testing.T) {
	d.store = make(map[string]map[string][]Document)
	if _, err := d.Custom(Document{}); nil == err {
		t.Fatal("this is not implemented yet")
	}
}
func Test_Update(t *testing.T) {
	d.store = make(map[string]map[string][]Document)
	for i := 0; i < 500; i++ {
		d.Insert(Document{"num": i})
	}
	doc, _ := d.GetOne(Document{"num": 3})
	d.Update(doc, Document{"num": 590})
	_, err := d.GetOne(Document{"num": 3})
	if nil == err {
		t.Fatal("error cannot be nil here")
	}
	_, err = d.GetOne(Document{"num": 590})
	if nil != err {
		t.Fatal("something weird is happening")
	}
}
func Test_UpdateMulti(t *testing.T) {

	d.store = make(map[string]map[string][]Document)
	for i := 0; i < 500; i++ {
		d.Insert(Document{"num": 12})
	}
	docs, err := d.Get(Document{"num": 12})
	if nil != err {
		t.Fatal("well, this is unexpected")
	}
	n, err := d.UpdateMulti(Document{"num": 12}, Document{"testField": "test", "num": 10})
	if nil != err {
		t.Fatal("something weird is happening")
	}
	if n != 500 {
		t.Fatal("n is supposed to be 500 WTH!")
	}
	if docs[0]["num"] != 10 || docs[10]["testField"] != "test" {
		t.Fatal("multi update unsuccessful")
	}
}
func Test_Save(t *testing.T) {

	d.store = make(map[string]map[string][]Document)
	err := d.Save(Document{"num": 15}, Document{"test": 12})
	if nil != err {
		t.Fatal("save unsuccessful", err)
	}
	doc, err := d.GetOne(Document{"num": 15})
	if nil != err {
		t.Fatal("save doesnt insert values")
	}
	if doc["num"] != 15 || doc["test"] != 12 {
		t.Fatal("save doesnt insert values correctly")
	}
	err = d.Save(doc, Document{"testValue": "testing"})
	if nil != err {
		t.Fatal(err)
	}
	doc, err = d.GetOne(Document{"num": 15})
	if nil != err {
		t.Fatal(err)
	}
	if doc["testValue"] != "testing" || doc["test"] != 12 {
		t.Fatal("save doesnt update properly")
	}

}
func Test_Remove(t *testing.T) {

	d.store = make(map[string]map[string][]Document)
	for i := 0; i < 500; i++ {
		d.Insert(Document{"num": i})
	}

	err := d.Remove(Document{"num": 10})
	if nil != err {
		t.Fatal("remove unsuccessful", err)
	}

	_, err = d.GetOne(Document{"num": 10})
	if nil == err {
		t.Fatal("remove doesnt remove values")
	}
}

func Test_RemoveAll(t *testing.T) {

	d.store = make(map[string]map[string][]Document)
	for i := 0; i < 500; i++ {
		d.Insert(Document{"num": i})
		d.Insert(Document{"num": i + 1})
		d.Insert(Document{"num": i + 2})
	}

	err := d.RemoveAll(Document{"num": 10})
	if nil != err {
		t.Fatal("remove unsuccessful", err)
	}

	_, err = d.Get(Document{"num": 10})
	if nil == err {
		t.Fatal("remove doesnt remove values")
	}
}
