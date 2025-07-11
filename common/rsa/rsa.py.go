package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// 生成 RSA 密钥对
func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// 加密数据
func encrypt(publicKey *rsa.PublicKey, plaintext []byte) (string, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 解密数据
func decrypt(privateKey *rsa.PrivateKey, ciphertext string) (string, error) {
	// 解码 Base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 解密
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// 将公钥转换为 PEM 格式
func exportPublicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})
	return string(pubPEM), nil
}

// 将私钥转换为 PEM 格式
func exportPrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	return string(privPEM)
}
