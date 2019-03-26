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
	ID        string
	MetaData  []string
	data      *json.Value
	Ancestors []string
}

func (f *Frame) Get(key string) string {
	var data []byte
	data = f.data.Get(key).MarshalTo(data)
	return string(data[:])
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
		ID:        context.GetKey(),
		MetaData:  context.GetKeys(),
		data:      context.GetParent(),
		Ancestors: context.GetAncestors()}
}
