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
	"fmt"
	"strings"
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func main() {
	args := os.Args[1:]
	custConf := config(args)
	
	conf := iris.YAML("./config/app.yml")
	if custConf.DbPath == "" {
		custConf.EsUrl = conf.Other["EsUrl"].(string)
		custConf.LogPath = conf.Other["LogPath"].(string)
		custConf.DbPath = conf.Other["DbPath"].(string)
		custConf.MailHost = conf.Other["MailHost"].(string)
		custConf.MailUser = conf.Other["MailUser"].(string)
		custConf.MailPasword = conf.Other["MailPasword"].(string)
		custConf.MailHtmlTplUrl = conf.Other["MailHtmlTplUrl"].(string)
	}
	fmt.Println(custConf.String())

	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{Status:true, IP:false, Method:true, Path:true}))

	f := newLogFile(custConf.LogPath)
	defer f.Close()

	app.Logger().AddOutput(newLogFile(custConf.LogPath))

	app.StaticWeb("/static", "./static")
	app.RegisterView(iris.HTML("./views", ".html").Layout("layout/layout.html").
		Delims("<<", ">>"))

	// Open DB
	daoDb, dbErr := dao.NewDao(custConf.DbPath)
	if dbErr != nil {
		app.Logger().Warn("open DB error: " + dbErr.Error())
	}
	defer daoDb.Db.Close()

	alertMail := mail.NewMail(custConf.MailUser, custConf.MailPasword, custConf.MailHost, custConf.MailHtmlTplUrl)
	apiDao := dao.NewApiDao(daoDb)
	ruleDao := dao.NewRuleDao(daoDb)
	noteDao := dao.NewNoteDao(daoDb)
	globalmailDao := dao.NewGlobalMailDao(daoDb)
	apiAlert := alert.NewAlert(apiDao, ruleDao, noteDao, globalmailDao, alertMail, conf.Other["EsUrl"].(string))
	app.Controller("/apis", new(controller.ApiController), apiDao, ruleDao, apiAlert)
	app.Controller("/rule", new(controller.RuleController), ruleDao)
	app.Controller("/notes", new(controller.NoteController), noteDao)
	app.Controller("/globalmails", new(controller.GlobalMailController), globalmailDao)

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

	app.Configure(iris.WithConfiguration(conf))
	if err := app.Run(iris.Addr(":8080"), iris.WithoutBanner); err != nil {
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
	EsUrl string `yaml:"EsUrl"`
	LogPath string `yaml:"LogPath"`
	DbPath string `yaml:"DbPath"`
	MailHost string `yaml:"MailHost"`
	MailUser string `yaml:"MailUser"`
	MailPasword string `yaml:"MailPasword"`
	MailHtmlTplUrl string `yaml:"MailHtmlTplUrl"`
}

func (c *customizeConfig) String() string {
	str := fmt.Sprintf("print configuration info: \n")
	str = str + fmt.Sprintf("EsUrl: %s \n", c.EsUrl)
	str = str + fmt.Sprintf("LogPath: %s \n", c.LogPath)
	str = str + fmt.Sprintf("DbPath: %s \n", c.DbPath)
	str = str + fmt.Sprintf("MailHtmlTplUrl: %s \n", c.MailHtmlTplUrl)
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