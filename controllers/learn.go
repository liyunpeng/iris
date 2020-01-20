package controllers

import (
	"context"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"html/template"
	elasticesearch1 "my.web/elasticesearch"
	gopoll "my.web/gopool"
	kafka1 "my.web/kafka"
	"my.web/lib"
	"my.web/zookeeper"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type LearnController struct{}

var learnView = mvc.View{
	Name: "learn/learn.html",
	Data: map[string]interface{}{
		"Title":     "Hello Page",
		//"MyMessage": "Welcome to my awesome website",
	},
}

func (c *LearnController) Get() mvc.Result {
	return learnView
}



func (m *LearnController) BeforeActivation(b mvc.BeforeActivation) {
	// b.Dependencies().Add/Remove
	// b.Router().Use/UseGlobal/Done // and any standard API call you already know

	// 1-> Method
	// 2-> Path
	// 3-> The controller's function name to be parsed as handler
	// 4-> Any handlers that should run before the MyCustomHandler
	//b.Handle("GET", "/something/{id:long}", "MyCustomHandler", anyMiddleware...)
	//b.Handle("GET","/users/info","QueryInfo")

	anyMiddlewareHere := func(ctx iris.Context) {
		ctx.Application().Logger().Warnf("Inside /custom_path")
		ctx.Next()
	}

	b.Handle("GET", "/path1",
		"CustomHandlerWithoutFollowingTheNamingGuide", anyMiddlewareHere)

	b.Handle("GET", "/elasticsearch",
		"F123", anyMiddlewareHere)

}

func (m *LearnController)  F123() string{
	elasticesearch1.Es()
	return  "abc"
}

func (c *LearnController) CustomHandlerWithoutFollowingTheNamingGuide() string {
	return "hello from the custom handler without following the naming guide"
}


func LearnIndex(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("view/html/learn.gtpl")

	w.Header().Set("Content-Type", "text/html")

	t.Execute(w, "")
}

func OnLearnAjax(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	method := r.Form.Get("method")

	fmt.Println("method=", method)

	if strings.Compare(method, "contextLearn") == 0 {
		//contextMain()
		lib.Context()
	} else if strings.Compare(method, "certButton") == 0 {
		fmt.Print("compare ", method)
		lib.GenerateCerteKeyPem()
	} else if strings.Compare(method, "arrslicerange") == 0 {
		fmt.Print("ArrayRange: ")
		lib.ArrayRange()

		fmt.Print("\n SliceRange: ")
		lib.SliceRange()

		fmt.Print("\n Slice1Range: ")
		lib.Slice1Range()

		fmt.Print("\n PointerRange: ")
		lib.PointerRange()

		fmt.Print("\n MapRange: ")
		lib.MapRange()

		fmt.Print("\n MapCounterRange: ")
		lib.MapCounterRange()
	} else if strings.Compare(method, "gopool") == 0 {
		gopoll.GopollMain()
		fmt.Println("真正可以开启的线程数=", runtime.GOMAXPROCS(0))
	} else if strings.Compare(method, "redis") == 0 {
		lib.TestAll()
	} else if strings.Compare(method, "unsafe") == 0 {
		lib.Unsafemain()
		lib.Unsafe2()
	} else if strings.Compare(method, "updateversion") == 0 {
		lib.UpdateVersion()
	} else if strings.Compare(method, "syncpool") == 0 {
		//lib.SyncPool2()
		//lib.PoolMain2()
		lib.PoolMain3()
	} else if strings.Compare(method, "scanport") == 0 {
		//lib.Scanner1("58.96.172.22", "1-520", "512")
		//lib.Scanner1("127.0.0.1", "8080", "2")  //实验结果开放
		//lib.Scanner1("127.0.0.1", "80", "1")
		lib.Scanner1("127.0.0.1", "80,3306", "1")
	} else if strings.Compare(method, "publish") == 0 {
		lib.Publish1(2)
	} else if strings.Compare(method, "consume") == 0 {
		lib.Consume(2)
	} else if strings.Compare(method, "kafkaPublish") == 0 {
		kafka1.Kafkaserver()
	} else if strings.Compare(method, "kafkaConsume") == 0 {
		kafka1.KafkaClient()
	} else if strings.Compare(method, "kafkaAsyncPublish") == 0 {
		kafka1.SaramaAsyncProducer()
	} else if strings.Compare(method, "zookeeper") == 0 {
		zookeeper.Zookeepermain()
	} else if strings.Compare(method, "zookeeper1") == 0 {
		zookeeper.Zookeepermain1()
	} else if strings.Compare(method, "elastisearch") == 0 {
		elasticesearch1.Es()
	}
}

func contextMain() {
	ctx, _ := context.WithTimeout(context.Background(), (10 * time.Second))
	go routineA(ctx)
}

func routineA(ctx context.Context) {
	ctxA, _ := context.WithTimeout(ctx, (5 * time.Second))
	ch := make(chan int)
	go routineB(ctxA, ch)

	select {
	case <-ctx.Done():
		fmt.Println("routineA Done")
		return
	case i := <-ch:
		fmt.Println(i)
	}
}

func routineB(ctx context.Context, ch chan int) {
	//模拟读取数据
	sumCh := make(chan int)
	go func(sumCh chan int) {
		sum := 10
		fmt.Println("routineB 又起的routine ")
		time.Sleep(10 * time.Second)
		fmt.Println("routineB 又起的routine 10秒后")
		sumCh <- sum
	}(sumCh)

	select {
	case <-ctx.Done():
		fmt.Println("routineB Done")
		<-sumCh
		return
	//case ch  <- <-sumCh: 注意这样会导致资源泄露
	case i := <-sumCh:
		fmt.Println("send", i)
		ch <- i
	}
}
