Nobody Cares is a microblogging engine written in Google Go. It takes its name from the fact that nobody cares about your microblog.

## Installation ##
Install [Go](http://www.golang.org) and set the `$GOROOT` environment variable. Use `goinstall` to install dependent packages: ~~[couch-go](http://couch-go.googlecode.com)~~ [diskv](http://diskv.googlecode.com), [web.go](http://github.com/hoisie/web.go), and [mustache.go](http://github.com/hoisie/mustache.go)

```
$ goinstall diskv.googlecode.com/hg
$ goinstall github.com/hoisie/web.go
$ goinstall github.com/hoisie/mustache.go
```

Then install `nobodycares` by checking it out via Mercurial and building it

```
$ hg clone https://nobodycares.googlecode.com/hg nobodycares
$ cd nobodycares
$ gomake
```

## Running ##
`nobodycares` has a number of runtime flags that dictate behavior; they're all pretty self-explanatory. Run

```
$ ./nobodycares -h
```

to see a list. In the default mode, nobodycares will initialize a diskv disk-backed database in the local directory (under diskv) and store entries there.

You can use the provided `pwhash` command inline when starting `nobodycares`

```
$ ./nobodycares -pwhash=`./pwhash -i=mypassword` ...
```

## Usage ##
Assuming `nobodycares` is configured to run on `127.0.0.1:9999`,

  * Create an Entry by going to `http://127.0.0.1:9999/post`, and using the password whose hash you specified at runtime
  * Edit or delete an Entry by appending (or prefixing) `edit` or `delete`, respectively, to its specific URL, eg.
    * `http://127.0.0.1:9999/10fc39d4/edit`
    * `http://127.0.0.1:9999/delete/10fc39d4`