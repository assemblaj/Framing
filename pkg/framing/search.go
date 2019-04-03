package framing

import (
	"sync"

	json "github.com/assemblaj/fastjson"
	trie "github.com/derekparker/trie"
)

type search struct {
	data trie.Trie
	lock sync.RWMutex
}

func buildSearch() search {
	return search{
		data: *trie.New()}
}

func (sch search) load(keywords []string) {
	for _, w := range keywords {
		sch.put(w, nil)
	}
}

func (sch search) put(key string, meta interface{}) {
	sch.lock.Lock()
	sch.data.Add(key, meta)
	sch.lock.Unlock()
}

func (sch search) get(value string) []string {
	sch.lock.RLock()
	result := sch.data.FuzzySearch(value)
	sch.lock.RUnlock()
	return result
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

func (sch search) isEmpty() bool {
	return sch == (search{})
}
