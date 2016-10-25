package util

import (
    "crypto/cipher"
    "crypto/aes"
)

//aes加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
    var iv = []byte(key)[:aes.BlockSize]
    encrypted := make([]byte, len(origData))
    aesBlockEncrypter, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
    aesEncrypter.XORKeyStream(encrypted, origData)
    return encrypted, nil
}

//aes解密

func AesDecrypt(crypted, key []byte) ([]byte, error) {
    var iv = []byte(key)[:aes.BlockSize]
    decrypted := make([]byte, len(crypted))
    var aesBlockDecrypter cipher.Block
    aesBlockDecrypter, err := aes.NewCipher([]byte(key))
    if err != nil {
        return []byte{}, err
    }
    aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
    aesDecrypter.XORKeyStream(decrypted, crypted)
    return decrypted, nil
}
