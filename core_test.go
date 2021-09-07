package request_test

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	request "github.com/smallcatx0/gequest"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

var echoServ = "http://postman-echo.com"

func TestHostNt(t *testing.T) {
	assert := assert.New(t)
	cli := request.New("self", "Target", 0).
		SetUri("http://not.host").
		SetMethod("Get")
	_, err := cli.Send()
	assert.Error(err)
}

func TestTimeOut(t *testing.T) {
	assert := assert.New(t)
	cli := request.New("self", "target", 1)
	cli.SetUri("http://www.baidu.com").SetMethod("get")

	_, err := cli.Send()
	assert.Error(err)
}

func TestAutoUri(t *testing.T) {
	assert := assert.New(t)
	cli := request.New("self", "postman-echo.com", 0).
		SetMethod("gEt").
		SetPath("/get")
	res, err := cli.Send()
	assert.NoError(err)
	if err != nil {
		return
	}
	resJson, _ := res.ToString()

	assert.Equal("http://postman-echo.com/get", gjson.Get(resJson, "url").String())
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	param := url.Values{
		"name": []string{"ming"}, "age": []string{"18"},
	}
	cli := request.New("self-Service-NAME", "Target-Service-Name", 0)
	res, err := cli.SetUri(echoServ).
		SetPath("/get").
		SetMethod("get").
		SetQuery(param).
		Debug(true).
		Send()
	assert.NoError(err)
	if err != nil {
		return
	}

	resJson, _ := res.ToString()

	assert.Equal(200, res.StatusCode)
	expectUrl := echoServ + "/get" + "?" + param.Encode()
	assert.Equal(expectUrl, gjson.Get(resJson, "url").String())
	assert.Equal("self-service-name", strings.ToLower(gjson.Get(resJson, "headers.service-name").String()))
	assert.Equal("target-service-name", strings.ToLower(gjson.Get(resJson, "headers.target-service-name").String()))
}

func TestPostRaw(t *testing.T) {
	assert := assert.New(t)
	cli := request.New("self-Service-NAME", "Target-Service-Name", 0).
		SetLoger(request.NewFileLogger("./log.log")).
		SetUri(echoServ).
		SetMethod("post").
		SetPath("/post").
		SetHeaders(map[string]string{"x-test-name": "tttta"}).
		Debug(true).
		SetBodyText("hello world")
	res, err := cli.Send()
	assert.NoError(err)
	if err != nil {
		return
	}
	resJson, _ := res.ToString()

	assert.Equal(200, res.StatusCode)
	assert.Equal(200, cli.Response().StatusCode)
	assert.Equal(200, cli.ResponseRaw().StatusCode)
	assert.Equal("hello world", gjson.Get(resJson, "data").String())
	assert.NotEqual("self-service-name", strings.ToLower(gjson.Get(resJson, "headers.service-name").String()))
	assert.NotEqual("target-service-name", strings.ToLower(gjson.Get(resJson, "headers.target-service-name").String()))
}

func TestPostJosn(t *testing.T) {
	assert := assert.New(t)
	param := map[string]string{"key": "value"}
	cli := request.New("self-Service-NAME", "Target-Service-Name", 0).
		SetUri(echoServ).
		SetMethod("POST").
		SetPath("/post").
		Debug(true).
		SetJson(param)
	res, err := cli.Send()
	assert.NoError(err)
	if err != nil {
		return
	}
	resJson, _ := res.ToString()

	assert.Equal(200, res.StatusCode)
	assert.Equal("self-service-name", strings.ToLower(gjson.Get(resJson, "headers.service-name").String()))
	assert.Equal("target-service-name", strings.ToLower(gjson.Get(resJson, "headers.target-service-name").String()))
	paramJson, _ := json.Marshal(param)
	assert.JSONEq(string(paramJson), gjson.Get(resJson, "json").String())
}

func TestWatchResp(t *testing.T) {
	assert := assert.New(t)
	cli := request.New("self-service-name", "target-service-name", 0).
		SetUri(echoServ).
		SetMethod("get").
		SetPath("/get")
	res, err := cli.Send()
	assert.NoError(err)
	_ = cli.String()
	// log.Print(cli.String())
	resJson, _ := res.ToString()
	assert.NotEmpty(resJson)
}

func TestRtry(t *testing.T) {
	assert := assert.New(t)
	cli := request.New("self", "target", 1)
	cli.SetUri("http://www.baidu.com").SetMethod("get")
	_, err := cli.SendRtry(10)
	assert.Error(err)
	_, err = cli.SetUri("http://aa.bbb.com/aaa").SetMethod("post").SendRtry(10)
	assert.Error(err)
}

func TestClearn(t *testing.T) {
	ass := assert.New(t)
	cli := request.New("self", "target", 5000)
	cli.SetUri(echoServ).
		SetPath("/get").
		SetQuery(url.Values{"name": []string{"kuiiii"}})
	cli.Clear()
	cli.SetQuery(url.Values{"age": []string{"18"}})
	errs := cli.Errors()
	ass.Equal(0, len(errs))
	res, err := cli.Send()
	ass.NoError(err)
	if err != nil {
		return
	}
	resJson, _ := res.ToString()
	// log.Print(resJson)
	ass.NotEqual("kuiiii", gjson.Get(resJson, "args.name").String())
	ass.Equal("18", gjson.Get(resJson, "args.age").String())
}

func TestSetHeaders(t *testing.T) {
	ass := assert.New(t)
	cli := request.New("self", "target", 0)
	cli.SetUri(echoServ).SetMethod("get").SetPath("/get")
	cli.SetHeader("x-gender", "boy")
	cli.SetHeaders(map[string]string{
		"x-name": "lili",
		"x-age":  "18",
	})
	r, err := cli.Send()
	ass.NoError(err)
	if err != nil {
		return
	}
	bodyStr, err := r.ToString()
	ass.NoError(err)
	if err != nil {
		return
	}
	bodyParse := gjson.Parse(bodyStr)
	ass.NotEqual("boy", bodyParse.Get("headers.x-gender").String())
	ass.Equal("lili", bodyParse.Get("headers.x-name").String())
	ass.Equal("18", bodyParse.Get("headers.x-age").String())

}

func TestAddHeaders(t *testing.T) {
	ass := assert.New(t)
	cli := request.New("self", "target", 0)
	cli.SetUri(echoServ).SetMethod("get").SetPath("/get")
	cli.SetHeader("x-gender", "boy")
	cli.AddHeaders(map[string]string{
		"x-name": "lili",
		"x-age":  "18",
	})
	r, err := cli.Send()
	ass.NoError(err)
	if err != nil {
		return
	}
	bodyStr, err := r.ToString()
	ass.NoError(err)
	if err != nil {
		return
	}
	bodyParse := gjson.Parse(bodyStr)
	ass.Equal("boy", bodyParse.Get("headers.x-gender").String())
	ass.Equal("lili", bodyParse.Get("headers.x-name").String())
	ass.Equal("18", bodyParse.Get("headers.x-age").String())

}
