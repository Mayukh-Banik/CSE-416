package utils

import (
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "errors"
	"crypto"
)

// ParsePublicKey parses a PEM encoded public key string
func ParsePublicKey(pemKey string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(pemKey))
    if block == nil || block.Type != "PUBLIC KEY" {
        return nil, errors.New("invalid public key format")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    return pub.(*rsa.PublicKey), nil
}

// VerifySignature verifies a signature against the challenge using the public key
func VerifySignature(pubKey *rsa.PublicKey, challenge, signature string) bool {
    hash := sha256.Sum256([]byte(challenge))
    signatureBytes, err := base64.StdEncoding.DecodeString(signature)
    if err != nil {
        return false
    }

    // Verify the signature using RSA-PSS
    err = rsa.VerifyPSS(pubKey, crypto.SHA256, hash[:], signatureBytes, nil)
    return err == nil
}
