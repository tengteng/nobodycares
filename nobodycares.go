package nobodycares

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
    Id   string
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
    store BackingStore
    password_hash string
)

// Generates unique ID
func GenerateID() string {
    // taken from Russ Cox 2010-02-24 post to golang-nuts
    f, _ := os.Open("/dev/urandom", os.O_RDONLY, 0)
    b := make([]byte, 16)
    f.Read(b)
    f.Close()
    return fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func Hash(password string) string {
    c := sha256.New()
    io.WriteString(c, password)
    return fmt.Sprintf("%x", c.Sum())
}

func Init(bs BackingStore, pwhash string) {
    store = bs
    if len(pwhash) <= 0 {
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

func LoadRange(startid string, limit int) ([]Entry, os.Error) {
    return store.LoadRange(startid, limit)
}

func Delete(id, pwhash string) os.Error {
    return store.Delete(id, pwhash)
}
