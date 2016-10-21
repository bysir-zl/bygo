package util

import (
    "crypto/des"
    "crypto/cipher"
    "bytes"
)


func DesEncrypt(origData, key []byte) ([]byte, error) {
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    origData = PKCS5Padding(origData, block.BlockSize())
    // origData = ZeroPadding(origData, block.BlockSize())
    blockMode := cipher.NewCBCEncrypter(block, key)
    crypted := make([]byte, len(origData))
    // 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
    // crypted := origData
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}

//des解密

func DesDecrypt(crypted, key []byte) ([]byte, error) {
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockMode := cipher.NewCBCDecrypter(block, key)
    //origData := make([]byte, len(crypted))
    origData := crypted
    blockMode.CryptBlocks(origData, crypted)
    //origData = PKCS5UnPadding(origData)

    origData = ZeroUnPadding(origData)
    return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
