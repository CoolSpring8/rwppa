// This file was partly from https://github.com/elazarl/goproxy/blob/master/examples/goproxy-customca/cert.go

// Copyright (c) 2012 Elazar Leibovich. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Elazar Leibovich. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/elazarl/goproxy"
)

var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIUdWndr4j2nJH3C+DfjL0sqG8PIMQwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMDA2MTgwNzQ1NDBaFw0yMzA0
MDgwNzQ1NDBaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggIiMA0GCSqGSIb3DQEB
AQUAA4ICDwAwggIKAoICAQDHiQ7JbOCvmEw+VYQY+4iLwtTHe9GyH3zXYujeeqC5
rcBU5MBD8NyqBjSDi/JPXEiCpBkYzmDtmqiCEmuoi8r6LADoJUdo2QUzu+S9UVt7
FqgXRNzZC5Enh1F9m/pkyz7+GtuQB4PIQGf74FNlza5Bd+9NmJrDYKxj6ydXoXTS
hTQrMEx3Wy8PzqpJw1006/So8Nh4jeu41DJr43bcO2J09/MpOd9pzB2GiDgrHu5M
5EAUJONWut62oWT42PYXzP2MRZPRmcd70/u5IX3LsV5Pi6XVqsFjjmIYnvKOa9rw
XICiCVdv5374CLwlTh4AnsstIjKpGOfWNx2cPXEDHtUQUXW+aaoNTR8gZOEs68qr
EpKz87xbl4pLo33nOFOSs0JU24nRGfLwOF3vN/usPqFEUyqquSWSOWfm3XPgKZlm
JrpmI/MgJ6YZunEujdvogT479RZ/IGq/Z/MhZaJj4iOcdq0iZqPUaKtV7v1lrl8S
1Z2O5fiVdoc/95n1D2ZQW7M5/+d9RiGuGSfVD95VA5WmQdGhRlV+MdeoGvvC+/F3
1pMl61YBvGc7jpPM6GA535Cmb9yXZFiKsutTHuDWEMngjfvlFV2qrxNy9PX+AQ+W
XbidvCZKQqxE6fgnFIQxsI82ZCkdddI9d9ifMt4qkFNeJDbyeNxPuea0qb4752er
5wIDAQABo1MwUTAdBgNVHQ4EFgQUYB+DRN17B+CMNJCRMRlvrR26HbAwHwYDVR0j
BBgwFoAUYB+DRN17B+CMNJCRMRlvrR26HbAwDwYDVR0TAQH/BAUwAwEB/zANBgkq
hkiG9w0BAQsFAAOCAgEAuWn+Vea4esUk5HCZH5KVd43jhoN5qZ0EOVB0jNQ+JayW
/n3FO+o6N+CTqT1gU11uUp1H54N24s6oj+a/VuFzy38Delu34qKA8Y+rSnM/IgxW
Uo6M1exgPQ0KqSZvbiJ0JamdIaO14OL5hS7ZzNWo3ZnvDaG04tfOAUPPAIqRt1cq
ilRemHc4WmJ/U/7Vf1HaJ+QWEMQUVg0WbDHfsVP2EY+L1oQD7ZEoPKkCLcd/pGeP
1eoGxV6vKJL7JD8Sdl1UoHPksp/1+LIB4s6JE68NrksY3TqF5vC8jKHiG67h0XJ8
WoWriNK/T6tSLj/ApAL1qeyz1nG1iccGJ4gaBOFH0kDl1vkU2qTbMQJv59cLUHWv
boevsBp356BiCBMYYkTvLAS4LDA27tvFcuzHtnaxQi22BnwZap9dvQJuxxtvfgCj
zeqmrCxwuuQpPvhHO4DJQxMxxeado6IyQmuVU6VJAnobBI2KIiGRDjldRqGe+aFe
6sTqW+NTs5ehkf0O12JTE9V2/yxALBd8BWEFhfaZfjU4g22TSF+1oPgd2r8qvmg6
fICwA9tna6N36NyDnic1V+xF+93wBxIfhl+ZtzMISzxYlh+zyfipvz4qJb0Hhk0S
u8FUK+ksijxQTFP/ZFXNw/r8ufC1mHsWbD0bCqRhb3eVblU1f9e2XLmW4LUgtII=
-----END CERTIFICATE-----`)

var caKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)

func setCA(caCert, caKey []byte) error {
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
