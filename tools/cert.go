package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	log "github.com/ronaksoft/rony/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

/*
   Creation Time: 2020 - Aug - 25
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// helper function to create a cert template with a serial number and other required fields
func CertTemplate() (*x509.Certificate, error) {
	// generate a random serial number (a real cert authority would have some logic behind this)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Yacht, Inc."}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour), // valid for an hour
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

func CreateCert(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (cert *x509.Certificate, certPEM []byte, err error) {
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	}
	// parse the resulting certificate so we can use it again
	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}
	// PEM encode the certificate (this is a standard TLS encoding)
	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}

func GenerateSelfSignedCerts(keyPath, certPath string) {
	// generate a new key-pair
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("generating random key:", zap.Error(err))
	}

	rootCertTmpl, err := CertTemplate()
	if err != nil {
		log.Fatal("creating cert template:", zap.Error(err))
	}
	// describe what the certificate will be used for
	rootCertTmpl.IsCA = true
	rootCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	rootCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	rootCertTmpl.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}

	rootCert, rootCertPEM, err := CreateCert(rootCertTmpl, rootCertTmpl, &rootKey.PublicKey, rootKey)
	if err != nil {
		log.Fatal("error creating cert", zap.Error(err))
	}
	// fmt.Printf("%s\n", rootCertPEM)
	// fmt.Printf("%#x\n", rootCert.Signature) // more ugly binary
	_ = rootCert

	// PEM encode the private key
	rootKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rootKey),
	})

	ioutil.WriteFile(certPath, rootCertPEM, os.ModePerm)
	ioutil.WriteFile(keyPath, rootKeyPEM, os.ModePerm)

	return
}

func GetCertificate(keyPath, certPath string) tls.Certificate {
	keyPEM, _ := ioutil.ReadFile(keyPath)
	certPEM, _ := ioutil.ReadFile(certPath)
	// Create a TLS cert using the private key and certificate
	rootTLSCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatal("invalid key pair", zap.Error(err))
	}
	return rootTLSCert

}
