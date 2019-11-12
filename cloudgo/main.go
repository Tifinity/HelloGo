package main

import (
	"github.com/Tifinity/cloudgo/services"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	// 日志等级
	app.Logger().SetLevel("debug")
	// run
	services.StartServices(app)
	// 程序内配置服务端
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.Configuration{
	    DisableInterruptHandler:           false,
	    DisablePathCorrection:             false,
	    EnablePathEscape:                  false,
	    FireMethodNotAllowed:              false,
	    DisableBodyConsumptionOnUnmarshal: false,
	    DisableAutoFireStatusCode:         false,
	    TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
	    Charset:                           "UTF-8",
	}))
}