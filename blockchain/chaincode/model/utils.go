package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	math_rand "math/rand"
	"strings"

	"github.com/tjfoc/gmsm/sm2"
)

var CHARS = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

/*RandAllString  生成随机字符串([a~zA~Z0~9])
  lenNum 长度
*/
func RandAllString(lenNum int) string {

	str := strings.Builder{}
	length := len(CHARS)
	for i := 0; i < lenNum; i++ {
		l := CHARS[math_rand.Intn(length)]
		str.WriteString(l)
	}
	return str.String()
}

func SM2Encrypt(msg string, keyGenSeed int64, encryptSeed int64) (string, error) {
	math_rand.Seed(keyGenSeed)
	keyStr := RandAllString(40)
	keyReader := bytes.NewReader([]byte(keyStr + keyStr))
	priv, err := sm2.GenerateKey(keyReader) // 生成密钥对
	if err != nil {
		fmt.Println(msg, err)
		return "", err
	}

	pub := &priv.PublicKey
	math_rand.Seed(encryptSeed)
	encStr := RandAllString(40)
	encReader := bytes.NewReader([]byte(encStr + encStr))
	ciphertxt, err := pub.EncryptAsn1([]byte(msg), encReader)
	fmt.Printf("ok:%x", string(ciphertxt))
	if err != nil {
		fmt.Println(msg, err)
		return "", err
	}
	return fmt.Sprintf("%x", string(ciphertxt)), nil
}

func SM2Decrypt(priv *sm2.PrivateKey, cipher string) string {
	msg, _ := hex.DecodeString(cipher)
	plaintxt, err := priv.DecryptAsn1(msg)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(plaintxt)
}

func SM2Sign(msg string, keyGenSeed int64, encryptSeed int64) (string, error) {
	math_rand.Seed(keyGenSeed)
	keyStr := RandAllString(40)
	keyReader := bytes.NewReader([]byte(keyStr + keyStr))
	priv, err := sm2.GenerateKey(keyReader) // 生成密钥对
	if err != nil {
		fmt.Println(msg, err)
		return "", err
	}
	math_rand.Seed(encryptSeed)
	encStr := RandAllString(40)
	encReader := bytes.NewReader([]byte(encStr + encStr))
	sign, err := priv.Sign(encReader, []byte(msg), nil)
	fmt.Printf("ok:%x", string(sign))
	if err != nil {
		fmt.Println(msg, err)
		return "", err
	}
	return fmt.Sprintf("%x", string(sign)), nil
}

func SM2Verify(msg string, sign string, keyGenSeed int64) (bool, error) {
	signByte, _ := hex.DecodeString(sign)
	math_rand.Seed(keyGenSeed)
	keyStr := RandAllString(40)
	keyReader := bytes.NewReader([]byte(keyStr + keyStr))
	priv, err := sm2.GenerateKey(keyReader) // 生成密钥对
	if err != nil {
		fmt.Println(msg, err)
		return false, err
	}
	ok := priv.Verify([]byte(msg), signByte)
	return ok, nil
}

func GetHash(msg string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	res := h.Sum(nil)
	return hex.EncodeToString(res), nil
}

func CheckHash(msg string, msgHash string) bool {
	msgHashCurrent, err := GetHash(msg)
	if err != nil {
		log.Fatal("Check Hash False", err)
		return false
	}

	return msgHash == msgHashCurrent
}

func DeleteSlice(a []string, elem string) []string {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}
