package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"xzlan/alert"
	"xzlan/controller"
	"xzlan/dao"
	"xzlan/mail"
)

func main() {
	args := os.Args[1:]
	custConf := config(args)

	if custConf.DbPath == "" {
		custConf.DbPath = "./alert.db"
	}
	if custConf.ServerPort == "" {
		custConf.ServerPort = "8001"
	}
	fmt.Println(custConf.String())

	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{Status:true, IP:false, Method:true, Path:true}))

	f := newLogFile(custConf.LogPath)
	defer f.Close()

	// Open DB
	daoDb, dbErr := dao.NewDao(custConf.DbPath)
	if dbErr != nil {
		app.Logger().Warn("open DB error: " + dbErr.Error())
	}
	defer daoDb.Db.Close()

	app.Logger().AddOutput(newLogFile(custConf.LogPath))

	app.HandleDir("/static", "./static", iris.DirOptions{
		Asset:      Asset,
		AssetInfo:  AssetInfo,
		AssetNames: AssetNames,
		ShowList:   true,
	})
	app.RegisterView(iris.HTML("./views", ".html").Layout("layout/layout.html").
		Delims("<<", ">>").Binary(Asset, AssetNames))

	//app.StaticWeb("/static", "./static")
	//app.RegisterView(iris.HTML("./views", ".html").Layout("layout/layout.html").
	//	Delims("<<", ">>"))

	alertMail := mail.NewMail(custConf.MailUser, custConf.MailPassword, custConf.MailHost, custConf.MailHTMLTplURL)
	apiDao := dao.NewAPIDao(daoDb)
	ruleDao := dao.NewRuleDao(daoDb)
	noteDao := dao.NewNoteDao(daoDb)
	globalmailDao := dao.NewGlobalMailDao(daoDb)
	apiAlert := alert.NewAlert(apiDao, ruleDao, noteDao, globalmailDao, alertMail, custConf.EsURL, custConf.EsIndex)
	mvc.New(app.Party("/apis")).Handle(new(controller.APIController)).Register(apiDao, ruleDao, apiAlert)
	mvc.New(app.Party("/rule")).Handle(new(controller.RuleController)).Register(ruleDao, apiDao)
	mvc.New(app.Party("/notes")).Handle(new(controller.NoteController)).Register(noteDao)
	mvc.New(app.Party("/globalmails")).Handle(new(controller.GlobalMailController)).Register(globalmailDao)

	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.View("index.html")
	})
	app.Get("/api", func(ctx iris.Context) {
		ctx.View("api/apis.html")
	})
	app.Get("/note", func(ctx iris.Context) {
		ctx.View("note/notes.html")
	})
	app.Get("/globalmail", func(ctx iris.Context) {
		ctx.View("globalmail/globalmails.html")
	})

	if err := app.Run(iris.Addr(":" + custConf.ServerPort), iris.WithoutBanner); err != nil {
		if err != iris.ErrServerClosed {
			app.Logger().Warn("Shutdown with error: " + err.Error())
		}
	}
}

func config(args []string) customizeConfig {
	if len(args) == 0 {
		fmt.Println("[warn] no profile specified, e.g. -config=/data/app.yml")
	}
	c := customizeConfig{}
	for _, v := range args {
		conf := strings.Split(v, "=")
		if strings.Index(conf[0], "config") == -1 {
			fmt.Println("no profile specified, e.g. -config=/data/app.yml")
			os.Exit(-1)
		}
		if _, err := os.Stat(conf[1]); os.IsNotExist(err) {
			fmt.Println("configuration file '" + conf[1] + "' does not exist")
			os.Exit(-1)
		}
		fileExt := filepath.Ext(conf[1])
		if fileExt != ".yml" {
			fmt.Println("configuration file '" + conf[1] + "' invalid, please use extension .yml")
			os.Exit(-1)
		}
		data, err := ioutil.ReadFile(conf[1])
		if err != nil {
			fmt.Println("configuration file '" + conf[1] + "' invalid")
			os.Exit(-1)
		}
		if err := yaml.Unmarshal(data, &c); err != nil {
			fmt.Println("configuration file '" + conf[1] + "' invalid")
			os.Exit(-1)
		}
		break
	}
	return c
}

type customizeConfig struct {
	ServerPort string `yaml:"ServerPort"`
	EsURL string `yaml:"EsURL"`
	EsIndex string `yaml:"EsIndex"`
	LogPath string `yaml:"LogPath"`
	DbPath string `yaml:"DbPath"`
	MailHost string `yaml:"MailHost"`
	MailUser string `yaml:"MailUser"`
	MailPassword string `yaml:"MailPassword"`
	MailHTMLTplURL string `yaml:"MailHTMLTplURL"`
}

func (c *customizeConfig) String() string {
	str := fmt.Sprintf("print configuration info: \n")
	str = str + fmt.Sprintf("ServerPort: %s \n", c.ServerPort)
	str = str + fmt.Sprintf("EsURL: %s \n", c.EsURL)
	str = str + fmt.Sprintf("LogPath: %s \n", c.LogPath)
	str = str + fmt.Sprintf("DbPath: %s \n", c.DbPath)
	str = str + fmt.Sprintf("MailHTMLTplURL: %s \n", c.MailHTMLTplURL)
	str = str + fmt.Sprintf("load configuration done \n")
	return str
}

func todayFilename() string {
	today := time.Now().Format("2006-01-02")
	return today + ".log"
}

func newLogFile(path string) *os.File {
	filename := path + todayFilename()
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}