package framing

import (
	"sync"

	"github.com/assemblaj/fastjson"
)

type storage struct {
	data map[string][]Frame
	lock sync.RWMutex
}

func buildStorage() storage {
	return storage{
		data: make(map[string][]Frame)}
}

func (st storage) get(value string) ([]Frame, bool) {
	st.lock.RLock()
	fr, ex := st.data[value]
	st.lock.RUnlock()
	return fr, ex
}

func (st storage) put(value string, frs []Frame) {
	st.lock.Lock()
	st.data[value] = frs
	st.lock.Unlock()
}

func (st storage) append(value string, frs []Frame) {
	st.lock.Lock()
	st.data[value] = append(st.data[value], frs...)
	st.lock.Unlock()
}

func (st storage) updateStorage(rawdata map[string][]*fastjson.Context) {
	fmap := buildFrames(rawdata)
	for k, v := range fmap {
		// Append new frames to existing
		// frame list if the maps share a value
		_, ex := st.get(k)
		if ex {
			st.append(k, v)
		} else {
			st.put(k, v)
		}
	}
}

func buildFrames(rawdata map[string][]*fastjson.Context) map[string][]Frame {
	fmap := make(map[string][]Frame)
	for v, cs := range rawdata {
		fmap[v] = []Frame{}
		for _, c := range cs {
			fmap[v] = append(fmap[v], buildFrame(c))
		}
	}
	return fmap
}

func (st storage) iter(f func(key, value interface{})) {
	st.lock.RLock()
	defer st.lock.RUnlock()
	for k, v := range st.data {
		st.lock.RUnlock()
		f(k, v)
		st.lock.RLock()
	}
}

func (st storage) isEmpty() bool {
	return st.data == nil
}
