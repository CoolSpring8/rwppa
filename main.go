package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/elazarl/goproxy"
	goproxy_html "github.com/elazarl/goproxy/ext/html"
)

var (
	rvpnURLMatcher = regexp.MustCompile(`/web/[0-3]/(https?)/[0-2]/`)

	movedLocationURLMatcher = regexp.MustCompile(`https://.*:443/web/[0-3]/(https?)/[0-2]/`)

	hasMovedLocationHeader = goproxy.RespConditionFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		return resp.Header.Get("Location") != ""
	})

	isWebRelatedText goproxy.RespCondition = goproxy.ContentTypeIs(
		"text/html",
		"text/css",
		"text/javascript", "application/javascript", "application/x-javascript",
		"text/xml",
		"text/json")
)

type reqData struct {
	rawURLWithPort    string
	rawURLWithoutPort string
}

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
			// store rawURL for later use
			rawURLWithPort := req.URL.String()
			hostWithoutPort, _, err := net.SplitHostPort(req.URL.Host)
			if err != nil {
				hostWithoutPort = req.URL.Host // http at port 80 causes "missing port in address"
			}
			rawURLWithoutPort := req.URL.Scheme + "://" + hostWithoutPort + req.URL.Path
			ctx.UserData = reqData{rawURLWithPort, rawURLWithoutPort}

			// new request target
			newURL, err := url.Parse("https://rvpn.zju.edu.cn/web/2/" + req.URL.Scheme + "/0/" + req.URL.Host + req.URL.Path) // port has been included in "Host" param
			if err != nil {
				return req, nil // this rarely happens?
			}
			req.URL = newURL

			// add cookie for web portal verification
			req.AddCookie(&http.Cookie{Name: "TWFID", Value: *twfid})

			return req, nil
		})

	proxy.OnResponse(hasMovedLocationHeader).DoFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			respLocation := resp.Header.Get("Location")
			newLocation := movedLocationURLMatcher.ReplaceAllString(respLocation, "$1://")
			resp.Header.Set("Location", newLocation)
			return resp
		})
	proxy.OnResponse(isWebRelatedText).Do(goproxy_html.HandleString(
		func(s string, ctx *goproxy.ProxyCtx) string {
			// fix link in page, and fix "src" issues in javascript files
			c1 := rvpnURLMatcher.ReplaceAllString(s, "$1://")
			rawURLWithPort := ctx.UserData.(reqData).rawURLWithPort
			rawURLWithoutPort := ctx.UserData.(reqData).rawURLWithoutPort
			c2 := strings.ReplaceAll(c1, rawURLWithPort[:strings.LastIndex(rawURLWithPort, "/")+1], "") // possible out of bounds?
			c3 := strings.ReplaceAll(c2, rawURLWithoutPort[:strings.LastIndex(rawURLWithoutPort, "/")+1], "")
			return c3
		}))

	fmt.Println("Current TWFID:" + *twfid)
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
