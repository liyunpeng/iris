package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"my.web/web/controllers"
	"my.web/web/datasource"
	"my.web/web/middleware"
	"my.web/web/repositories"
	"my.web/web/services"
	datamodels2 "my.web/web/datamodels"
)

//注意mvc.Application，它不是iris.Application。
func movies(app *mvc.Application) {
	//添加基本身份验证（admin：password）中间件
	//用于基于/movies的请求。
	app.Router.Use(middleware.BasicAuth)
	//使用数据源中的一些（内存）数据创建我们的电影资源库。
	repo := repositories.NewMovieRepository(datasource.Movies)
	//创建我们的电影服务，我们将它绑定到电影应用程序的依赖项中
	movieService := services.NewMovieService(repo)
	app.Register(movieService)
	//为我们的电影控制器服务
	//请注意，您可以为多个控制器提供服务
	//你也可以使用`movies.Party（relativePath）`或`movies.Clone（app.Party（...））创建子mvc应用程序
	// 如果你想。
	app.Handle(new(controllers.MovieController))
}

func WebMain() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	app.RegisterView(iris.HTML("./src/my.web/web/views/", ".html"))

	/*
		添加两个内置处理程序
		可以从任何与http相关的panics中恢复并将请求记录到终端。
	*/

	app.Use(recover.New())
	app.Use(logger.New())

	// 根路由路径 "/"
	mvc.New(app).Handle(new(controllers.HomeController))

	mvc.Configure(app.Party("/basic"), controllers.BasicMVC)

	mvc.Configure(app.Party("/movies"), movies)

	mvc.New(app.Party("/learn")).Handle(new(controllers.LearnController))

	m := make (map[int64] datamodels2.User, 2)

	//s := sessions.New(sessions.Config{})
	//s.Start
	c := controllers.UserController{
		Session: &sessions.Session{},
		Service: services.NewUserService(repositories.NewUserRepository(m)),
	}

	mvc.New(app.Party("/user")).Handle(c)

	mvc.New(app.Party("/hello")).Handle(new(controllers.HelloController))


	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
