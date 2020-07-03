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

// Rwppa is a tool that exposes an HTTP proxy service intercepting requests to pass through ZJU RVPN web portal.
// In other words, on receiving requests, they will be sent to rvpn.zju.edu.cn and corresponding results will be passed back to the browser,
// or any other HTTP-proxy-capable requesters, like clients that only utilizes HTTP protocol.

// In short, users are given the ability to browse ZJU intranet sites,
// with a ZJU internet service account required, and via ZJU RVPN web portal (view rvpn.zju.edu.cn on phones to see).
// Hopefully it can replace the role of Sangfor EasyConnect, to a certain extent.

// Only HTTP(s) protocol is supported. WebSocket, FTP, SSH and any other ones are not.

// This program is powered by an MITM-like procedure, so use it at your own risk.

// Most of proxy functionalities are imported from https://github.com/elazarl/goproxy/ , An HTTP proxy library for Go.

package main

import (
	"fmt"

	"github.com/coolspring8/rwppa/internal/login"
	"github.com/coolspring8/rwppa/internal/proxy"
)

func main() {
	var username string
	var password string
	var listenAddr string
	_, err := fmt.Scanln(&username)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Scanln(&password)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Scanln(&listenAddr)
	if err != nil {
		panic(err)
	}
	twfid := login.LoginRVPNWebPortal(username, password)
	proxy.StartProxyServer(listenAddr, twfid)
}
