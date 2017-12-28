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

Configuration `configs/app.yml`

```
Other: {
  EsUrl: "", # elasticsearch url, e.g. http://localhost:9200 
  LogPath: "", 
  DbPath: "", # boltdb db filepath, e.g. /data/db/app.db
  MailHost: "", # email host, e.g. smtp.163.com
  MailUser: "", # email username
  MailPasword: "" # email password
  MailHtmlTplUrl: "" # email content template, e.g. /data/template.html
}
```

Run app

```
$ go run main.go
```

Then open `http://localhost:8080` in your browser.
