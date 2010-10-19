package main

import (
	"fmt"
	"flag"
	"strconv"
	"log"
	"github.com/hoisie/web.go"
)

var title *string = flag.String("title", "Nobody Cares", "title of the microblog")
var url *string = flag.String("url", "http://127.0.0.1", "base URL for the microblog (important!)")
var host *string = flag.String("host", "0.0.0.0", "web server bind host/address")
var port *int = flag.Int("port", 9999, "web server bind port")
var max_entries *int = flag.Int("max_entries", 10, "max entries per page")
var pwhash *string = flag.String("pwhash", "", "sha256 hash for password")
var couch_host *string = flag.String("couch_host", "127.0.0.1", "CouchDB server")
var couch_port *int = flag.Int("couch_port", 5984, "CouchDB port")
var couch_name *string = flag.String("couch_name", "ncdb", "CouchDB database name")

func main() {
	flag.Parse()
	Init(NewCouchStore(*couch_host, strconv.Itoa(*couch_port), *couch_name), *pwhash)
	log.Stderrf("nobodycares engine starting up...")

	web.Get("/", get_root)
	web.Get("/from/([0-9a-f]+)", get_from)
	web.Get("/post", get_post)
	web.Get("/edit/([0-9a-f]+)", get_edit)
	web.Get("/([0-9a-f]+)/edit", get_edit)
	web.Get("/delete/([0-9a-f]+)", get_delete)
	web.Get("/([0-9a-f]+)/delete", get_delete)
	web.Get("/([0-9a-f]+)", get_specific_id)
	web.Get("/rss", get_rss)
	web.Get("/css/(.*)", get_css)

	web.Post("/post", post_post)
	web.Post("/edit", post_edit)
	web.Post("/delete", post_delete)

	web.Run(fmt.Sprintf("%s:%d", *host, *port))
}
