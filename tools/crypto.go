package tools

import (
	"crypto/aes"
	"crypto/cipher"
)

/*
   Creation Time: 2019 - Oct - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// AES256CCMEncrypt encrypts the msg according with key and iv
func AES256GCMEncrypt(key, iv []byte, msg []byte) ([]byte, error) {
	var block cipher.Block
	if b, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		block = b
	}
	var encrypted []byte
	if aesGCM, err := cipher.NewGCM(block); err != nil {
		return nil, err
	} else {
		encrypted = aesGCM.Seal(msg[:0], iv, msg, nil)
	}
	return encrypted, nil
}

// AES256GCMDecrypt decrypts the msg according with key and iv
func AES256GCMDecrypt(key, iv []byte, msg []byte) ([]byte, error) {
	var block cipher.Block
	if b, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		block = b
	}
	var decrypted []byte
	if aesGCM, err := cipher.NewGCM(block); err != nil {
		return nil, err
	} else {
		decrypted, err = aesGCM.Open(nil, iv, msg, nil)
		if err != nil {
			return nil, err
		}
	}
	return decrypted, nil
}
