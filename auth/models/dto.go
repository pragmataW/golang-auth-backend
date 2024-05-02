package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

type ReqUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ReqUser
	VerifyCode int
	CreatedAt  time.Time
	IsVerified bool
}

type Encryptor struct {
	Key []byte
}

func NewEncryptor(key string) *Encryptor {
	return &Encryptor{Key: []byte(key)}
}

func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	bPlainText := []byte(plaintext)

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipherText := make([]byte, len(bPlainText))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText, bPlainText)

	result := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func (e *Encryptor) Decrypt(encryptedText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
