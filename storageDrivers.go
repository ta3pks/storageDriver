package storageDriver

import "io"

type (
	//Document is general data structure where keys are string and values are anything you want as long as the underlaying driver supports it
	Document = map[string]interface{}
)

//Meta contains the metaData
type Meta interface {
	//DB sets the db of to a given dbname
	DB(dbname string) error
	//Table changes table/collection of the meta
	Table(colName string) error
	//Clone returns a copy of previous meta which uses the same underlaying connection in most cases
	Clone() Meta
	//Driver returns the actual driver which can be later queried
	Driver() (StorageDriver, error)
}
type (
	//Saver inserts or updates data
	Saver interface {
		Save(Query Document, Doc Document) error
	}
	// Getter Either returns a single doc (GetOne) or multiple (Get)
	Getter interface {
		Get(Query Document) ([]Document, error)
		GetOne(Query Document) (Document, error)
		Custom(Query interface{}) ([]Document, error)
	}
	//Updater updates the data and returns the updated document.
	//If there are no documents to update returns an error
	//UpdateMulti returns also updatedDocuments number
	Updater interface {
		Update(Query Document, UpdateFields Document) error
		UpdateMulti(Query Document, UpdateFields Document) (int, error)
	}
	//Inserter inserts the document and returns an error if cannot insert
	//InsertMulti fails on the first error and returns the error stopping the execution
	//On the other hand InsertMultiNoFail doesnt fail on error and returns a slice of errors occured during the execution
	//Also you may pass an optional io.Writer to see the errors in realtime
	Inserter interface {
		Insert(Doc Document) error
		InsertMulti(Docs []Document) error
		InsertMultiNoFail(Docs []Document, ErrorOut ...io.Writer) []error
	}
	// Remover removes the document and returns an error if cannot remove
	Remover interface {
		Remove(Query Document) error
	}
	StorageDriver interface {
		Saver
		Getter
		Updater
		Inserter
		Remover
		AggregateMongo([]map[string]interface{}) ([]Document, error)
		Cursor() Cursor
		Lt(Doc Document) Document
		Gt(Doc Document) Document
		Gte(Doc Document) Document
		Lte(Doc Document) Document
		In(key string, values []interface{}) Document
		Between(key string, values [2]interface{}) Document
		Not(Doc Document) Document
		Regex(key string, value string) Document
	}
)
type Cursor interface {
	And(Doc Document) Cursor
	Or([]interface{}) Cursor
	Select(fieldNames ...string) Cursor
	Sort(Doc ...string) Cursor
	One(Doc interface{}) error
	Limit(num int) Cursor
	Skip(num int) Cursor
	All(Doc interface{}) error
	Count(num *int) error
	Distinct(key string, result interface{}) error
}
type DummyCursor struct{}

func (d DummyCursor) And(Doc Document) Cursor                       { return d }
func (d DummyCursor) Or([]interface{}) Cursor                       { return d }
func (d DummyCursor) Select(fieldNames ...string) Cursor            { return d }
func (d DummyCursor) Sort(Doc ...string) Cursor                     { return d }
func (d DummyCursor) One(Doc interface{}) error                     { return nil }
func (d DummyCursor) Limit(num int) Cursor                          { return d }
func (d DummyCursor) Skip(num int) Cursor                           { return d }
func (d DummyCursor) All(Doc interface{}) error                     { return nil }
func (d DummyCursor) Count(num *int) error                          { return nil }
func (d DummyCursor) Distinct(key string, result interface{}) error { return nil }
