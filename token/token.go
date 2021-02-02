package token

import (
    "bytes"
    "encoding/json"
    "crypto/aes"
    "crypto/cipher"
)

type TokenContent struct {
    CreateTime uint64 `json:"CreateTime"`
    ID string `json:"ID"` 
    Message string `json:"Message"`
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    origData = pkcs7Padding(origData, blockSize)
    blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
    crypted := make([]byte, len(origData))
    blockMode.CryptBlocks(crypted, origData)
    return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
    origData := make([]byte, len(crypted))
    blockMode.CryptBlocks(origData, crypted)
    origData = pkcs7UnPadding(origData)
    return origData, nil
}

func GenerateToken(TokenContent *TokenContent, key []byte) ([]byte, error) {
    tokenContentStr, err := json.Marshal(TokenContent)
    if err != nil {
        return nil, err
    }
    var tokenStr []byte
    tokenStr, err = AesEncrypt(tokenContentStr, key)
    return tokenStr, nil
}

func CheckToken(tokenStr []byte , key []byte) (*TokenContent, error) {
    tokenContentStr, err := AesDecrypt(tokenStr, key)
    if err != nil {
        return nil, nil
    }
    var token TokenContent 
    err = json.Unmarshal(tokenContentStr, &token)
    if err != nil {
        return nil, nil
    }
    return &token, err
}