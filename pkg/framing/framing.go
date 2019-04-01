package framing

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/assemblaj/fastjson"
	json "github.com/assemblaj/fastjson"
	trie "github.com/derekparker/trie"
)

type DB struct {
	search  trie.Trie
	storage map[string][]Frame
}

func NewFramingDB() *DB {
	return &DB{
		storage: make(map[string][]Frame)}
}

func (db *DB) Load(input io.Reader) error {
	b, err := ioutil.ReadAll(input)
	if err != nil {
		fmt.Println(err)
		return err
	}
	d := string(b[:])

	p := json.MappedParser()
	rd := p.GetFramingData(d)
	if rd == nil {
		return errors.New("Parsing Error")
	}
	db.storage = buildFrames(rd)

	vals := parsedKeywords(rd)
	db.search = *buildSearch(&vals)
	return nil
}

type SearchParams struct {
	value         string
	exact         bool
	caseSensitive bool
}

func (db *DB) getWithSearch(s SearchParams) (bool, *[]Frame) {
	kws := db.search.FuzzySearch(s.value)

	if !s.caseSensitive {
		lcval := strings.ToLower(s.value)
		lckws := db.search.FuzzySearch(lcval)
		kws = *concatUnique(&kws, &lckws)
	}

	if len(kws) == 0 {
		return false, nil
	}

	var frs []Frame
	for _, w := range kws {
		wfrs, ex := db.storage[w]
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
	frame, exists := db.storage[s.value]
	return exists, &frame
}

// TODO :
// func (db *DB) GetRelationships() map[string][]*Frame {

func (db *DB) GetDistincMetaData() map[string][]*Frame {
	mmap := make(map[string][]*Frame)
	for v, fs := range db.storage {
		if !hasDictinctFrames(&fs) {
			continue
		}
		exist, dmdfs := db.GetDistinct(v)
		if exist {
			mmap[v] = dmdfs
		}
	}
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

func parsedKeywords(m map[string][]*json.Context) []string {
	ks := make([]string, len(m))
	i := 0
	for k := range m {
		ks[i] = k
		i++
	}
	return ks
}

func buildFrames(rd map[string][]*fastjson.Context) map[string][]Frame {
	fmap := make(map[string][]Frame)
	for v, cs := range rd {
		fmap[v] = []Frame{}
		for _, c := range cs {
			fmap[v] = append(fmap[v], buildFrame(c))
		}
	}
	return fmap
}

func buildSearch(vals *[]string) *trie.Trie {
	search := trie.New()
	for _, v := range *vals {
		search.Add(v, nil) // may come to use this someday
	}
	return search
}

func hasDictinctFrames(fs *[]Frame) bool {
	if len(*fs) < 2 {
		return false
	}

	fst := (*fs)[0]
	df := MetaDataString(&(fst.MetaData))

	// return true the first time you see
	// a different frame
	for _, f := range (*fs)[1:] {
		mds := MetaDataString(&(f.MetaData))
		if df != mds {
			return true
		}
	}

	return false
}

func hasGroups(fs *[]Frame) bool {
	mdsl := []string{}
	for _, f := range *fs {
		mds := MetaDataString(&(f.MetaData))
		if inSlice(&mdsl, mds) {
			return true
		}
		mdsl = append(mdsl, mds)
	}
	return false
}

func inSlice(sl *[]string, s string) bool {
	for _, v := range *sl {
		if s == v {
			return true
		}
	}
	return false
}

func appendIfNew(sl *[]string, s string) *[]string {
	var nsl []string
	if !inSlice(sl, s) {
		nsl = append(*sl, s)
		return &nsl
	}
	return sl
}

func concatUnique(sl *[]string, sl2 *[]string) *[]string {
	for _, v := range *sl2 {
		sl = appendIfNew(sl, v)
	}
	return sl
}

type collectFunc func(fr Frame, cmap map[string]bool) bool

func collect(f collectFunc, fs *[]Frame) []*Frame {
	tmp := make(map[string]bool)
	cfs := []*Frame{}
	// do not use pointers from values returned from range
	for i, fr := range *fs {
		mds := MetaDataString(&fr.MetaData)
		if f(fr, tmp) == true {
			tmp[mds] = true
			cfs = append(cfs, &((*fs)[i]))
		}
	}

	return cfs
}

func collectDistinct(f Frame, cmap map[string]bool) bool {
	mds := MetaDataString(&f.MetaData)
	if _, v := cmap[mds]; !v {
		return true
	}
	return false
}

type groupFunc func(fr Frame, gmap map[string][]*Frame) (bool, string)

func group(f groupFunc, fs *[]Frame) map[string][]*Frame {
	gfs := make(map[string][]*Frame)
	for i, fr := range *fs {
		b, k := f(fr, gfs)
		if b == true {
			gfs[k] = append(gfs[k], &((*fs)[i]))
		}
	}
	return gfs
}

func groupMetadata(fr Frame, gmap map[string][]*Frame) (bool, string) {
	mds := MetaDataString(&fr.MetaData)
	return true, mds
}
