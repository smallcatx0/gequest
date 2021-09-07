package request

import (
	"sync"
)

type Resp struct {
	Index int
	Core  *Core
	Err   error
}

// 并发请求
func MultRequest(requests ...*Core) []Resp {
	var wg sync.WaitGroup
	wg.Add(len(requests))
	ch := make(chan Resp, len(requests))
	for i, req := range requests {
		go func(i int, req *Core, ch chan Resp) {
			defer wg.Done()
			_, err := req.Send()
			ch <- Resp{
				Index: i,
				Core:  req,
				Err:   err,
			}
		}(i, req, ch)
	}
	wg.Wait()
	close(ch)
	// 将返回转化为数组
	resps := make([]Resp, len(requests))
	for one := range ch {
		resps[one.Index] = one
	}
	return resps
}
