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
// or any other HTTP-proxy-capable requesters, like clients that only utilize HTTP protocol.

// In short, users are given the ability to access ZJU intranet sites,
// with a ZJU internet service account required, and via ZJU RVPN web portal (view rvpn.zju.edu.cn on phones to see).
// Hopefully it can replace the role of Sangfor EasyConnect, to a certain extent.

// Only HTTP(s) protocol is supported. WebSocket, FTP, SSH and any other ones are not.

// This program is powered by an MITM-like procedure, so use it at your own risk.

// Most of proxy functionalities are imported from https://github.com/elazarl/goproxy/ , An HTTP proxy library for Go.

package main

import (
	"fmt"
	"os"

	"github.com/coolspring8/rwppa/internal/proxy"
	"github.com/coolspring8/rwppa/pkg/rvpn"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	// Stay calm with this temporary solution, we are going to refactor all these codes on input.
	var username string
	var password string
	var listenAddr string
	fmt.Println("Username:")
	_, err := fmt.Scanln(&username)
	if err != nil {
		panic(err)
	}
	fmt.Println("Password:")
	pw, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	password = string(pw)
	fmt.Println("")
	fmt.Println("Listen Address:")
	_, err = fmt.Scanln(&listenAddr)
	if err != nil {
		panic(err)
	}

	w := rvpn.WebPortal{Username: username, Password: password}
	twfid, err := w.DoLogIn()
	if err != nil {
		panic(err)
	}

	proxy.StartProxyServer(listenAddr, *twfid)
}
