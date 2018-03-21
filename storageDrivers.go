package storageDrivers

import "io"

type (
	//Document is general data structure where keys are string and values are anything you want as long as the underlaying driver supports it
	Document = map[string]interface{}
)
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
	StorageDriver interface {
		Saver
		Getter
		Updater
		Inserter
	}
)

func NewMapDriver() StorageDriver {
	var driver = new(mapDriver)

	driver.store = make([]Document, 0, 100)
	return driver
}
