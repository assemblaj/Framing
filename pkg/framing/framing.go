package framing

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/assemblaj/fastjson"
	json "github.com/assemblaj/fastjson"
)

type DB struct {
	storage map[string][]Frame
}

func NewFramingDB() *DB {
	return &DB{
		storage: make(map[string][]Frame)}
}

func (db *DB) Get(value string) (bool, *[]Frame) {
	frame, exists := db.storage[value]
	return exists, &frame
}

func (db *DB) Load(input io.Reader) error {
	b, err := ioutil.ReadAll(input)
	if err != nil {
		fmt.Println(err)
		return err
	}
	data := string(b[:])

	p := json.MappedParser()
	rawData := p.GetFramingData(data)
	if rawData == nil {
		return errors.New("Parsing Error")
	}

	db.storage = buildFrames(rawData)
	return nil
}

type Frame struct {
	Subject   string
	MetaData  []string
	data      *json.Value
	Ancestors []string
}

func (f *Frame) Get(key string) (bool, string) {
	var d []byte

	v := f.data.Get(key)
	if v == nil {
		return false, ""
	}

	d = v.MarshalTo(d)
	return true, string(d[:])
}

func buildFrames(rawData map[string][]*fastjson.Context) map[string][]Frame {
	frameMap := make(map[string][]Frame)
	for value, contexts := range rawData {
		for _, context := range contexts {
			frameMap[value] = append(frameMap[value], buildFrame(context))
		}
	}
	return frameMap
}

func buildFrame(context *json.Context) Frame {
	return Frame{
		Subject:   context.GetKey(),
		MetaData:  context.GetKeys(),
		data:      context.GetParent(),
		Ancestors: context.GetAncestors()}
}
