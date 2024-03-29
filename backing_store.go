package main

import (
	"os"
	"fmt"
	"io"
	"crypto/sha256"
)

const (
	NCTIME = "2006/01/02 15:04:05"
)

type Entry struct {
	Id   string // semantically same as CouchDB "_id" but must convert!
	Date string
	Body string
}

type BackingStore interface {
	Save(e Entry, pwhash string) os.Error
	Load(id string) (Entry, os.Error)
	LoadRange(startid string, limit int) ([]Entry, os.Error)
	Delete(id string, pwhash string) os.Error
}

var (
	store         BackingStore
	password_hash string
)

func Hash(password string) string {
	c := sha256.New()
	io.WriteString(c, password)
	return fmt.Sprintf("%x", c.Sum())
}

func Init(bs BackingStore, pwhash string) {
	store = bs
	if pwhash == "" {
		panic("invalid password hash")
	}
	password_hash = pwhash
}

func PasswordHash() string {
	return password_hash
}

func Save(e Entry, pwhash string) os.Error {
	return store.Save(e, pwhash)
}

func Load(id string) (Entry, os.Error) {
	return store.Load(id)
}

func LoadRange(fromid string, limit int) ([]Entry, os.Error) {
	return store.LoadRange(fromid, limit)
}

func Delete(id, pwhash string) os.Error {
	return store.Delete(id, pwhash)
}
