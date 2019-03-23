package framing

import (
	"fmt"

	json "github.com/assemblaj/fastjson"
)

type DB struct {
	storage map[string]*Frame
}

func (db *DB) Get(value string) (bool, *Frame) {
	frame, exists := db.storage[value]
	return exists, frame
}

type Frame struct {
	Id        string
	MetaData  []string
	Ancestors *Frame
}

func (db *DB) Load(data string) {
	p := json.MappedParser()

	v, err := p.Parse(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v.String())
}
