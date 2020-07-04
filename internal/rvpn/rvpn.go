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

// Package rvpn deals with things related to ZJU RVPN web portal.
package rvpn

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	// endpointURL is where ZJU RVPN web portal's login interface locates.
	// TODO: find out if some of the parameters are not necessary.
	endpointURL string = "https://rvpn.zju.edu.cn/por/login_psw.csp?type=cs&dev=android-phone&dev=android-phone&language=zh_CN"
)

// WebPortal refers to ZJU RVPN web portal.
type WebPortal struct {
	Username string // Username is ZJU network service account username.
	Password string // Password is ZJU network service account password.
}

// LogIn uses username and password to get a TWFID.
// TWFID is used by the web portal for authentication.
// Incorrect or empty username and password will simply lead to a useless TWFID.
// TODO: verify, throw an error if not working.
func (webPortal WebPortal) LogIn() (*string, error) {
	data := url.Values{}
	data.Set("svpn_name", webPortal.Username)
	data.Set("svpn_password", webPortal.Password)

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpointURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	twfid := resp.Cookies()[0].Value

	return &twfid, nil
}
