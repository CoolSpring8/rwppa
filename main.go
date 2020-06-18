package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/elazarl/goproxy"
	goproxy_html "github.com/elazarl/goproxy/ext/html"
)

var someURLMatcher = regexp.MustCompile(`/web/[0-3]/https?/[0-2]/`)

func main() {
	addr := flag.String("addr", ":8080", "proxy listen address")
	twfid := flag.String("id", "", "TWFID cookie got from rvpn.zju.edu.cn")
	flag.Parse()

	setCA(caCert, caKey)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			req.AddCookie(&http.Cookie{Name: "TWFID", Value: *twfid})

			newURL, err := url.Parse("https://rvpn.zju.edu.cn/web/2/" + req.URL.Scheme + "/0/" + req.URL.Host + req.URL.Path) // port has been included in "Host" param
			if err != nil {
				return req, nil // TODO: better handling
			}
			req.URL = newURL

			return req, nil
		})

	proxy.OnResponse(goproxy_html.IsWebRelatedText).Do(goproxy_html.HandleString(
		func(s string, ctx *goproxy.ProxyCtx) string {
			c := someURLMatcher.ReplaceAllString(s, "//")
			return c
		}))

	fmt.Println("Current TWFID:" + *twfid)
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
