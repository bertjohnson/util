package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
	"strings"
	"time"
)

// GenerateAESKey creates a new AES key for encryption.
func GenerateAESKey(keySize int) (string, error) {
	key := make([]byte, keySize)
	n, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	if n != keySize {
		return "", errors.New("fewer than 64 bytes were read in")
	}

	return hex.EncodeToString(key[:]), nil
}

// DecryptAES decrypts an AES-protected message using GCM.
func DecryptAES(key string, inputData []byte) ([]byte, error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	if len(inputData) < 12 {
		return nil, errors.New("nonce not included within encrypted data")
	}

	nonce := inputData[0:12]

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	outputData, err := aesGCM.Open(nil, nonce, inputData[12:], nil)
	if err != nil {
		return nil, err
	}

	return outputData, nil
}

// EncryptAES encrypts a message with AES using GCM.
func EncryptAES(key string, inputData []byte) ([]byte, error) {
	decodedKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	outputData := aesGCM.Seal(nil, nonce, inputData, nil)
	return append(nonce, outputData...), nil
}

// GenerateX509Certificate generates an X509 certificate for TLS.
// Based on https://golang.org/src/crypto/tls/generate_cert.go.
func GenerateX509Certificate(hostname string, organizationName string, validFrom time.Time, validFor time.Duration, isCertAuthority bool, rsaBits int, ecdsaCurve string) ([]byte, []byte, error) {
	// Validate inputs.
	if hostname == "" {
		return nil, nil, errors.New("hostname not specified")
	}
	if validFor < 0 {
		return nil, nil, errors.New("validFor required")
	}

	// Generate private key.
	var privateKey interface{}
	var err error
	switch strings.ToLower(ecdsaCurve) {
	case "", "rsa":
		privateKey, err = rsa.GenerateKey(rand.Reader, rsaBits)
	case "p224":
		privateKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "p256":
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "p384":
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "p521":
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, nil, errors.New("unrecognized elliptic curve: " + ecdsaCurve)
	}
	if err != nil {
		return nil, nil, err
	}

	// Calculate dates.
	validTo := validFrom.Add(validFor)

	// Generate serial number.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	// Generate certificate.
	template := x509.Certificate{
		BasicConstraintsValid: true,
		DNSNames:              []string{hostname},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  isCertAuthority,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotBefore:             validFrom,
		NotAfter:              validTo,
		SerialNumber:          serialNumber,
		Subject: pkix.Name{
			Organization: []string{organizationName},
		},
	}
	if isCertAuthority {
		template.KeyUsage = template.KeyUsage | x509.KeyUsageCertSign
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey)
	if err != nil {
		return nil, nil, err
	}
	crtBuffer := new(bytes.Buffer)
	keyBuffer := new(bytes.Buffer)
	err = pem.Encode(crtBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, nil, err
	}
	pemBlockForPrivateKey, err := pemBlockForKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	err = pem.Encode(keyBuffer, pemBlockForPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	return crtBuffer.Bytes(), keyBuffer.Bytes(), nil
}

// pemBlockForKey returns a PEM-encoded key.
func pemBlockForKey(priv interface{}) (*pem.Block, error) {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}, nil
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}, nil
	default:
		return nil, errors.New("unknown key type")
	}
}

// publicKey returns the public key associated with a private key.
func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}
