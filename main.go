package main

import (
	"github.com/kataras/iris"
	"time"
	"os"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/middleware/logger"
	"xzlan/dao"
	"xzlan/controller"
	"xzlan/alert"
	"xzlan/mail"
)

func todayFilename() string {
	today := time.Now().Format("2006-01-02")
	return today + ".log"
}

func newLogFile(path string) *os.File {
	filename := "/Users/tiger/project/logs/xzlan/" + todayFilename()
	// open an output file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func main() {
	conf := iris.YAML("./configs/app.yml")

	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{Status:true, IP:false, Method:true, Path:true}))

	app.Logger().Info("LogPath: ", conf.Other["LogPath"])
	app.Logger().Info("DbPath: ", conf.Other["DbPath"])
	//f := newLogFile(conf.Other["LogPath"])
	//defer f.Close()
	//app.Logger().AddOutput(newLogFile())

	app.StaticWeb("/static", "./static")
	app.RegisterView(iris.HTML("./views", ".html").Layout("layout/layout.html").Delims("<<", ">>"))

	// Open DB
	daPath := conf.Other["DbPath"]
	daoDb, dbErr := dao.NewDao(daPath.(string))
	if dbErr != nil {
		app.Logger().Warn("open DB error: " + dbErr.Error())
	}
	defer daoDb.Db.Close()

	alertMail := mail.NewMail(conf.Other["MailUser"].(string), conf.Other["MailPasword"].(string),
		conf.Other["MailHost"].(string))
	apiDao := dao.NewApiDao(daoDb)
	ruleDao := dao.NewRuleDao(daoDb)
	apiAlert := alert.NewAlert(apiDao, ruleDao, alertMail, conf.Other["EsUrl"].(string))
	app.Controller("/apis", new(controller.ApiController), apiDao, apiAlert)
	app.Controller("/rule", new(controller.RuleController), ruleDao)

	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	app.Get("/api", func(ctx iris.Context) {
		ctx.View("metric/apis.html")
	})

	app.Configure(iris.WithConfiguration(conf))
	if err := app.Run(iris.Addr(":8080"), iris.WithoutBanner); err != nil {
		if err != iris.ErrServerClosed {
			app.Logger().Warn("Shutdown with error: " + err.Error())
		}
	}
}