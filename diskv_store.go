package main

import (
	"os"
	"diskv.googlecode.com/hg"
)

const (
	max_cache_sz = 1024*1024*100 // 100MB
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
	return []byte{}, os.NewError("not yet implemented")
}

func unmarshal(buf []byte) (Entry, os.Error) {
	return Entry{}, os.NewError("not yet implemented")
}

func (p DiskvStore) Save(e Entry, pwhash string) os.Error {
	if pwhash != PasswordHash() {
		return os.NewError("invalid password")
	}
	buf, err := marshal(e)
	if err != nil {
		return err
	}
	if err = p.Store.Write(diskv.KeyType(e.Id), buf); err != nil {
		return err
	}
	return nil
}

func (p DiskvStore) Load(id string) (Entry, os.Error) {
	buf, err := p.Store.Read(diskv.KeyType(id))
	if err != nil {
		return Entry{}, err
	}
	e, err := unmarshal(buf)
	if err != nil {
		return Entry{}, err
	}
	return e, nil
}

func (p DiskvStore) LoadRange(startid string, limit int) ([]Entry, os.Error) {
	return []Entry{}, os.NewError("not yet implemented")
}

func (p DiskvStore) Delete(id string, pwhash string) os.Error {
	if pwhash != PasswordHash() {
		return os.NewError("invalid password")
	}
	if err := p.Store.Erase(diskv.KeyType(id)); err != nil {
		return err
	}
	return nil
}

