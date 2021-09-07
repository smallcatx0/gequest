# xthk-request

## 用法

快速开始

```go
package main

import (
	"fmt"
	"net/url"

	request "github.com/smallcatx0/gequest"
)

func main() {

	client := request.New("self", "Target", 0)
	param := url.Values{
		"name": []string{"xiaoming"}, "age": []string{"20"},
	}
	res, err := client.SetUri("http://postman-echo.com").
		SetMethod("get").
		SetPath("/get").
		SetQuery(param).
		Send()
	if err != nil {
		fmt.Println(err)
	} else {
		resJson, _ := res.ToString()
		fmt.Println(res.Status, resJson)
        // {"args":{"age":"20","name":"xiaoming"},"headers":{"x-forwarded-proto":"http","x-forwarded-port":"80","host":"postman-echo.com","x-amzn-trace-id":"Root=1-6093d836-5fec68e95cbd47fd14d8f301","user-agent":"Go-http-client/1.1","service-name":"self","target-service-name":"Target","accept-encoding":"gzip"},"url":"http://postman-echo.com/get?age=20&name=xiaoming"}
	}
}
```



#### GET query

```go
// http://postman-echo.com/get?age=20&name=xiaoming
request.New("self", "Target", 0).
    SetUri("http://postman-echo.com").
    SetMethod("get").
    SetPath("/get").
    SetQuery(url.Values{
		"name": []string{"xiaoming"}, "age": []string{"20"},
	}).
    Send()
```

#### POST json

```go
request.New("self", "Target", 0).
    SetUri("http://postman-echo.com").
    SetMethod("PosT").
    SetPath("/post").
    SetJson(map[string]string{"key": "value"}).
    Send()

// {"args":{},"data":{"key":"value"},"files":{},"form":{},"headers":{"x-forwarded-proto":"http","x-forwarded-port":"80","host":"postman-echo.com","x-amzn-trace-id":"Root=1-6093d9c8-4264d7bb6027b4d673073aa2","content-length":"15","user-agent":"Go-http-client/1.1","content-type":"application/json","service-name":"self-Service-NAME","target-service-name":"Target-Service-Name","accept-encoding":"gzip"},"json":{"key":"value"},"url":"http://postman-echo.com/post"}
```



#### POST raw

```go
request.New("self-Service-NAME", "Target-Service-Name", 0).
    SetUri("http://postman-echo.com").
    SetMethod("post").
    SetPath("/post").
    SetHeaders(map[string]string{"x-test-name": "tttta"}).
    SetBodyText("hello world")

// {"args":{},"data":"hello world","files":{},"form":{},"headers":{"x-forwarded-proto":"http","x-forwarded-port":"80","host":"postman-echo.com","x-amzn-trace-id":"Root=1-6093da71-19575f164785c9ec024ca3f4","content-length":"11","user-agent":"Go-http-client/1.1","content-type":"text/plain","x-test-name":"tttta","accept-encoding":"gzip"},"json":null,"url":"http://postman-echo.com/post"}

```



