package storageDrivers

import (
	"fmt"
	"io"
	"sync"
)

type mapDriver struct {
	sync.Mutex
	store []Document
}

func (d *mapDriver) Get(Query Document) ([]Document, error) {
	var docs = make([]Document, 0)
	var match bool = true
	d.Lock()
	defer d.Unlock()
	for _, DBDoc := range d.store {
		for k, v := range Query {
			if val, ok := DBDoc[k]; !ok || val != v {
				match = false
				goto next
			}
		}
		if match {
			docs = append(docs, DBDoc)
		}
	next:
		match = true
	}
	var err error
	if len(docs) == 0 {
		err = fmt.Errorf("no documents found")
	}
	return docs, err
}

func (d *mapDriver) Insert(doc Document) error {
	d.Lock()
	defer d.Unlock()
	d.store = append(d.store, doc)
	return nil
}

func (d *mapDriver) GetOne(Query Document) (Document, error) {
	var match bool = true
	d.Lock()
	defer d.Unlock()
	for _, DBDoc := range d.store {
		for k, v := range Query {
			if val, ok := DBDoc[k]; !ok || val != v {
				match = false
				goto next
			}
		}
		if match {
			return DBDoc, nil
		}
	next:
		match = true
	}
	return nil, fmt.Errorf("no documents found")
}
func (d *mapDriver) Custom(_ interface{}) ([]Document, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (d *mapDriver) InsertMulti(docs []Document) error {
	for _, doc := range docs {
		d.Insert(doc)
	}
	return nil
}

func (d *mapDriver) InsertMultiNoFail(docs []Document, _ ...io.Writer) []error {
	d.InsertMulti(docs)
	return nil
}

func (d *mapDriver) Update(Query Document, UpdatedFields Document) error {
	doc, err := d.GetOne(Query)
	if nil != err {
		return err
	}
	d.Lock()
	defer d.Unlock()
	for k, v := range UpdatedFields {
		doc[k] = v
	}
	return nil
}
func (d *mapDriver) UpdateMulti(Query, UpdatedFields Document) (int, error) {
	docs, err := d.Get(Query)
	if nil != err {
		return 0, err
	}
	d.Lock()
	defer d.Unlock()
	for _, doc := range docs {
		for k, v := range UpdatedFields {
			doc[k] = v
		}
	}
	return len(docs), nil
}
func (d *mapDriver) Save(Query, Doc Document) error {
	doc, err := d.GetOne(Query)
	if nil != err {
		dd := make(Document)
		for k, v := range Query {
			dd[k] = v
		}
		for k, v := range Doc {
			dd[k] = v
		}
		return d.Insert(dd)
	}
	return d.Update(doc, Doc)

	return nil
}
