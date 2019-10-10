package memo_test

import (
	"testing"

	"../memo"
)

var httpGetBody = memotest.HTTPGetBody

func TestConcurrentCancelAll(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.ConcurrentCancelAll(t, m)
}

func TestConcurrentCancelOdd(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.ConcurrentCancelOdd(t, m)
}

// TODO: リクエスト自体をキャンセルするようになっていないので、キャンセルされない
// 並行でリクエストすると奇数番目のレスポンスのキャッシュを返してしまう。
func TestConcurrentCancelEven(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.ConcurrentCancelEven(t, m)
}

func TestConcurrentCancelAllAfter1000ms(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.ConcurrentCancelAllWait(t, m, 1000)
}

func Test(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.Sequential(t, m)
}

func TestConcurrent(t *testing.T) {
	m := memo.New(httpGetBody)
	defer m.Close()
	memotest.Concurrent(t, m)
}
