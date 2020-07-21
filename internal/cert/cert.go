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

// Package cert provides operation functions on CA cert and key.
package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

// GetCA returns a CA cert, a CA key or an error in file reading or writing process.
// If one of rootCA.crt and rootCA.key does not exist,
// new ones will be generated, wrote to current dir and returned.
// TODO: do not hardcode filename.
func GetCA() ([]byte, []byte, error) {
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
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
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
