package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/pflag"
)

var newDriverName *string

func init() {
	newDriverName = pflag.StringP("new", "n", "", "create a new driver with the provided name")
	pflag.Parse()
}

func main() {
	_, err := os.Stat(*newDriverName)
	if nil == err {
		fmt.Println("driver with name", *newDriverName, "exists in current folder")
		os.Exit(-1)
	}
	tpl := template.New("default")
	tpl, err = tpl.Parse(_tpl)
	if nil != err {
		panic(err)
	}
	err = os.Mkdir(*newDriverName, 0755)
	if nil != err {
		panic(err)
	}
	f, err := os.OpenFile(*newDriverName+"/"+*newDriverName+".go", os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		panic(err)
	}
	defer f.Close()
	err = tpl.Execute(f, *newDriverName)
	if nil != err {
		panic(err)
	}
}

var _tpl = `
package {{.}}

import (
	"fmt"
	"io"
	"github.com/nikosEfthias/storageDriver"
)

type doc=map[string]interface{} 
type crs struct {
	or    []interface{}
	and   doc
	queue []func()
	q     interface{}
}

type driver struct {
	session interface{} 
	db      interface{} 
	col     interface{} 
	cursor  *crs
}
func(d *driver) Connect (addr string)(storageDriver.Meta,error){
	return fmt.Errorf("not implemented")
} 
func Aggregate(opt *storageDriver.GroupOptions)storageDriver.Aggregator{
	return fmt.Errorf("not implemented")
}
func (d *driver) DB(name string) error {
	return fmt.Errorf("not implemented")
}
func (d *driver) Table(name string) error {
	return fmt.Errorf("not implemented")
}
func (d *driver) Clone() storageDriver.Meta {
	return fmt.Errorf("not implemented")
}
func (d *driver) Driver() (storageDriver.Driver, error) {
	return fmt.Errorf("not implemented")
}
func (d *driver) Lt(Doc doc) doc {
	return fmt.Errorf("not implemented")
}
func (d *driver) Lte(Doc doc) doc {
	return fmt.Errorf("not implemented")
}
func (d *driver) Gte(Doc doc) doc {
	return fmt.Errorf("not implemented")
}
func (d *driver) Gt(Doc doc) doc {
	return fmt.Errorf("not implemented")
}
func (d *driver) In(key string, values []interface{}) doc {
	return fmt.Errorf("not implemented")
}
func (d *driver) Between(key string, values [2]interface{}) doc {
	return fmt.Errorf("not implemented")
}
func (d *driver) Not(Doc doc) doc {
	return fmt.Errorf("not implemented")
}
func (d *driver) Regex(key, value string) doc {
	return fmt.Errorf("not implemented")
}
func (d *driver) Cursor() storageDriver.Cursor {
	return fmt.Errorf("not implemented")
}
func (d *driver) And(Doc doc) storageDriver.Cursor {
	return fmt.Errorf("not implemented")
}

func (d *driver) Or(Doc []interface{}) storageDriver.Cursor {
	return fmt.Errorf("not implemented")
}

func (d *driver) Select(fields ...string) storageDriver.Cursor {
	return fmt.Errorf("not implemented")
}

func (d *driver) Sort(Doc ...string) storageDriver.Cursor {
	return fmt.Errorf("not implemented")
}

func (d *driver) Limit(num int) storageDriver.Cursor {
	return fmt.Errorf("not implemented")
}

func (d *driver) Skip(num int) storageDriver.Cursor {
	return fmt.Errorf("not implemented")
}

func (d *driver) One(Doc interface{}) error {
	return fmt.Errorf("not implemented")
}

func (d *driver) Count(num *int) error {
	return fmt.Errorf("not implemented")
}

func (d *driver) All(Doc interface{}) error {
	return fmt.Errorf("not implemented")
}

func (d *driver) Distinct(key string, result interface{}) error {
	return fmt.Errorf("not implemented")
}
func (d *driver) Save(Query, Doc doc) error {
	return fmt.Errorf("not implemented")
}
func (d *driver) Get(query doc) ([]doc, error) {
	return fmt.Errorf("not implemented")
}
func (d *driver) GetOne(query doc) (doc, error) {
	return fmt.Errorf("not implemented")
}
func (d *driver) Update(query, updateFields doc) error {
	return fmt.Errorf("not implemented")
}
func (d *driver) UpdateMulti(query, updateFields doc) (int, error) {
	return fmt.Errorf("not implemented")
}
func (d *driver) Insert(Doc doc) error {
	return fmt.Errorf("not implemented")
}
func (d *driver) InsertMulti(docs []doc) error {
	return fmt.Errorf("not implemented")
}
func (d *driver) InsertMultiNoFail(docs []doc, ErrorOut ...io.Writer) []error {
	return fmt.Errorf("not implemented")
}
func (d *driver) Remove(query doc) error {
	return fmt.Errorf("not implemented")
}
// vim: ft=go:
`
