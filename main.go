package main

import (
	"github.com/kataras/iris"
	"time"
	"os"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/middleware/logger"
	"xzlan/dao"
	"xzlan/controller"
)

func todayFilename() string {
	today := time.Now().Format("2006-01-02")
	return today + ".log"
}

func newLogFile() *os.File {
	filename := "/Users/tiger/project/logs/xzlan/" + todayFilename()
	// open an output file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

func main() {
	//f := newLogFile()
	//defer f.Close()

	conf := iris.YAML("./configs/app.yml")

	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{Status:true, IP:false, Method:true, Path:true}))
	//app.Logger().AddOutput(newLogFile())

	app.Logger().Info("LogPath: ", conf.Other["LogPath"])

	app.StaticWeb("/static", "./static")
	app.RegisterView(iris.HTML("./views", ".html").Layout("layout/layout.html"))

	// Open DB
	metricDao, dbErr := dao.NewDao("/Users/tiger/project/logs/go/xzlan.db")
	if dbErr != nil {
		app.Logger().Warn("open DB error: " + dbErr.Error())
	}
	defer metricDao.Db.Close()

	app.Controller("/apis", new(controller.MetricController), metricDao)

	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		//ctx.HTML("<h1>Welcome</h1>")
		ctx.View("index.html")
	})

	app.Get("/api", func(ctx iris.Context) {
		ctx.View("metric/apis.html")
	})

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})

	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		name := ctx.URLParam("name")
		ctx.Application().Logger().Info(name)
		ctx.JSON(iris.Map{"message": "Hello Iris!" + " " + name})
	})

	app.Configure(iris.WithConfiguration(conf))
	if err := app.Run(iris.Addr(":8080"), iris.WithoutBanner); err != nil {
		if err != iris.ErrServerClosed {
			app.Logger().Warn("Shutdown with error: " + err.Error())
		}
	}
}