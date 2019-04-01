package framing

import (
	"os"
	"testing"
)

var basicTestData = "../../testdata/small.json"

func TestLoad(t *testing.T) {
	r, err := os.Open(basicTestData)
	if err != nil {
		t.Fatalf("Unable to open file. ")
	}

	Framing := NewFramingDB()
	err = Framing.Load(r)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}
}

func TestGetFramingMap(t *testing.T) {
	r, err := os.Open(basicTestData)
	if err != nil {
		t.Fatalf("Unable to open file. ")
	}

	Framing := NewFramingDB()
	err = Framing.Load(r)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}

	exists, _ := Framing.Get(SearchParams{value: "file"})
	if !exists {
		t.Fatalf("Frames do not exist for value 'file' ")
	}

}

func TestGetValue(t *testing.T) {
	r, err := os.Open(basicTestData)
	if err != nil {
		t.Fatalf("Unable to open file. ")
	}

	Framing := NewFramingDB()
	err = Framing.Load(r)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}

	exists, frames := Framing.Get(SearchParams{value: "file"})
	if !exists {
		t.Fatalf("Frames do not exist for value 'file' ")
	}

	exists, _ = (*frames)[0].Get("popup")
	if exists == false {
		t.Fatalf("Value not found.")
	}
}

func TestMetaDataConvert(t *testing.T) {
	tst := []string{"1", "2", "3", "4"}
	tmds := MetaDataString(&tst)
	tmdsl := MetaDataSlice(tmds)
	tmds2 := MetaDataString(&tmdsl)
	if tmds != tmds2 {
		t.Fatalf("Value not found.")
	}
}

var complexTestData = "../../testdata/complex.json"

func TestGetDistincMetaData(t *testing.T) {
	r, err := os.Open(complexTestData)
	if err != nil {
		t.Fatalf("Unable to open file. ")
	}

	Framing := NewFramingDB()
	err = Framing.Load(r)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}

	fmap := Framing.GetDistincMetaData()
	if len(fmap) != 0 {
		t.Fatalf("Result map shouldn't have frames for complex.json")
	}

}

func TestGetDistinct(t *testing.T) {
	r, err := os.Open(complexTestData)
	if err != nil {
		t.Fatalf("Unable to open file. ")
	}

	Framing := NewFramingDB()
	err = Framing.Load(r)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}

	e, _ := Framing.GetDistinct("Apple")
	if !e {
		t.Fatalf("Value not found . ")
	}
}

func TestGroupByMetaData(t *testing.T) {
	r, err := os.Open(complexTestData)
	if err != nil {
		t.Fatalf("Unable to open file. ")
	}

	Framing := NewFramingDB()
	err = Framing.Load(r)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}

	e, _ := Framing.GroupByMetaData("Apple")
	if !e {
		t.Fatalf("Value not found . ")
	}
}

func TestSearch(t *testing.T) {
	r, err := os.Open(complexTestData)
	if err != nil {
		t.Fatalf("Unable to open file. ")
	}

	Framing := NewFramingDB()
	err = Framing.Load(r)
	if err != nil {
		t.Fatalf("Framing data load failure. ")
	}
}
