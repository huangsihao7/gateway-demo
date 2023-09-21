package cmd

import (
	"bytes"
	"gateway-detail/load_balance/lb_factory"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	Addr      = "127.0.0.1:2224"
	transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, //连接超时
			KeepAlive: 30 * time.Second, //长连接超时时间
		}).DialContext,
		MaxIdleConns:          100,              //最大空闲连接
		IdleConnTimeout:       90 * time.Second, //空闲超时时间
		TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
		ExpectContinueTimeout: 1 * time.Second,  //100-continue状态码超时时间
	}
)

func NewMultipleHostsReverseProxy(lb lb_factory.LoadBalance) *httputil.ReverseProxy {
	// 请求协调者
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			log.Fatal("get next addr fail")

		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			log.Fatal(err)
		}
		//改造req
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	//更改内容
	modifyFunc := func(resp *http.Response) error {
		//请求以下命令：curl 'http://127.0.0.1:2002/error'
		if resp.StatusCode != 200 {
			//获取内容
		}
		//追加内容
		oldPayload, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		newPayload := []byte("Hello :" + string(oldPayload))
		resp.Body = io.NopCloser(bytes.NewBuffer(newPayload))
		resp.ContentLength = int64(len(newPayload))
		resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(newPayload)), 10))
		return nil
	}
	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		//todo 如果是权重的负载则调整临时权重
		http.Error(w, "ErrorHandler error:"+err.Error(), 500)
	}
	return &httputil.ReverseProxy{Director: director, Transport: transport, ModifyResponse: modifyFunc, ErrorHandler: errFunc}
}
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
