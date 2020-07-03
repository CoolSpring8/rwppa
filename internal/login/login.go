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

package login

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	// endpointURL is where ZJU RVPN web portal's login interface locates.
	// TODO: try to find out if some of the parameters are not necessary.
	endpointURL string = "https://rvpn.zju.edu.cn/por/login_psw.csp?type=cs&dev=android-phone&dev=android-phone&language=zh_CN"
)

// loginRVPNWebPortal takes username and password, returns TWFID.
// TWFID is used by the web portal for authentication.
// Incorrect or empty username and password will simply lead to a useless TWFID.
// TODO: verify before returning. throw an error if not working.
func LoginRVPNWebPortal(username, password string) string {
	data := url.Values{}
	data.Set("svpn_name", username)
	data.Set("svpn_password", password)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", endpointURL, strings.NewReader(data.Encode()))

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	twfid := resp.Cookies()[0].Value

	return twfid
}
