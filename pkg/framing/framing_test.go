package framing

import (
	"strings"
	"testing"
)

var smallJSON = `
		{"menu": {
			"id": "file",
			"value": "File",
			"popup": {
			"menuitem": [
				{"value": "New", "onclick": "CreateNewDoc()"},
				{"value": "Open", "onclick": "OpenDoc()"},
				{"value": "Close", "onclick": "CloseDoc()"}
			]
			}
		}}	
`

func TestLoad(t *testing.T) {
	s := strings.NewReader(smallJSON)
	Framing := NewFramingDB()
	err := Framing.Load(s)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}
}

func TestGetFramingMap(t *testing.T) {
	s := strings.NewReader(smallJSON)

	Framing := NewFramingDB()
	err := Framing.Load(s)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}

	exists, _ := Framing.Get("file")
	if !exists {
		t.Fatalf("Frames do not exist for value 'file' ")
	}

}

func TestGetValue(t *testing.T) {
	s := strings.NewReader(smallJSON)

	Framing := NewFramingDB()
	err := Framing.Load(s)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}

	exists, frames := Framing.Get("file")
	if !exists {
		t.Fatalf("Frames do not exist for value 'file' ")
	}

	exists, _ = (*frames)[0].Get("popup")
	if exists == false {
		t.Fatalf("Value not found.")
	}
}
