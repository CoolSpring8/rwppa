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

package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/elazarl/goproxy"
)

// getCA returns a CA cert, a CA key or an error in file reading or writing process.
// If one of rootCA.crt and rootCA.key does not exist,
// new ones will be generated, wrote to current dir and returned.
// TODO: do not hardcode filename.
func getCA() ([]byte, []byte, error) {
	if fileExists("rootCA.crt") && fileExists("rootCA.key") {
		caCert, err := ioutil.ReadFile("rootCA.crt")
		if err != nil {
			return nil, nil, err
		}
		caKey, err := ioutil.ReadFile("rootCA.key")
		if err != nil {
			return nil, nil, err
		}
		return caCert, caKey, err
	}

	// files do not exist, make new ones instead
	// source: https://golang.org/src/crypto/tls/generate_cert.go
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	panicOnErr(err)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	notBefore := time.Now()
	notAfter := notBefore.Add(3 * 365 * 24 * time.Hour) // 3 years
	panicOnErr(err)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"github.com/coolspring8/rwppa"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	panicOnErr(err)
	certOut, err := os.Create("rootCA.crt")
	panicOnErr(err)
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	panicOnErr(err)
	err = certOut.Close()
	panicOnErr(err)

	keyOut, err := os.OpenFile("rootCA.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	panicOnErr(err)
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	panicOnErr(err)
	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	panicOnErr(err)
	err = keyOut.Close()
	panicOnErr(err)

	caCert, err := ioutil.ReadFile("rootCA.crt")
	if err != nil {
		return nil, nil, err
	}
	caKey, err := ioutil.ReadFile("rootCA.key")
	if err != nil {
		return nil, nil, err
	}
	return caCert, caKey, err
}

// setCA takes CA cert and CA key, and sets up goproxy CA.
// Returns error in parsing and setting CA.
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

// fileExists checks if a file exists and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// panicOnErr panics on err not equal to nil.
// TODO: better handling of errors, don't panic or you should recover at some time.
func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
