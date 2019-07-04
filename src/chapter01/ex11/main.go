package main

import (
"fmt"
"io"
"io/ioutil"
"net/http"
"os"
"time"
)

// より長い引数リストでfetchallを試してみる

// 下記のコマンドでhttp://www.alexa.com/topsitesから上位50件のサイトを取得
// $ curl https://www.alexa.com/topsites | grep 'href="/siteinfo/' | cut -d"/" -f3 | cut -d'"' -f1 | xargs -IX echo http://X | tr "\n" " "
// http://google.com http://youtube.com http://facebook.com http://baidu.com http://wikipedia.org http://yahoo.com http://qq.com http://taobao.com http://tmall.com http://google.co.in http://twitter.com http://instagram.com http://amazon.com http://sohu.com http://reddit.com http://vk.com http://live.com http://jd.com http://yandex.ru http://weibo.com http://sina.com.cn http://360.cn http://login.tmall.com http://google.co.jp http://blogspot.com http://netflix.com http://google.com.hk http://linkedin.com http://pornhub.com http://google.com.br http://google.co.uk http://pages.tmall.com http://twitch.tv http://yahoo.co.jp http://csdn.net http://mail.ru http://google.ru http://alipay.com http://google.fr http://google.de http://office.com http://ebay.com http://microsoft.com http://bing.com http://t.co http://microsoftonline.com http://xvideos.com http://aliexpress.com http://livejasmin.com http://msn.com

// ウェブサイトが応答しないときは、30secでタイムアウトした
// Get http://microsoftonline.com: dial tcp 52.178.167.109:80: i/o timeout
// 30.00s elapsed


func main() {
	start := time.Now()

	ch := make(chan string)
	for _, url := range os.Args[2:] {
		go fetch(url, ch)
	}
	for range os.Args[2:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
