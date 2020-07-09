package rsa

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/fuloge/basework/api"
	cfg "github.com/fuloge/basework/configs"
	"io/ioutil"
)

var (
	publicKey  []byte
	privateKey []byte
	pub        *rsa.PublicKey
	priv       *rsa.PrivateKey
)

func init() {
	var err error
	privateKey, err = ioutil.ReadFile(cfg.EnvConfig.Authkey.PrivateKey)
	if err != nil {
		panic("read private pem fail")
	}

	publicKey, err = ioutil.ReadFile(cfg.EnvConfig.Authkey.Publickey)
	if err != nil {
		panic("read public pem fail")
	}

	block, _ := pem.Decode(publicKey)
	if block == nil {
		panic("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err.Error())
	}

	pub = pubInterface.(*rsa.PublicKey)

	block, _ = pem.Decode(privateKey)
	if block == nil {
		panic("private key error")
	}

	priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err.Error())
	}
}

// 加密
func RsaEncrypt(origData []byte) ([]byte, *api.Errno) {
	partLen := pub.N.BitLen()/8 - 11
	chunks := split([]byte(origData), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, chunk)
		if err != nil {
			fmt.Println(api.RSAEncERR, err.Error())
			return nil, api.RSAEncERR
		}

		buffer.Write(bytes)
	}

	return buffer.Bytes(), nil
}

// 解密
func RsaDecrypt(ciphertext string) (string, *api.Errno) {
	partLen := pub.N.BitLen() / 8
	chunks := split([]byte([]byte(ciphertext)), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, priv, chunk)
		if err != nil {
			fmt.Println(api.RSADecERR, err.Error())
			return "", api.RSADecERR
		}
		buffer.Write(decrypted)
	}

	return buffer.String(), nil
}

// 数据加签
func Sign(data string) (string, error) {
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(sign), err
}

// 数据验签
func Verify(data string, sign string) error {
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed, decodedSign)
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
