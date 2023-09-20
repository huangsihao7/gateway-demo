package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

// 当前服务器地址
var addr = "127.0.0.1:2222"

func main() {

	//设置下游的地址
	rs1 := "http://127.0.0.1:2003"
	url1, err := url.Parse(rs1)
	if err != nil {
		log.Println(err)
	}

	proxy := NewSingleHostReverseProxy(url1)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery

	//设置新的请求信息
	director := func(req *http.Request) {
		//url_rewrite
		//127.0.0.1:2002/dir/abc ==> 127.0.0.1:2003/base/abc ??
		//127.0.0.1:2002/dir/abc ==> 127.0.0.1:2002/abc
		//127.0.0.1:2002/abc ==> 127.0.0.1:2003/base/abc
		re, _ := regexp.Compile("^/dir(.*)")
		req.URL.Path = re.ReplaceAllString(req.URL.Path, "$1")

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		//target.Path : /base
		//req.URL.Path : /dir
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}

	//修改返回值
	modifyFunc := func(res *http.Response) error {
		if res.StatusCode != 200 {
			return errors.New("error statusCode")

		}
		oldPayload, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		newPayLoad := []byte("hello " + string(oldPayload))
		res.Body = io.NopCloser(bytes.NewBuffer(newPayLoad))
		res.ContentLength = int64(len(newPayLoad))
		res.Header.Set("Content-Length", fmt.Sprint(len(newPayLoad)))
		return nil
	}
	//错误处理
	errorHandler := func(res http.ResponseWriter, req *http.Request, err error) {
		res.Write([]byte(err.Error()))
	}
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyFunc, ErrorHandler: errorHandler}
}

// 重新组合URL
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
