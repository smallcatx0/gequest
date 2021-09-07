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
	res := request.MultRequest(requests...)
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
