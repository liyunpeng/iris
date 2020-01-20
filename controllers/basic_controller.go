package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

type LoggerService interface {
	Log(string)
}

type prefixedLogger struct {
	prefix string
}

func (s *prefixedLogger) Log(msg string) {
	fmt.Printf("11111111111111 %s: %s\n", s.prefix, msg)
}

type BasicController struct {
	Logger  LoggerService
	Session *sessions.Session
}

func (c *BasicController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/custom", "Custom")
}

func (c *BasicController) AfterActivation(a mvc.AfterActivation) {
	if a.Singleton() {
		panic("BasicController should be stateless,a request-scoped,we have a 'Session' which depends on the context.")
	}
}

func (c *BasicController) Get() string {
	count := c.Session.Increment("count", 1)
	body := fmt.Sprintf("Hello from BasicController\nTotal visits from you: %d", count)
	c.Logger.Log(body)
	return body
}

func (c *BasicController) Custom() string {
	return "custom"
}


type basicSubController struct {
	Session *sessions.Session
}

func (c *basicSubController) Get() string {
	count := c.Session.GetIntDefault("count", 1)
	return fmt.Sprintf("Hello from basicSubController.\nRead-only visits count: %d", count)
}

func BasicMVC(app *mvc.Application) {
	//当然，你可以在MVC应用程序中使用普通的中间件。
	app.Router.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("222222222222222 Path: %s", ctx.Path())
		ctx.Next()
	})
	//把依赖注入，controller(s)绑定
	//可以是一个接受iris.Context并返回单个值的函数（动态绑定）
	//或静态结构值（service）。
	app.Register(
		sessions.New(sessions.Config{}).Start,
		&prefixedLogger{prefix: "1111111111111 DEV"},
	)
	app.Handle(new(BasicController))
	//所有依赖项被绑定在父 *mvc.Application
	//被克隆到这个新子身上，父的也可以访问同一个会话。
	// GET: http://localhost:8080/basic/sub
	app.Party("/sub").Handle(new(basicSubController))
}
