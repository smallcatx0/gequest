package request_test

import (
	"testing"

	request "github.com/smallcatx0/gequest"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestMultReq(t *testing.T) {
	requests := make([]*request.Core, 0, 5)
	for i := 0; i < 20; i++ {
		requests = append(requests,
			request.New("", "postman-echo.com", 0).
				SetMethod("post").
				SetPath("/post").
				SetJson(map[string]int{"index": i}),
		)
	}
	res := request.MultRequest(requests...)
	for i, one := range res {
		assert.NoError(t, one.Err)
		if one.Err != nil {
			panic(one.Err)
		}
		r := res[i].Core.Response()
		resJson, _ := r.ToString()
		assert.Equal(t, 200, r.StatusCode)
		assert.Equal(t, i, int(gjson.Get(resJson, "json.index").Int()))
	}
}
