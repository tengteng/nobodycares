package main

import (
	"os"
	"fmt"
	"diskv.googlecode.com/hg"
	"json"
	"time"
)

type DiskvStore struct {
	Store *diskv.OrderedCachedStore
}

func IDTransform(id diskv.KeyType) []string {
	return []string{string(id)} // TODO
}

func UTCGreater(a, b interface{}) bool {
	var ai, bi int64
	fmt.Sscanf(string(a.(diskv.KeyType)), "%x", &ai)
	fmt.Sscanf(string(b.(diskv.KeyType)), "%x", &bi)
	return ai > bi
}

func NewDiskvStore(basedir string, maxsz uint32) DiskvStore {
	cs, err := diskv.NewOrderedCachedStore(basedir, IDTransform, maxsz, UTCGreater)
	if err != nil {
		panic(fmt.Sprintf("couldn't create diskv store: %s", err))
	}
	return DiskvStore{cs}
}

func marshal(e Entry) ([]byte, os.Error) {
	return json.Marshal(e)
}

func unmarshal(buf []byte) (Entry, os.Error) {
	e := Entry{}
	err := json.Unmarshal(buf, &e)
	return e, err
}

func generate_id() string {
	i := time.UTC().Seconds()
	s := fmt.Sprintf("%x", i)
	return s
}

func (p DiskvStore) Save(e Entry, pwhash string) os.Error {
	if pwhash != PasswordHash() {
		return os.NewError("invalid password")
	}
	if len(e.Id) <= 0 {
		e.Id = generate_id()
	}
	buf, err := marshal(e)
	if err != nil {
		return err
	}
	return p.Store.Write(diskv.KeyType(e.Id), buf)
}

func (p DiskvStore) Load(id string) (Entry, os.Error) {
	buf, err := p.Store.Read(diskv.KeyType(id))
	if err != nil {
		return Entry{}, err
	}
	return unmarshal(buf)
}

func (p DiskvStore) LoadRange(startid string, limit int) ([]Entry, os.Error) {
	keys, err := p.Store.KeysFrom(diskv.KeyType(startid), limit)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
		return []Entry{}, err
	}
	entries := make([]Entry, len(keys))
	for i, k := range keys {
		buf, err := p.Store.Read(k)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
			return []Entry{}, err
		}
		entry, err := unmarshal(buf)
		if err != nil {
			panic(fmt.Sprintf("%v", err))
			return []Entry{}, err
		}
		entries[i] = entry
	}
	return entries, nil
}

func (p DiskvStore) Delete(id string, pwhash string) os.Error {
	if pwhash != PasswordHash() {
		return os.NewError("invalid password")
	}
	return p.Store.Erase(diskv.KeyType(id))
}

