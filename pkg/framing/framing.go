package framing

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	json "github.com/assemblaj/fastjson"
)

type DB struct {
	parser  json.MappedParserPool
	storage storage
	search  search
}

func NewFramingDB() *DB {
	return &DB{
		storage: buildStorage(),
		search:  buildSearch()}
}

func (db *DB) Load(input io.Reader) error {
	p := db.parser.MappedParserGet()

	b, err := ioutil.ReadAll(input)
	if err != nil {
		fmt.Println(err)
		return err
	}
	d := string(b[:])

	rawdata := p.GetFramingData(d)
	if rawdata == nil {
		return errors.New("Parsing Error")
	}

	if db.storage.isEmpty() {
		db.storage = buildStorage()
	}
	db.storage.updateStorage(rawdata)

	if db.search.isEmpty() {
		db.search = buildSearch()
	}
	kws := parsedKeywords(rawdata)
	db.search.load(kws)

	db.parser.MappedParserPut(p)

	return nil
}

type SearchParams struct {
	value         string
	exact         bool
	caseSensitive bool
}

func (db *DB) getWithSearch(s SearchParams) (bool, *[]Frame) {
	kws := db.search.get(s.value)

	if !s.caseSensitive {
		lcval := strings.ToLower(s.value)
		lckws := db.search.get(lcval)
		kws = *concatUnique(&kws, &lckws)
	}

	if len(kws) == 0 {
		return false, nil
	}

	var frs []Frame
	for _, w := range kws {
		wfrs, ex := db.storage.get(w)
		if ex {
			frs = append(frs, wfrs...)
		}
	}

	return true, &frs
}

func (db *DB) Get(s SearchParams) (bool, *[]Frame) {
	if !s.exact {
		return db.getWithSearch(s)
	}
	frame, exists := db.storage.get(s.value)
	return exists, &frame
}

// TODO :
// func (db *DB) GetRelationships() map[string][]*Frame {

func (db *DB) GetDistincMetaData() map[string][]*Frame {
	mmap := make(map[string][]*Frame)
	db.storage.iter(func(key, value interface{}) {
		fs := value.([]Frame)
		if !hasDictinctFrames(&fs) {
			return
		}
		v := value.(string)
		exist, dmdfs := db.GetDistinct(v)
		if exist {
			mmap[v] = dmdfs
		}
	})
	return mmap
}

func (db *DB) GetDistinct(value string) (bool, []*Frame) {
	ex, fs := db.Get(SearchParams{value: value})
	if !ex {
		return false, nil
	}
	return true, collect(collectDistinct, fs)
}

func (db *DB) GroupByMetaData(value string) (bool, map[string][]*Frame) {
	ex, fs := db.Get(SearchParams{value: value})
	if !ex {
		return false, nil
	}
	return true, group(groupMetadata, fs)
}

type Frame struct {
	Subject   string
	MetaData  []string
	data      json.Value
	Ancestors []string
}

func buildFrame(context *json.Context) Frame {
	return Frame{
		Subject:   context.GetKey(),
		MetaData:  context.GetKeys(),
		data:      context.GetParent(),
		Ancestors: context.GetAncestors()}
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

func (f *Frame) String() string {
	return fmt.Sprintf("«Subject:%s, MetaData:%v»", f.Subject, f.MetaData)
}

func printFmap(fmap map[string][]*Frame) {
	for md, fs := range fmap {
		fmt.Printf("%s:[\n", md)
		for _, f := range fs {
			fmt.Print("\t")
			fmt.Print(f)
			fmt.Print("\n")
		}
		fmt.Println("]")
	}
}
func printFmap2(fmap map[string][]Frame) {
	for md, fs := range fmap {
		fmt.Printf("%s:[\n", md)
		for _, f := range fs {
			fmt.Print("\t")
			fmt.Print(f)
			fmt.Print("\n")
		}
		fmt.Println("]")
	}
}

var metaDataDelimiter = "|"

func MetaDataString(md *[]string) string {
	return strings.Join(*md, metaDataDelimiter)
}

func MetaDataSlice(md string) []string {
	return strings.Split(md, metaDataDelimiter)
}
