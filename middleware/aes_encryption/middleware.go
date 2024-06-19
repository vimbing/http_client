package aesencryption

import (
	"crypto/aes"
	"crypto/cipher"
)

func pKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

func (m *Middleware) encrypt(body []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.Key)

	if err != nil {
		return []byte{}, err
	}

	var encrypted []byte

	mode := cipher.NewCBCEncrypter(block, []byte(m.Iv))
	mode.CryptBlocks(encrypted, body)

	return encrypted, nil
}

func (m *Middleware) decrypt(body []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.Key)

	if err != nil {
		return []byte{}, err
	}

	var decrypted []byte

	mode := cipher.NewCBCDecrypter(block, []byte(m.Iv))
	mode.CryptBlocks(decrypted, body)

	return pKCS5UnPadding(decrypted), nil
}
