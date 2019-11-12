package services

import (
	//"fmt"
	"github.com/kataras/iris"
)

type Form struct {
	Something string
}

func StartServices(app *iris.Application) {
	// 静态文件
	app.HandleDir("/static", "./static")
	// 使用模板加载视图
	app.RegisterView(iris.HTML("./templates", ".html").Reload(true))
	
	app.Get("/hello", hello)
	app.Get("/unknown", unknown)
	app.Post("/info", getInfo)
	app.OnErrorCode(iris.StatusNotFound, notFound)
}

func hello(ctx iris.Context) {
	ctx.View("hello.html")
}

func notFound(ctx iris.Context) {
    ctx.JSON(iris.Map{
		"error": "404 Not Found",
	})
}

func unknown(ctx iris.Context) {
	ctx.StatusCode(501)
	ctx.JSON(iris.Map{
		"error": "501 ~开发中~",
	})
}

func getInfo(ctx iris.Context) {
	form := Form{}
	err := ctx.ReadForm(&form)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"error": "Something Wrong",
		})
	} else {
		ctx.ViewData("something", form.Something)
		ctx.View("info.html")
	}
}