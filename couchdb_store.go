package main

import (
	"os"
	"couch-go.googlecode.com/hg"
	"log"
	"fmt"
)

type CouchStore struct {
	Database couch.Database
}

type CouchEntry struct {
	Id   string "_id"
	Type string
	Date string
	Body string
}

type CouchFullEntry struct {
	Id   string "_id"
	Rev  string "_rev"
	Type string
	Date string
	Body string
}

func NewCouchStore(host, port, dbname string) CouchStore {
	db, err := couch.NewDatabase(host, port, dbname)
	if err != nil {
		panic(fmt.Sprintf("couldn't create or load CouchDB: %v", err))
	}
	id_rev := new(couch.IdAndRev)
	if _, err := db.Retrieve("_design/entry", id_rev); err != nil {
		type EntryView struct {
			Id    string
			Views map[string]interface{}
		}
		vv := map[string]interface{}{"by_date": map[string]string{"map": "function(doc) { if (doc.Type == 'Entry') { emit(doc.Date, doc) } }"}}
		ev := EntryView{"_design/entry", vv}
		if _, err := db.Edit(ev); err != nil {
			panic(fmt.Sprintf("couldn't Insert necessary view to CouchDB: %v", err))
		}
		if _, err := db.Retrieve("_design/entry", id_rev); err != nil {
			panic(fmt.Sprintf("couldn't Retrieve necessary view from CouchDB: %v", err))
		}
	}
	return CouchStore{db}
}

func (p CouchStore) Save(e Entry, pwhash string) os.Error {
	if pwhash != PasswordHash() {
		return os.NewError("invalid password")
	}
	// Save must overwrite existing Entry, if it exists
	id_rev := new(couch.IdAndRev)
	if rev, err := p.Database.Retrieve(e.Id, &id_rev); err == nil && e.Id == id_rev.Id {
		// Already exists: overwrite
		full_e := CouchFullEntry{Id: e.Id, Rev: rev, Type: "Entry", Date: e.Date, Body: e.Body}
		log.Printf("Save: overwriting %s: %v", e.Id, full_e)
		if _, err := p.Database.Edit(full_e); err != nil {
			return err
		}
	} else {
		// Doesn't exist: insert new
		log.Printf("Save: inserting new Entry: %s", e.Id)
		couch_e := CouchEntry{Id: e.Id, Type: "Entry", Date: e.Date, Body: e.Body}
		if _, _, err := p.Database.Insert(couch_e); err != nil {
			return err
		}
	}
	return nil
}

func (p CouchStore) Load(id string) (Entry, os.Error) {
	e := CouchEntry{}
	_, err := p.Database.Retrieve(id, &e)
	return Entry{Id: e.Id, Date: e.Date, Body: e.Body}, err
}

func (p CouchStore) LoadRange(fromid string, limit int) ([]Entry, os.Error) {
	fromdate := ""
	if len(fromid) > 0 {
		// need fromdate as well
		fromentry := new(CouchEntry)
		if _, err := p.Database.Retrieve(fromid, fromentry); err == nil {
			fromdate = fromentry.Date
		} else {
			log.Printf("CouchStore: LoadRange: error retrieving %s: %v\n", fromid, err)
		}
	}
	options := map[string]interface{}{"limit": limit, "descending": true}
	if len(fromid) > 0 && len(fromdate) > 0 {
		options["startkey"] = fromdate
		options["startkey_docid"] = fromid
	}
	a := make(map[int]string)
	err := p.Database.Query("_design/entry/_view/by_date", options, a)
	if err != nil {
		log.Printf("CouchStore: LoadRange: error during Query: %v\n", err)
		return make([]Entry, 0), err
	}
	ea := make([]Entry, len(a))
	for i := 0; i < len(a); i++ {
		ce := new(CouchEntry)
		if _, err := p.Database.Retrieve(a[i], ce); err != nil {
			log.Printf("CouchStore: LoadRange: error retrieving %s: %v\n", a[i], err)
		} else {
			ea[i] = Entry{Id: ce.Id, Date: ce.Date, Body: ce.Body}
		}
	}
	return ea, nil
}

func (p CouchStore) Delete(id, pwhash string) os.Error {
	if pwhash != PasswordHash() {
		return os.NewError("invalid password")
	}
	id_rev := couch.IdAndRev{}
	if rev, err := p.Database.Retrieve(id, &id_rev); err == nil && id == id_rev.Id {
		if err := p.Database.Delete(id_rev.Id, rev); err != nil {
			return err
		}
	} else {
		return os.NewError("no such id")
	}
	return nil
}
