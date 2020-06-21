// Copyright (C) 2020  CoolSpring8

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
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
	rvpnURLMatcher *regexp.Regexp = regexp.MustCompile(`/web/[0-3]/(https?)/[0-2]/`)

	movedLocationURLMatcher *regexp.Regexp = regexp.MustCompile(`https://.*:443/web/[0-3]/(https?)/[0-2]/`)

	isOPTIONSRequest goproxy.ReqCondition = goproxy.ReqConditionFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			return req.Method == "OPTIONS"
		})

	hasMovedLocationHeader goproxy.RespCondition = goproxy.RespConditionFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
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

func startProxyServer(listenAddr string, twfid string) {
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
			req.AddCookie(&http.Cookie{Name: "TWFID", Value: twfid})

			return req, nil
		})
	proxy.OnRequest(isOPTIONSRequest).DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			resp := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "")
			resp.Header.Add("Access-Control-Allow-Credentials", "true")
			resp.Header.Add("Access-Control-Allow-Headers", "authorization")
			resp.Header.Add("Access-Control-Allow-Methods", "GET, POST, HEAD, DELETE, OPTIONS")
			resp.Header.Add("Access-Control-Allow-Origin", "*")
			return req, resp
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
			c := rvpnURLMatcher.ReplaceAllString(s, "$1://")
			rawURLWithPort := ctx.UserData.(reqData).rawURLWithPort
			rawURLWithoutPort := ctx.UserData.(reqData).rawURLWithoutPort
			c = strings.ReplaceAll(c, rawURLWithPort[:strings.LastIndex(rawURLWithPort, "/")+1], "") // possible out of bounds?
			c = strings.ReplaceAll(c, rawURLWithoutPort[:strings.LastIndex(rawURLWithoutPort, "/")+1], "")
			return c
		}))

	fmt.Println("Current TWFID:" + twfid)
	fmt.Println("Listen Address:" + listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, proxy))
}
