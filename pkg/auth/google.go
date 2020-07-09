package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/fuloge/basework/api"
	"image/png"
	"strings"
	"time"
)

type GoogleAuth struct {
}

func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}

func (this *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}

func (this *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}

func (this *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}

func (this *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}

func (this *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (this *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}

func (this *GoogleAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := this.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := this.toUint32(hashParts)
	return number % 1000000
}

// 获取秘钥
func (this *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, this.un())
	return strings.ToUpper(this.base32encode(this.hmacSha1(buf.Bytes(), nil)))
}

// 获取动态码
func (this *GoogleAuth) GetCode(secret string) (string, *api.Errno) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := this.base32decode(secretUpper)
	if err != nil {
		fmt.Println(api.GoogleAuthGetErr.Message, err.Error())
		return "", api.GoogleAuthGetErr
	}
	number := this.oneTimePassword(secretKey, this.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}

// 获取动态码二维码内容
func (this *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}

// 获取动态码二维码图片地址,这里是第三方二维码api
//func (this *GoogleAuth) GetQrcodeUrl(user, secret string) string {
//	qrcode := this.GetQrcode(user, secret)
//	width := "200"
//	height := "200"
//	data := url.Values{}
//	data.Set("data", qrcode)
//	return "https://api.qrserver.com/v1/create-qr-code/?" + data.Encode() + "&size=" + width + "x" + height + "&ecc=M"
//}

// 验证动态码
func (this *GoogleAuth) VerifyCode(secret, code string) (bool, *api.Errno) {
	_code, err := this.GetCode(secret)
	if err != nil {
		fmt.Println(api.GoogleAuthVerifyErr.Message, err.Message)
		return false, api.GoogleAuthVerifyErr
	}
	return _code == code, nil
}

// 在浏览器显示二维码
func (this *GoogleAuth) CreateQRcode(user, secret string, size int) string {
	qrc, _ := qr.Encode(this.GetQrcode(user, secret), qr.M, qr.Auto)
	qrc, _ = barcode.Scale(qrc, size, size)
	b := bytes.Buffer{}
	png.Encode(&b, qrc)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(b.Bytes())
}

//func WritePng(filename string, img image.Image) {
//	file, err := os.Create(filename)
//	if err != nil {
//		log.Fatal(err)
//	}
//	err = png.Encode(file, img)
//	// err = jpeg.Encode(file, img, &jpeg.Options{100})      //图像质量值为100，是最好的图像显示
//	if err != nil {
//		log.Fatal(err)
//	}
//	file.Close()
//	log.Println(file.Name())
//}
