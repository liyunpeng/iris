package viewmodels
//应该有视图模型，客户端将能够看到的结构例：
import (
"github.com/kataras/iris/_examples/mvc/login/datamodels"
"github.com/kataras/iris/context"
)

type User struct {
	datamodels.User
}
func (m User) IsValid() bool {
	/*做一些检查，如果有效则返回true ... */
	return m.ID > 0
}

//Iris能够将任何自定义数据结构转换为HTTP响应调度程序，
//从理论上讲，如果真的有必要，可以使用以下内容;

//Dispatch实现`kataras/iris/mvc＃Result`接口。
//将`User` 作为受控的http响应发送。
//如果其ID为零或更小，则返回404未找到错误
//否则返回其json表示，
//（就像控制器的函数默认为自定义类型一样）。
//不要过度，应用程序的逻辑不应该在这里。可以在这里添加简单的检查验证
//这只是一个展示
//想象一下设计更大的应用程序时此功能將很有帮助。
//调用控制器方法返回值的函数
//是`User` 的类型。
func (m User) Dispatch(ctx context.Context) {
	if !m.IsValid() {
		ctx.NotFound()
		return
	}
	ctx.JSON(m, context.JSON{Indent: " "})
}
