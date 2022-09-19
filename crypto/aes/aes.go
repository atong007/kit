package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type AesGCM struct {
	encKey []byte
}

func NewAesGCM(encKey string) (*AesGCM, error) {
	key, err := hex.DecodeString(encKey)
	if err != nil {
		return nil, err
	}
	return &AesGCM{
		encKey: key,
	}, nil
}

func (a *AesGCM) EncryptedCode(code string) (enCode string, err error) {
	encKeyLen := len(a.encKey)

	if encKeyLen < 16 {
		err = errors.New("the key must be 16 bytes long")
		return
	}

	encKeySized := a.encKey
	if encKeyLen > 16 {
		encKeySized = encKeySized[:16]
	}

	block, err := aes.NewCipher(encKeySized)
	if err != nil {
		return
	}

	nonceSize := 12
	nonce := make([]byte, nonceSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	aesGcm, err := cipher.NewGCMWithNonceSize(block, nonceSize)
	if err != nil {
		return
	}

	plainText := []byte(code)
	cipherText := aesGcm.Seal(nil, nonce, plainText, nil)

	nonceLen := len(nonce)
	output := make([]byte, 1+nonceLen+len(cipherText))
	i := 0
	output[i] = byte(nonceLen)
	i++
	copyWithRange(cipherText, 0, output, i, nonceLen)
	i += nonceLen
	copyWithRange(nonce, 0, output, i, nonceLen)
	i += nonceLen

	copyWithRange(cipherText, nonceLen, output, i, len(cipherText)-nonceLen)
	return hex.EncodeToString(output), nil
}

// copyWithRange from pkg into a dest byte array
func copyWithRange(src []byte, srcI int, dest []byte, destI int, copyLen int) {
	srcI2 := srcI + copyLen
	copy(dest[destI:], src[srcI:srcI2])
}

func (a *AesGCM) DecryptedCode(enCode string) (code string, err error) {
	input, err := hex.DecodeString(enCode)
	if err != nil {
		return
	}

	encKeyLen := len(a.encKey)
	if encKeyLen < 16 {
		err = fmt.Errorf("the key must be 16 bytes long")
		return
	}

	encKeySized := a.encKey
	if encKeyLen > 16 {
		encKeySized = a.encKey[:16]
	}

	i := 0
	nonceLen := int(input[i])
	i++
	i += nonceLen

	if nonceLen != 12 {
		err = errors.New("nonce length is not correct")
		return
	}

	iv := make([]byte, nonceLen)
	copyWithRange(input, i, iv, 0, nonceLen)
	i += nonceLen

	cipherTextLen := len(input) - i + nonceLen
	cipherText := make([]byte, cipherTextLen)

	copyWithRange(input, 1, cipherText, 0, nonceLen)
	copyWithRange(input, i, cipherText, nonceLen, cipherTextLen-nonceLen)

	c, err := aes.NewCipher(encKeySized)
	if err != nil {
		return
	}

	dec, err := cipher.NewGCMWithNonceSize(c, nonceLen)
	if err != nil {
		return
	}

	output, err := dec.Open(nil, iv, cipherText, nil)
	if err != nil {
		return
	}

	return string(output), nil
}
