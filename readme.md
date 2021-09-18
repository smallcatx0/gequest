# gequest

简体中文 | [English](readme.en.md)

一个轻量级的、语义化的、链式操作的 golang http 客户端封装

## 用法

>  要求：
>  - golang >= 1.14

### 安装

```go
go get -u "github.com/smallcatx0/gequest"

import (
    request "github.com/smallcatx0/gequest"
)
```

### 快速开始

```go
package main

import (
	"log"
	"net/url"

	request "github.com/smallcatx0/gequest"
	"github.com/tidwall/gjson"
)

func main() {
	get()
	log.Println("============================")
	post()
	log.Println("============================")
	orderMulti()
}

var EchoServ = "http://postman-echo.com"

// 发起 get 请求
func get() {
	cli := request.New("you-app-name", "Target-Service-Name", 5000)
	param := url.Values{}
	param.Add("name", "lilei")
	cli.SetMethod("get").SetUri(EchoServ).SetPath("/get")
	cli.SetQuery(param).SetHeader("user-agent", "gequest/1.1")
	r, err := cli.Send()
	if err != nil {
		log.Print(err)
	}
	bodyStr, err := r.ToString()
	if err != nil {
		log.Print(err)
	}
	log.Print(bodyStr)
	/* ====output====
	2021/09/07 14:19:10 {"args":{"name":"lilei"},"headers":{"x-forwarded-proto":"http","x-forwarded-port":"80","host":"postman-echo.com","x-amzn-trace-id":"Root=1-6137045e-739a474709ad1f4665e068ea","user-agent":"gequest/1.1","service-name":"you-app-name","target-service-name":"Target-Service-Name","accept-encoding":"gzip"},"url":"http://postman-echo.com/get?name=lilei"}
	*/
}

// 发起 post 请求
func post() {
	cli := request.New("you-app-name", "Target-Service-Name", 5000)
	cli.SetMethod("post").SetUri(EchoServ).SetPath("/post")
	cli.SetJson(map[string]interface{}{
		"name": "hanmeimei",
		"age":  "18",
	})
	r, err := cli.SendRtry(3)
	if err != nil {
		log.Print(err)
	}
	bodyStr, err := r.ToString()
	if err != nil {
		log.Print(err)
	}
	log.Print(bodyStr)
	/* === output ===
	2021/09/07 14:22:11 {"args":{},"data":{"age":"18","name":"hanmeimei"},"files":{},"form":{},"headers":{"x-forwarded-proto":"http","x-forwarded-port":"80","host":"postman-echo.com","x-amzn-trace-id":"Root=1-61370513-431a128438b68b8a0feb6c74","content-length":"31","user-agent":"Go-http-client/1.1","content-type":"application/json","service-name":"you-app-name","target-service-name":"Target-Service-Name","accept-encoding":"gzip"},"json":{"age":"18","name":"hanmeimei"},"url":"http://postman-echo.com/post"}
	*/
}

// 有序并发请求
func orderMulti() {
	requests := make([]*request.Core, 0, 5)
	// 准备请求
	for i := 0; i < 5; i++ {
		requests = append(requests,
			request.New("", "postman-echo.com", 0).
				SetMethod("post").SetPath("/post").
				SetJson(map[string]int{"index": i}),
		)
	}
	// 并发发起请求
	res := request.MultRequest(3, requests...)
	// res 的结果按请求顺序排布
	for i, one := range res {
		if one.Err != nil {
			log.Fatal(one.Err)
		}
		r := res[i].Core.Response()
		resStr, _ := r.ToString()
		// log.Print(resStr)
		log.Print(gjson.Get(resStr, "json.index").Int()) // 0 1 2 3 4
	}
}
```



## 文档

### 创建 `request.Core` 实例

函数签名：

```go
func New(serviceName, targetServiceName string, timeOutMs int) *Core

// 用法
cli := request.New("你的服务名", "目标服务名", 超时时间)
```

当未设置请求域名时 会以 `http://目标服务名` 作为请求域名

超时时间单位为毫秒，设置为0表示不设置超时时间

### 设置请求方法

```go
// 函数签名
func (c *Core) SetMethod(method string) *Core
```

请求方法不区分大小写，这意味着 `post` 你可以写成 `PoSt` `pOsT` `POst` `posT`

### 设置请求地址

函数签名

```go
// 设置请求域
func (c *Core) SetUri(uri string) *Core
// 设置请求路径
func (c *Core) SetPath(path string) *Core
```

最终的请求地址=请求域(uri)+请求路径(path)，若请求域为空则会以[创建实例](#创建  实例)时的第二个参数作为请求域

### 设置请求参数

函数签名：

```go
// 设置Query参数
func (c *Core) SetQuery(param url.Values) *Core 
// 设置json参数 会将param json序列化并在请求头中添加 "Content-Type":"application/json"
func (c *Core) SetJson(param interface{}) *Core 
// 设置请求体 二进制数据
func (c *Core) SetBody(data []byte) *Core 
// 设置请求体字符串 会在请求头中添加 "Content-Type":"application/json"
func (c *Core) SetBodyText(body string) 
// 设置请求头 重复设置同一个key会覆盖
func (c *Core) SetHeader(k string, v string) *Core 
func (c *Core) SetHeaders(headers map[string]string) *Core 
func (c *Core) AddHeaders(headers map[string]string) *Core 
```

### 发起请求

函数签名：

```go
func (c *Core) Send() (r *Response, err error)
//发起请求并在异常时重试
func (c *Core) SendRtry(times int) (r *Response, err error)
```

`Send`与`SendRtry` 方法除了会返回 `*request.Response` 也会将结果保存至`c.request`中。通过[获取响应](#获取响应)方法可再次获得

### 获取响应

函数签名：

```go
func (c *Core) Response() *Response 
func (c *Core) ResponseRaw() *http.Response 

type Response struct {
	*http.Response
}
func (r *Response) ReadAll() ([]byte, error) 
func (r *Response) ToString() (string, error)
```

`request.Respones` 为标准库`http.Response`的继承，所以 `http.Response` 能用的方法 `request.Respones` 也能用

除了调用[发起请求](#发起请求)后获取到的返回值，还可以使用`c.Response()`方法获取到保存在实例中的响应

**但是：请注意 http.Response.Body 是一个流 读一次就没了**除非你再写回去参考`func (r *Response) watchAll() (body []byte, err error) `方法

### 日志

函数签名：

```go
// 开启debug才会写日志 
func (c *Core) Debug(debug bool) *Core 

// 兼容标准库log的接口
var _ Logger = log.New(os.Stdout, "", 0)
type Logger interface {
	Print(v ...interface{})
}

// 接管日志 
func (c *Core) SetLoger(loger Logger) *Core 
```

开启debug才会写日志。每次请求都会记录请求参数响应数据。默认日志写在控制台

`var ConsoleLog Logger = log.New(os.Stderr, "[request] ", log.Ldate|log.Ltime|log.Lshortfile)`

### 错误

函数签名：

```go
// 返回链式过程中发生的任何错误
func (c *Core) Errors() []error
```

### 并发请求

函数签名:

```go
type Resp struct {
	Index int
	Core  *Core
	Err   error
}
func MultRequest(retime int, requests ...*Core) []Resp
```

传入多个准备好的请求实例，返回的结果会以请求顺序排布。



## 单元测试

```bash
go test -timeout 30s -coverprofile=/tmp/gequest-cover github.com/smallcatx0/gequest -v

...
PASS
coverage: 86.1% of statements
ok      github.com/smallcatx0/gequest   6.951s  coverage: 86.1% of statements
```

