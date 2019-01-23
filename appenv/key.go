package appenv

import (
	"fmt"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func DecodeRsaKey(buf []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(buf)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("Failed to decode PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse PKI public key: %v", err)
	}
	if rsaPubKey, ok := pub.(*rsa.PublicKey); ok {
		return rsaPubKey, nil
	} else {
		return nil, fmt.Errorf("Failed to dereference a PublicKey value.")
	}
}
