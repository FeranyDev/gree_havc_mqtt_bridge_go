package gree

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

var defaultKey = "a3K8Bx%2r8Y7#xDh"

type DeviceInfo struct {
	T       string   `json:"t"`
	Cid     string   `json:"cid"`
	Bc      string   `json:"bc"`
	Brand   string   `json:"brand"`
	Catalog string   `json:"catalog"`
	Mac     string   `json:"mac"`
	Mid     string   `json:"mid"`
	Model   string   `json:"model"`
	Name    string   `json:"name"`
	Lock    int      `json:"lock"`
	Series  string   `json:"series"`
	Vender  string   `json:"vender"`
	Ver     string   `json:"ver"`
	Key     string   `json:"key"`
	Dat     []int    `json:"dat"`
	Cols    []string `json:"cols"`
	Val     []int    `json:"val"`
	P       []int    `json:"p"`
	R       int      `json:"r"`
	Opt     []string `json:"opt"`
}

func Decrypt(input UDPInfo, key string) DeviceInfo {
	if key == "" {
		key = defaultKey
	}

	deviceInfo := DeviceInfo{}
	decodeString, err := base64.StdEncoding.DecodeString(input.Pack)
	if err != nil {
		log.Errorf("[BASE64] base64解码失败：%s", err)
		return DeviceInfo{}
	}
	decrypted, err := ECBDecrypt(decodeString, []byte(key))
	if err != nil {
		log.Errorf("[KEY] 解密失败：%s", err)
		return deviceInfo
	}
	err = json.Unmarshal(decrypted, &deviceInfo)
	if err != nil {
		log.Errorf("[JSON] json失败：%s", err)
	}
	return deviceInfo

}

func Encrypt(output interface{}, key string) string {
	if key == "" {
		key = defaultKey
	}
	marshal, err := json.Marshal(output)
	if err != nil {
		log.Errorf("[JSON] json序列化失败：%s", err)
		return ""
	}
	encrypted, _ := ECBEncrypt(marshal, []byte(key))
	return base64.StdEncoding.EncodeToString(encrypted)
}

func ECBDecrypt(crypted, key []byte) ([]byte, error) {
	if !validKey(key) {
		return nil, fmt.Errorf("秘钥长度错误,当前传入长度为 %d", len(key))
	}
	if len(crypted) < 1 {
		return nil, fmt.Errorf("源数据长度不能为0")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(crypted)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("源数据长度必须是 %d 的整数倍，当前长度为：%d", block.BlockSize(), len(crypted))
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())

	for index := 0; index < len(crypted); index += block.BlockSize() {
		block.Decrypt(tmpData, crypted[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	dst, err = PKCS5UnPadding(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func ECBEncrypt(src, key []byte) ([]byte, error) {
	if !validKey(key) {
		return nil, fmt.Errorf("秘钥长度错误, 当前传入长度为 %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src) < 1 {
		return nil, fmt.Errorf("源数据长度不能为0")
	}
	src = PKCS5Padding(src, block.BlockSize())
	if len(src)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("源数据长度必须是 %d 的整数倍，当前长度为：%d", block.BlockSize(), len(src))
	}
	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(src); index += block.BlockSize() {
		block.Encrypt(tmpData, src[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	return dst, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])

	if length < unpadding {
		return nil, fmt.Errorf("invalid unpadding length")
	}
	return origData[:(length - unpadding)], nil
}

func validKey(key []byte) bool {
	k := len(key)
	switch k {
	default:
		return false
	case 16, 24, 32:
		return true
	}
}
