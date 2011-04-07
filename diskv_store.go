package main

import (
	"os"
	"diskv.googlecode.com/hg"
	"json"
)

type DiskvStore struct {
	Store *diskv.CachedStore
}

func IDTransform(id diskv.KeyType) []string {
	return []string{string(id)} // TODO
}

func NewDiskvStore(basedir string, maxsz uint32) (DiskvStore, os.Error) {
	cs, err := diskv.NewCachedStore(basedir, IDTransform, maxsz)
	if err != nil {
		return DiskvStore{nil}, err
	}
	return DiskvStore{cs}, nil
}

func marshal(e Entry) ([]byte, os.Error) {
	return json.Marshal(e)
}

func unmarshal(buf []byte) (Entry, os.Error) {
	e := Entry{}
	err := json.Unmarshal(buf, e)
	return e, err
}

func (p DiskvStore) Save(e Entry, pwhash string) os.Error {
	if pwhash != PasswordHash() {
		return os.NewError("invalid password")
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

