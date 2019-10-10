package memo

import "errors"

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res    result
	ready  chan struct{}
	cancel <-chan struct{}
}

type request struct {
	key      string
	response chan<- result
	cancel   <-chan struct{}
}

type Memo struct{ requests chan request }

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string, cancel <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response, cancel}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil || e.canceled() {
			e = &entry{ready: make(chan struct{}), cancel: req.cancel}
			cache[req.key] = e
			go e.call(f, req.key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	e.res.value, e.res.err = f(key)
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	select {
	case <-e.ready:
		response <- e.res
	case <-e.cancel:
		response <- result{nil, errors.New("request is canceled")}
	}
}

func (e *entry) canceled() bool {
	select {
	case <-e.cancel:
		return true
	default:
		return false
	}
}
