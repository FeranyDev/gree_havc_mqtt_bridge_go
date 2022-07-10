package app

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

var defaultKey = "a3K8Bx%2r8Y7#xDh"

type DeviceInfo struct {
	T       string `json:"t"`
	Cid     string `json:"cid"`
	Bc      string `json:"bc"`
	Brand   string `json:"brand"`
	Catalog string `json:"catalog"`
	Mac     string `json:"mac"`
	Mid     string `json:"mid"`
	Model   string `json:"model"`
	Name    string `json:"name"`
	Lock    int    `json:"lock"`
	Series  string `json:"series"`
	Vender  string `json:"vender"`
	Ver     string `json:"ver"`
	Key     string `json:"key"`
}

func generateMd5Key(str string) []byte {
	checksum := md5.Sum([]byte(str))
	return checksum[0:]
}

func aesEncryptECB(plaintext string, key string) string {
	origData := []byte(plaintext)
	checksum := generateMd5Key(key)
	cipher, _ := aes.NewCipher(checksum)
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return hex.EncodeToString(encrypted)
}
func aesDecryptECB(ciphertext string, key string) (string, error) {
	checksum := generateMd5Key(key)
	cipher, _ := aes.NewCipher(checksum)
	encrypted, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, len(encrypted))
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return string(decrypted[:trim]), nil
}

func decrypt(input, key string) DeviceInfo {
	if key == "" {
		key = defaultKey
	}
	plaintext := input
	deviceInfo := DeviceInfo{}
	fmt.Printf("原    文：%s\n", plaintext)
	encrypted := aesEncryptECB(plaintext, key)
	fmt.Printf("加密结果：%s\n", encrypted)
	err := json.Unmarshal([]byte(encrypted), &deviceInfo)
	if err != nil {
		log.Errorf("解析json失败：%s", err)
	}
	return deviceInfo

}
func encrypt(output interface{}, key string) string {
	if key == "" {
		key = defaultKey
	}
	marshal, err := json.Marshal(output)
	if err != nil {
		log.Errorf("json序列化失败：%s", err)
		return ""
	}
	decrypted, err := aesDecryptECB(string(marshal), key)
	fmt.Printf("解密结果：%s\n", decrypted)
	if err != nil {
		panic(err)
	}
	return decrypted
}
