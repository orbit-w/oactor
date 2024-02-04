package actor

/*
   @Author: orbit-w
   @File: request
   @2024 2月 周日 22:18
*/

const (
	Call      = 1
	AsyncCall = 2
)

type Request struct {
	category int
	msg      any
	ch       chan *Response
}

func (req *Request) Response(msg any, err error) {
	req.ch <- &Response{
		msg: msg,
		err: err,
	}
}

func (req *Request) Return() {
	reqPool.Put(req)
}

func (req *Request) Done() <-chan *Response {
	return req.ch
}

type Response struct {
	msg any
	err error
}
