package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// ------------Private Key conversion
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(priv)
}

// PrivateKeyToString: return base64 encoded version
func PrivateKeyToString(priv *rsa.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(PrivateKeyToBytes(priv))
}

func Base64ToPrivateKey(priv string) (*rsa.PrivateKey, error) {
	privDec, err := base64.StdEncoding.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(privDec)
}

// ------------ Public Key conversion
// Base64ToPublicKey: base64 encoded to public key
func Base64ToPublicKey(pub string) (*rsa.PublicKey, error) {
	pubDec, err := base64.StdEncoding.DecodeString(pub)
	if err != nil {
		return nil, err
	}
	pb, err := x509.ParsePKIXPublicKey(pubDec)
	if err != nil {
		return nil, err
	}

	if pub, ok := pb.(*rsa.PublicKey); ok {
		return pub, nil
	}

	return nil, fmt.Errorf("Expected *rsa.PublicKey, got %T", pb)
}

// Convert public key to bytes (pem)
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

// PrivateKeyToString: return base64 encoded version
func PublicKeyToString(pub *rsa.PublicKey) (string, error) {
	pubByte, err := PublicKeyToBytes(pub)
	if err != nil {
		return "", err
	}
	return EncodeBytesToBase64String(pubByte), nil
}

func GenerateRSAKeyPair(keySizeOptional ...int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	keySize := 2048
	if len(keySizeOptional) > 0 {
		keySize = keySizeOptional[0]
	}
	privkey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

// --------------ENCRYPT

// EncryptRSAWithPublicKey: Encrypt message (return as base64 String) using public key (base64 string)
func EncryptRSAWithPublicKey(msg string, pubB64 string) (string, error) {
	pub, err := Base64ToPublicKey(pubB64)
	if err != nil {
		return "", err
	}
	hash := sha512.New()

	// Chunk the message into smaller parts
	var chunkSize = pub.N.BitLen()/8 - 2*hash.Size() - 2
	var result []byte
	chunks := chunkBy[byte]([]byte(msg), chunkSize)
	for _, chunk := range chunks {
		ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, chunk, nil)
		if err != nil {
			return "", err
		}
		result = append(result, ciphertext...)
	}

	return EncodeBytesToBase64String(result), nil
}

// --------------DECRYPT

// DecryptRSAWithPrivateKey: Decrypt message using private key (base64 string)
func DecryptRSAWithPrivateKey(ciphertextB64 string, privB64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}
	priv, err := Base64ToPrivateKey(privB64)
	if err != nil {
		return "", err
	}
	hash := sha512.New()
	dec_msg := []byte("")

	for _, chnk := range chunkBy[byte](ciphertext, priv.N.BitLen()/8) {
		plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, chnk, nil)
		if err != nil {
			return "", err
		}
		dec_msg = append(dec_msg, plaintext...)
	}

	return string(dec_msg), nil
}
