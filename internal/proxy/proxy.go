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

// Package proxy provides utilities to set up a proxy.
package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/coolspring8/rwppa/internal/cert"
	"github.com/elazarl/goproxy"
	goproxy_html "github.com/elazarl/goproxy/ext/html"
)

var (
	// RVPNURLMatcher detects RVPN-modified URLs.
	RVPNURLMatcher *regexp.Regexp = regexp.MustCompile(`/web/[0-3]/(https?)/[0-2]/`)

	// MovedLocationURLMatcher detects RVPN-modified URL in the 3xx response's Location header.
	MovedLocationURLMatcher *regexp.Regexp = regexp.MustCompile(`https://.*:443/web/[0-3]/(https?)/[0-2]/`)

	// IsOPTIONSRequest checks whether the request's method is OPTIONS.
	IsOPTIONSRequest goproxy.ReqCondition = goproxy.ReqConditionFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
			return req.Method == "OPTIONS"
		})

	// HasMovedLocationHeader checks whether the request has a Location header.
	HasMovedLocationHeader goproxy.RespCondition = goproxy.RespConditionFunc(
		func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
			if resp == nil {
				return false
			}
			return resp.Header.Get("Location") != ""
		})

	// IsWebRelatedText checks whether the response's content type is one of html, js, css, xml and json.
	IsWebRelatedText goproxy.RespCondition = goproxy.ContentTypeIs(
		"text/html",
		"text/css",
		"text/javascript", "application/javascript", "application/x-javascript",
		"text/xml",
		"text/json")
)

// reqData stores a request's raw URL, and a port-stripped one.
// rawURLWithoutPort is for URLs like https://example.com:443/ ,
// where in href-s, port 443 does not show up but the URLs just need to be processed.
type reqData struct {
	rawURLWithPort    string
	rawURLWithoutPort string
}

// StartProxyServer starts a proxy server listening at the given address with the given TWFID.
func StartProxyServer(listenAddr, twfid string) {
	caCert, caKey, err := cert.GetCA()
	if err != nil {
		panic(err)
	}
	err = SetCA(caCert, caKey)
	if err != nil {
		panic(err)
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().Do(goproxy.FuncReqHandler(SendToRVPNHandler(twfid)))
	proxy.OnRequest(IsOPTIONSRequest).DoFunc(OPTIONSRequestHandler)

	proxy.OnResponse(HasMovedLocationHeader).DoFunc(MovedLocationHandler)
	proxy.OnResponse(IsWebRelatedText).Do(WebTextHandler())

	fmt.Println("Current TWFID:" + twfid)
	fmt.Println("Listen Address:" + listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, proxy))
}

// SetCA takes CA cert and CA key, and sets up goproxy CA.
// Returns error in parsing and setting CA.
func SetCA(caCert, caKey []byte) error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

// SendToRVPNHandler sends the original request to RVPN web portal.
func SendToRVPNHandler(twfid string) func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		// store rawURL for later use
		rawURLWithPort := req.URL.String()
		hostWithoutPort, _, err := net.SplitHostPort(req.URL.Host)
		if err != nil {
			hostWithoutPort = req.URL.Host // http at port 80 causes "missing port in address"
		}
		rawURLWithoutPort := req.URL.Scheme + "://" + hostWithoutPort + req.URL.Path // Query isn't here, but it doesn't matter
		ctx.UserData = reqData{rawURLWithPort, rawURLWithoutPort}

		// new request target
		// TODO: add other missed fields in "type URL struct" if necessary, like "User"
		// port has been included in "Host"
		if req.URL.RawQuery != "" {
			newURL, err := url.Parse("https://rvpn.zju.edu.cn/web/2/" + req.URL.Scheme + "/0/" + req.URL.Host + req.URL.Path + "?" + req.URL.RawQuery)
			if err != nil {
				return req, nil // this rarely happens?
			}
			req.URL = newURL
		} else {
			newURL, err := url.Parse("https://rvpn.zju.edu.cn/web/2/" + req.URL.Scheme + "/0/" + req.URL.Host + req.URL.Path)
			if err != nil {
				return req, nil
			}
			req.URL = newURL
		}

		// add cookie for web portal verification
		req.AddCookie(&http.Cookie{Name: "TWFID", Value: twfid})

		return req, nil
	}
}

// OPTIONSRequestHandler skips sending the request and returns a response designed for CORS preflight request.
func OPTIONSRequestHandler(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	resp := goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "")
	resp.Header.Add("Access-Control-Allow-Credentials", "true")
	resp.Header.Add("Access-Control-Allow-Headers", "authorization, content-type")
	resp.Header.Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD, OPTIONS")
	resp.Header.Add("Access-Control-Allow-Origin", "*")
	return req, resp
}

// MovedLocationHandler fixes redirections.
func MovedLocationHandler(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	respLocation := resp.Header.Get("Location")
	newLocation := MovedLocationURLMatcher.ReplaceAllString(respLocation, "$1://")
	resp.Header.Set("Location", newLocation)
	return resp
}

// WebTextHandler fixes links in page and "src" issues in javascript files.
// This solution, however, may prevent content streaming. Fix it?
func WebTextHandler() goproxy.RespHandler {
	return goproxy_html.HandleString(
		func(s string, ctx *goproxy.ProxyCtx) string {
			c := RVPNURLMatcher.ReplaceAllString(s, "$1://")
			rawURLWithPort := ctx.UserData.(reqData).rawURLWithPort
			rawURLWithoutPort := ctx.UserData.(reqData).rawURLWithoutPort
			c = strings.ReplaceAll(c, rawURLWithPort[:strings.LastIndex(rawURLWithPort, "/")+1], "") // possible out of bounds?
			c = strings.ReplaceAll(c, rawURLWithoutPort[:strings.LastIndex(rawURLWithoutPort, "/")+1], "")
			return c
		})
}
