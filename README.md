ELK Alarm components

Monitor application logs and alert by using ES APIs.

But I think the right way to do this is [here (https://github.com/grafana/grafana/issues/5893)](https://github.com/grafana/grafana/issues/5893).

So I just use *go* to do something.

xzlan depends on these

* iris
* github.com/boltdb/bolt
* gopkg.in/olivere/elastic.v5
* layui

## Getting Started

Configuration `./conf/app.yml`

```
Other: {
  ServerProt: "" # default 8001
  EsUrl: "", # elasticsearch url, e.g. http://localhost:9200
  EsIndex: "", # elasticsearch logstash index, e.g. "logstash-", no date
  LogPath: "", 
  DbPath: "", # boltdb db filepath default ./alert.db, e.g. /data/db/app.db
  MailHost: "", # email host, e.g. smtp.163.com
  MailUser: "", # email username
  MailPassword: "" # email password
  MailHtmlTplUrl: "" # email content template, e.g. /data/template.html
}
```

Run app

```
$ go get -u github.com/jteeuwen/go-bindata/...
$ go-bindata static/... views/...
$ go build
$ ./xzlan -config=./conf/app.yml
```

Then open `http://localhost:8001` in your browser.
