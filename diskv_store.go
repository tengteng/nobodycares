package main

import (
	"os"
	"fmt"
	"diskv.googlecode.com/hg"
	"json"
	"time"
)

type DiskvStore struct {
	Store *diskv.CachedStore
}

func IDTransform(id diskv.KeyType) []string {
	return []string{string(id)} // TODO
}

func NewDiskvStore(basedir string, maxsz uint32) DiskvStore {
	cs, err := diskv.NewCachedStore(basedir, IDTransform, maxsz)
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
	err := json.Unmarshal(buf, e)
	return e, err
}

func generate_id() string {
	return fmt.Sprintf("%x", time.UTC().Seconds())
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
	return []Entry{}, os.NewError("not yet implemented")
}

func (p DiskvStore) Delete(id string, pwhash string) os.Error {
	if pwhash != PasswordHash() {
		return os.NewError("invalid password")
	}
	return p.Store.Erase(diskv.KeyType(id))
}

