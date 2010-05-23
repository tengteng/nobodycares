package nobodycares

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
    Type string
    Id   string
    Date string
    Body string
}

type CouchFullEntry struct {
    Id   string
    Rev  string
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
            Id string
            Views map[string]interface{}
        }
        vv := map[string]interface{}{ "by_date": map[string]string{ "map": "function(doc) { if (doc.Type == 'Entry') { emit(doc.Date, doc) } }" } }
        ev := EntryView{"_design/entry", vv}
        if id, _, err := db.Insert(ev); err != nil || id != ev.Id {
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
    id_rev := couch.IdAndRev{}
    if rev, err := p.Database.Retrieve(e.Id, &id_rev); err == nil && e.Id == id_rev.Id {
        // Already exists: overwrite
        full_e := CouchFullEntry{e.Id, rev, "Entry", e.Date, e.Body}
        if _, err := p.Database.Edit(full_e); err != nil {
            return err
        }
    } else {
        // Doesn't exist: insert new
        couch_e := CouchEntry{"Entry", e.Id, e.Date, e.Body}
        if _, _, err := p.Database.Insert(couch_e); err != nil {
            return err
        }
    }
    return nil
}

func (p CouchStore) Load(id string) (Entry, os.Error) {
    e := Entry{}
    _, err := p.Database.Retrieve(id, &e)
    return e, err
}

func (p CouchStore) LoadRange(startid string, limit int) ([]Entry, os.Error) {
    exclude_first := false
    if len(startid) > 0 {
        exclude_first = true
        limit = limit + 1
    }
    options := map[string]interface{}{"limit": limit, "descending": true}
    if exclude_first {
        options["startkey_docid"] = startid
    }
    a, err := p.Database.Query("_design/entry/_view/by_date", options)
    if err != nil {
        log.Stderrf("CouchStore: LoadRange: error during Query: %v\n", err)
        return make([]Entry, 0), err
    }
    if exclude_first && len(a) == limit {
        a = a[1:]
    }
    ea := make([]Entry, len(a))
    for i := 0; i < len(a); i++ {
        if _, err := p.Database.Retrieve(a[i], &ea[i]); err != nil {
            log.Stderrf("CouchStore: LoadRange: error retrieving %s: %v\n", a[i], err)
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
