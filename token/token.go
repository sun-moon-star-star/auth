package token

import (
    "encoding/json"
    "auth/crypt"
)

type TokenContent struct {
    CreateTime uint64 `json:"CreateTime"`
    ID string `json:"ID"` 
    Message string `json:"Message"`
}

func GenerateToken(TokenContent *TokenContent, key []byte) ([]byte, error) {
    tokenContentStr, err := json.Marshal(TokenContent)
    if err != nil {
        return nil, err
    }
    var tokenStr []byte
    tokenStr, err = crypt.AesEncrypt(tokenContentStr, key)
    return tokenStr, nil
}

func CheckToken(tokenStr []byte , key []byte) (*TokenContent, error) {
    tokenContentStr, err := crypt.AesDecrypt(tokenStr, key)
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