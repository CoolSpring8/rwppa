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
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	// PreLogInURL is where a normal browser user will run into first,
	// and it is where we need to collect some variable necessary in the logging in request.
	PreLogInURL = "https://rvpn.zju.edu.cn/por/login_psw.csp"
	// LogInURL is where ZJU RVPN web portal's login interface locates.
	// "sfrnd" seems to be a constant.
	LogInURL = "https://rvpn.zju.edu.cn/por/login_psw.csp?sfrnd=2346912324982305&encrypt=1"
	// LogInSimpleURL is where ZJU RVPN web portal's simpler login interface locates.
	LogInSimpleURL = "https://rvpn.zju.edu.cn/por/login_psw.csp?type=cs&dev=android-phone&dev=android-phone&language=zh_CN"
	// RSAn is a constant taken from the script in the page.
	RSAn = "B81AB4511CCC90B27170266122DCA5496C7D6ECCEFE65830071B487C0457403B5FCAB8C788DF37F882E897984E28250E6B11879403D54F46355F2F0802BB776EC041035F50BA5A77221FED2A91D24BF2FE44160653A2824F650E458EB3AE15A4446514C89A7EE6213F4B3687C9AB13E6ABCE919676C37E7DDF0580DBDD5643642CFA2BE8513F523E9759CA3351C944A6533752728260C8EDEB9CD59FFA08CE57FB9B109CFB0881858EBE36E384D13D4D3DB80768C36CB62FDF67799AECA4EA23D101FAC43FF2C8B1165AC4A15D2A2BDC3A987298975AFAB5A9CF4B38B8DC27ECA0278335B1ACEB197FE583D4E1648FE8922B34E2B276A23F9E61A29058225839"
	// RSAe is a constant taken from the script in the page. 0x10001 is 65537 in decimal.
	RSAe = 65537
	// UserAgent is used as the User-Agent.
	UserAgent = "Mozilla/5.0 (Linux; Android 8.1.0; Redmi 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Mobile Safari/537.36"
)

// WebPortal refers to ZJU RVPN web portal.
type WebPortal struct {
	Username string // Username is ZJU network service account username.
	Password string // Password is ZJU network service account password.
}

// LogIn uses username and password to get a TWFID.
// TWFID is used by the web portal for authentication.
// Make sure to check the returned error value.
// This function performs a complicated series of operations inside,
// aiming to imitate normal phone browser users' behavior.
func (webPortal WebPortal) LogIn() (*string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	req, err := http.NewRequest("GET", PreLogInURL, nil)
	if err != nil {
		return nil, errors.New("parsing request for web portal pre-login error")
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("network error during pre-login to web portal")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	// TODO: Must we neglect defer statement and choose to close it directly?
	resp.Body.Close()

	// Not so good, but it works.
	randCodeFinder := regexp.MustCompile(`id="svpn_req_randcode" value="(\d{3,4})"`)
	randCode := randCodeFinder.FindSubmatch(body)
	if randCode == nil {
		return nil, errors.New("randcode not found during pre-login")
	}

	// twfidPreValue(TWFID) is always in the response's Set-Cookie header,
	// regardless of success or failure on the process.
	twfidPreValue := resp.Cookies()[0].Value

	n := new(big.Int)
	if _, ok := n.SetString(RSAn, 16); !ok {
		return nil, errors.New("parsing RSA n required for web portal login error")
	}
	pubKey := rsa.PublicKey{N: n, E: RSAe}
	e, err := rsa.EncryptPKCS1v15(rand.Reader, &pubKey, []byte(webPortal.Password))
	if err != nil {
		return nil, errors.New("encrypting password used in web portal login error")
	}
	encryptedPassword := hex.EncodeToString(e)

	data := url.Values{"svpn_req_randcode": {string(randCode[1])}, "svpn_name": {webPortal.Username}, "svpn_password": {encryptedPassword}}
	req, err = http.NewRequest("POST", LogInURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.New("parsing request for web portal login error")
	}
	req.AddCookie(&http.Cookie{Name: "TWFID", Value: twfidPreValue})
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://rvpn.zju.edu.cn")
	req.Header.Set("Referer", "https://rvpn.zju.edu.cn/por/login_psw.csp")
	req.Header.Set("User-Agent", UserAgent)
	resp, err = client.Do(req)
	if err != nil {
		return nil, errors.New("network error during login to web portal")
	}
	defer resp.Body.Close()

	// It is probably not worth parsing HTML, extracting the script piece,
	// and parsing or executing it just to access a tiny error variable. Is it?
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	errorInfoFinder := regexp.MustCompile(`error_info: '(.+)',`)
	error := errorInfoFinder.FindSubmatch(body)
	if error != nil {
		return nil, errors.New(string(error[1]))
	}

	// twfid(TWFID) is a new usable value returned in the response's Set-Cookie header.
	twfid := resp.Cookies()[0].Value

	return &twfid, nil
}

// LogInSimple also uses username and password to get a TWFID,
// but it has a simpler logic.
func (webPortal WebPortal) LogInSimple() (*string, error) {
	data := url.Values{}
	data.Set("svpn_name", webPortal.Username)
	data.Set("svpn_password", webPortal.Password)

	client := &http.Client{}
	req, err := http.NewRequest("POST", LogInSimpleURL, strings.NewReader(data.Encode()))
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
