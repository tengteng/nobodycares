package main

import (
    "log"
    "fmt"
    "time"
    "bytes"
    "template"
    "nobodycares"
    "github.com/hoisie/web.go"
    "github.com/hoisie/mustache.go"
    "flag"
    "strconv"
)

var title *string = flag.String("title", "Nobody Cares", "title of the microblog")
var url *string = flag.String("url", "http://127.0.0.1", "base URL for the microblog (important!)")
var host *string = flag.String("host", "0.0.0.0", "web server binding host/address")
var port *int = flag.Int("port", 9999, "web server port")
var max_entries *int = flag.Int("max_entries", 10, "max entries per page")
var password_hash *string = flag.String("pwhash", "a5ad895656074bb12930374348bf903460016bcf430bf7039d7e34f0c505a7b1", "sha256 hash for password")
var couch_host *string = flag.String("couch_host", "127.0.0.1", "CouchDB server")
var couch_port *int = flag.Int("couch_port", 5984, "CouchDB port")
var couch_name *string = flag.String("couch_name", "ncdb", "CouchDB database name")

const page_str = `
<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html;charset=utf-8" />
  <title>%s</title>
  <link rel="stylesheet" href="/css/master.css" type="text/css" media="screen" />
  <link rel="alternate" type="application/rss+xml" title="RSS feed" href="/rss" />
</head>
<body>
<div class="c">
<h1>%s</h1>
<br/>
%s
</div>
</body>
</html>
`

const css_str = `
body { font-family: monospace; text-transform:uppercase; }
a { color:#777; }
.c { width: 400px; margin: 0 auto; margin-top: 50px; }
`

const edit_form_template = `
<form method="POST" action="{{action}}">
    <input type="hidden" name="form_id" value="{{id}}" />
    <p>
        <input type="text" name="form_date" value="{{date}}" />
    </p>
    <p>
        <textarea name="form_body" style="width:350px; height:200px;">{{body}}</textarea>
    </p>
    <p>
        <input type="password" name="form_password" value="" />
    </p>
    <p>
        <input type="Submit" value="{{button_label}}" />
    </p>
</form>
`

const entry_template = `
<p>
    <a href="/{{Id}}">{{Date}}</a><br/>
    {{Body}}
</p>
`

const rss_template = `
<rss version="2.0">
    <channel>
        <title>{{title}}</title>
        <link>{{url}}</link>
        <language>en-us</language>
        <pubDate>{{most_recent_date}}</pubDate>
        <lastBuildDate>{{most_recent_date}}</lastBuildDate>
        <docs>http://blogs.law.harvard.edu/tech/rss</docs>
        <generator>NobodyCares microblog engine</generator>
        {{#entries}}<item>
            <title>{{Date}}</title>
            <link>{{Guid}}</link>
            <description>{{Body}}</description>
            <pubdate>{{RssDate}}</pubdate>
            <guid>{{Guid}}</guid>
        </item>
        {{/entries}}
    </channel>
</rss>
`

func htmlize(input string) string {
    b := []byte(input)
    output := bytes.NewBufferString("")
    template.HTMLEscape(output, b)
    return output.String()
}

func page(content string) string {
    return fmt.Sprintf(page_str, *title, *title, content)
}

func edit_form(action, id, date, body, button_label string) string {
    t := edit_form_template
    m := make(map[string]interface{})
    m["action"] = action
    m["id"] = id
    if len(date) > 0 {
        m["date"] = date
    } else {
        m["date"] = time.LocalTime().Format(nobodycares.NCTIME)
    }
    m["body"] = body
    m["button_label"] = button_label
    s := mustache.Render(t, m)
    return s
}


func get_root(ctx *web.Context) {
    get_from(ctx, "")
}

func get_from(ctx *web.Context, id string) {
    log.Stderrf("get_root\n")
    p := `
        {{#entries}}` + entry_template + `{{/entries}}
        <br/>
        <!-- get_from id {{id}} -->
        {{#from_id}}
            <a href="/from/{{from_id}}">Older &gt;</a>
        {{/from_id}}
    `
    t := page(p)
    m := make(map[string]interface{})
    entries, _ := nobodycares.LoadRange(id, *max_entries)
    for i, _ := range entries {
        entries[i].Body = htmlize(entries[i].Body)
    }
    m["entries"] = entries
    m["id"] = id
    if len(entries) == *max_entries {
        m["from_id"] = entries[len(entries)-1].Id
    }
    s := mustache.Render(t, m)
    ctx.WriteString(s)
}

func get_post(ctx *web.Context) {
    log.Stderrf("get_post\n")
    ctx.WriteString(page(edit_form("/post", "", "", "", "Post")))
}

func get_edit(ctx *web.Context, id string) {
    log.Stderrf("get_edit %s\n", id)
    if e, err := nobodycares.Load(id); err == nil {
        ctx.WriteString(page(edit_form("/edit", e.Id, e.Date, e.Body, "Edit")))
    } else {
        ctx.WriteString(page("<p>Invalid ID</p>"))
    }
}

func get_delete(ctx *web.Context, id string) {
    log.Stderrf("get_delete %s\n", id)
    if e, err := nobodycares.Load(id); err == nil {
        ctx.WriteString(page(edit_form("/delete", e.Id, e.Date, e.Body, "Really delete")))
    } else {
        ctx.WriteString(page("<p>Invalid ID</p>"))
    }
}

func get_specific_id(ctx *web.Context, id string) {
    log.Stderrf("get_specific_id %s\n", id)
    if e, err := nobodycares.Load(id); err == nil {
        t := entry_template
        m := map[string]interface{}{"Id": e.Id, "Date": e.Date, "Body": e.Body}
        s := mustache.Render(t, m)
        ctx.WriteString(page(s))
    } else {
        ctx.WriteString(page(fmt.Sprintf("<p>Invalid ID</p> <!--%v-->", err)))
    }
}

func nctime_to_rsstime(nctime string) string {
    if t, err := time.Parse(nobodycares.NCTIME, nctime); err == nil {
        return t.Format(time.RFC1123)
    }
    log.Stderrf("nctime_to_rsstime: failed to convert '%s'\n", nctime)
    return nctime
}

func get_rss(ctx *web.Context) {
    log.Stderrf("get_rss\n")
    ctx.SetHeader("Content-Type", "application/rss+xml", false)
    t := rss_template
    m := map[string]interface{}{"title": *title, "url": *url}
    if entries, err := nobodycares.LoadRange("", *max_entries); err == nil {
        type RSSEntry struct {
            nobodycares.Entry
            Guid    string
            RssDate string
        }
        rss_entries := make([]RSSEntry, len(entries))
        for i, _ := range entries {
            re := RSSEntry{entries[i], fmt.Sprintf("%s/%s", *url, entries[i].Id), nctime_to_rsstime(entries[i].Date)}
            rss_entries[i] = re
        }
        m["entries"] = rss_entries
        m["most_recent_date"] = nctime_to_rsstime(entries[0].Date)
        s := mustache.Render(t, m)
        ctx.WriteString(s)
    } else {
        ctx.WriteString(page(fmt.Sprintf("<p>Error generating RSS: %s</p>", err)))
    }
}

func get_css(ctx *web.Context, path string) {
    log.Stderrf("get_css\n")
    ctx.SetHeader("Content-Type", "text/css", false)
    ctx.WriteString(css_str)
}


func post_post(ctx *web.Context) {
    date, date_ok := ctx.Request.Params["form_date"]
    body, body_ok := ctx.Request.Params["form_body"]
    pass, pass_ok := ctx.Request.Params["form_password"]
    if date_ok && body_ok && pass_ok && len(date[0]) > 0 && len(body[0]) > 0 {
        if err := nobodycares.Save(nobodycares.Entry{nobodycares.GenerateID(), date[0], body[0]}, nobodycares.Hash(pass[0])); err == nil {
            ctx.WriteString(page("<p>Post successful</p>"))
        } else {
            ctx.WriteString(page(fmt.Sprintf("<p>Post failed: %v</p>", err)))
        }
    } else {
        ctx.WriteString(page("<p>Invalid form data</p>"))
    }
}

func post_edit(ctx *web.Context) {
    id, id_ok := ctx.Request.Params["form_id"]
    date, date_ok := ctx.Request.Params["form_date"]
    body, body_ok := ctx.Request.Params["form_body"]
    pass, pass_ok := ctx.Request.Params["form_password"]
    if id_ok && date_ok && body_ok && pass_ok && len(id[0]) > 0 && len(date[0]) > 0 && len(body[0]) > 0 {
        if err := nobodycares.Save(nobodycares.Entry{id[0], date[0], body[0]}, nobodycares.Hash(pass[0])); err == nil {
            ctx.WriteString(page("<p>Edit successful</p>"))
        } else {
            ctx.WriteString(page(fmt.Sprintf("<p>Edit failed: %v</p>", err)))
        }
    } else {
        ctx.WriteString(page("<p>Invalid form data</p>"))
    }
}

func post_delete(ctx *web.Context) {
    id, id_ok := ctx.Request.Params["form_id"]
    pass, pass_ok := ctx.Request.Params["form_password"]
    if id_ok && pass_ok && len(id[0]) > 0 {
        if err := nobodycares.Delete(id[0], nobodycares.Hash(pass[0])); err == nil {
            ctx.WriteString(page("<p>Delete successful</p>"))
        } else {
            ctx.WriteString(page(fmt.Sprintf("<p>Delete failed: %v</p>", err)))
        }
    } else {
        ctx.WriteString(page("<p>Invalid form data</p>"))
    }
}

func main() {
    flag.Parse()
    nobodycares.Init(nobodycares.NewCouchStore(*couch_host, strconv.Itoa(*couch_port), *couch_name), *password_hash)
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
