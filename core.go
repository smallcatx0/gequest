package request

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// var ALLOW_METHODS = []string{"GET", "POST", "PUT", "DELETE"}
var ALLOW_METHODS = map[string]int{"GET": 1, "POST": 1, "PUT": 1, "DELETE": 1}

type Core struct {
	serviceName       string
	targetServiceName string
	uri               string
	path              string
	method            string
	headers           map[string]string
	json              string
	body              []byte
	query             url.Values
	errs              []error
	client            *http.Client
	response          *http.Response
}

func (c *Core) String() string {
	res := map[string]interface{}{
		"serviceName":       c.serviceName,
		"targetServiceName": c.targetServiceName,
		"uri":               c.uri,
		"path":              c.path,
		"method":            c.method,
		"headers":           c.headers,
		"json":              c.json,
		"body":              string(c.body),
		"query":             c.query,
		"errs":              c.errs,
	}
	if c.response != nil {
		res["response"], _ = c.Response().ToString()
	}
	jstr, _ := json.Marshal(res)
	return string(jstr)
}

type Response struct {
	*http.Response
}

// Read response body into a byte slice.
func (r *Response) ReadAll() ([]byte, error) {
	var reader io.ReadCloser
	var err error
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = r.Body
	}

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// Read response body into string.
func (r *Response) ToString() (string, error) {
	bytes, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func New(serviceName, targetServiceName string, timeOutMs int) *Core {
	c := &Core{
		client: &http.Client{
			Timeout: time.Millisecond * time.Duration(timeOutMs),
		},
		headers: make(map[string]string),
	}
	if serviceName != "" {
		c.serviceName = serviceName
		c.SetHeader("Service-Name", serviceName)
	}
	if targetServiceName != "" {
		c.targetServiceName = targetServiceName
		c.SetHeader("Target-Service-Name", targetServiceName)
	}
	return c
}

func (c *Core) SetMethod(method string) *Core {
	method = strings.ToUpper(method)
	_, ok := ALLOW_METHODS[method]
	if ok {
		c.method = method
	} else {
		c.errs = append(c.errs, errors.New("method is not allow instandard http method"))
	}
	return c
}

func (c *Core) SetUri(uri string) *Core {
	c.uri = uri
	return c
}
func (c *Core) SetPath(path string) *Core {
	c.path = path
	return c
}

func (c *Core) SetJson(param interface{}) *Core {
	c.headers["Content-Type"] = "application/json"
	jsonData, _ := json.Marshal(param)
	c.json = string(jsonData)
	return c
}

func (c *Core) SetBodyText(body string) *Core {
	c.body = []byte(body)
	c.SetHeader("Content-Type", "text/plain")
	return c
}

func (c *Core) SetBody(data []byte) *Core {
	c.body = data
	return c
}

func (c *Core) SetQuery(param url.Values) *Core {
	c.query = param
	return c
}

func (c *Core) SetHeader(k string, v string) *Core {
	c.headers[k] = v
	return c
}

func (c *Core) SetHeaders(headers map[string]string) *Core {
	c.headers = headers
	return c
}

func (c *Core) Send() (r *Response, err error) {
	if len(c.errs) != 0 {
		return nil, c.errs[0]
	}

	// 组装body
	var payload io.Reader
	if c.json != "" {
		payload = strings.NewReader(c.json)
	} else if len(c.body) != 0 {
		payload = bytes.NewReader(c.body)
	}
	req, err := http.NewRequest(c.method, c.getUrl(), payload)
	if err != nil {
		return nil, err
	}
	// 组装header
	for k, v := range c.headers {
		req.Header.Add(k, v)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	c.response = res
	return &Response{res}, nil
}

func (c *Core) Response() *Response {
	return &Response{c.response}
}

func (c *Core) ResponseRaw() *http.Response {
	return c.response
}

func (c *Core) Errors() []error {
	return c.errs
}

func (c *Core) getUrl() string {
	domain := c.uri
	path := c.path
	if c.uri == "" {
		domain = "http://" + c.targetServiceName
	}
	queryStr := ""
	if len(c.query) != 0 {
		queryStr = "?" + c.query.Encode()
	}
	return domain + path + queryStr
}
