package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type HomeController struct{}

func (c *HomeController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

func (c *HomeController) GetPing() string {
	return "pong"
}

func (c *HomeController) GetTest1() string {
	return "test1"
}

func (c *HomeController) GetHello() interface{} {
	return map[string]string{"message": "Hello Iris!"}
}

func (c *HomeController) BeforeActivation(b mvc.BeforeActivation) {
	anyMiddlewareHere := func(ctx iris.Context) {
		ctx.Application().Logger().Warnf("Inside /custom_path")
		ctx.Next()
	}
	b.Handle(
		"GET",
		"/custom_path",
		"CustomHandlerWithoutFollowingTheNamingGuide", anyMiddlewareHere)

}

func (c *HomeController) CustomHandlerWithoutFollowingTheNamingGuide() string {
	return "hello from the custom handler without following the naming guide"
}
