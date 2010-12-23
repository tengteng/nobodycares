# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

NOBODYCARES_GOFILES=\
	nobodycares.go\
	web_interface.go\
	backing_store.go\
	couchdb_store.go\

PWHASH_GOFILES=\
	pwhash.go\
	backing_store.go\

all: nobodycares pwhash

nobodycares: $(NOBODYCARES_GOFILES)
	$(GC) -o $@.$O $(NOBODYCARES_GOFILES)
	$(LD) -o $@ $@.$O

pwhash: $(PWHASH_GOFILES)
	$(GC) -o $@.$O $(PWHASH_GOFILES)
	$(LD) -o $@ $@.$O

clean:
	rm -f *.[$(OS)] $(CLEANFILES) nobodycares pwhash

format: $(NOBODYCARES_GOFILES) $(PWHASH_GOFILES)
	gofmt -w $^
	
