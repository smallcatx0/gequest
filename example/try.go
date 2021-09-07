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
	}
}
