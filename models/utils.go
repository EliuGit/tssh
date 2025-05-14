package models

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var commonKeyBytes = []byte("thisis32byteslongsecretkey123456")

// pkcs7Padding 填充明文到AES块大小的整数倍
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pkcs7UnPadding 去除PKCS7填充
func pkcs7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("pkcs7: unpadding error, empty data")
	}
	unpadding := int(origData[length-1])
	if unpadding > length || unpadding == 0 { // 检查padding值是否合法
		return nil, errors.New("pkcs7: unpadding error, invalid padding value")
	}
	// 检查填充字节是否都相同
	for i := 0; i < unpadding; i++ {
		if origData[length-unpadding+i] != byte(unpadding) {
			return nil, errors.New("pkcs7: unpadding error, invalid padding bytes")
		}
	}
	return origData[:(length - unpadding)], nil
}

// EncryptAESCBC 使用AES CBC模式加密
// key 必须是16, 24, 或 32字节長度 (对应 AES-128, AES-192, AES-256)
// plaintext 是要加密的原始数据
// 返回值是 IV + ciphertext (IV被预置在密文前)
func EncryptAESCBC(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	blockSize := block.BlockSize()
	paddedPlaintext := pkcs7Padding(plaintext, blockSize)
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)
	return append(iv, ciphertext...), nil
}

// DecryptAESCBC 使用AES CBC模式解密
// key 必须是16, 24, 或 32字节長度 (对应 AES-128, AES-192, AES-256)
// ciphertextWithIV 是包含IV的密文 (IV在最前面)
// 返回原始明文
func DecryptAESCBC(key []byte, ciphertextWithIV []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	blockSize := block.BlockSize()
	if len(ciphertextWithIV) < blockSize {
		return nil, errors.New("ciphertext is too short (missing IV)")
	}
	iv := ciphertextWithIV[:blockSize]
	ciphertext := ciphertextWithIV[blockSize:]
	if len(ciphertext)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	decryptedText := make([]byte, len(ciphertext))
	mode.CryptBlocks(decryptedText, ciphertext)
	unpaddedText, err := pkcs7UnPadding(decryptedText)
	if err != nil {
		return nil, fmt.Errorf("failed to unpad: %w", err)
	}
	return unpaddedText, nil
}

func EncryptString(plaintext string) (string, error) {
	plaintextBytes := []byte(plaintext)
	encryptedBytes, err := EncryptAESCBC(commonKeyBytes, plaintextBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
func DecryptString(ciphertext string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := DecryptAESCBC(commonKeyBytes, ciphertextBytes)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}
