package token

import (
	"auth/crypto"
	"auth/random"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

type TokenID uint64

type Token struct {
	ID         TokenID           `json:"ID"`
	CreateTime uint64            `json:"CreateTime"`
	ExpireTime uint64            `json:"ExpireTime"`
	Info       map[string]string `json:"Info"`
	Signature  string            `json:"Signature"`
}

func GenerateTokenID() TokenID {
	return TokenID(random.RandomUint64())
}

func copyInfo(info map[string]string) map[string]string {
	if info == nil {
		return nil
	}

	infoNew := make(map[string]string)

	for k, v := range info {
		infoNew[k] = v
	}

	return infoNew
}

func generateTokenNoCopyInfo(info map[string]string, expireSeconds uint32, key []byte) *Token {
	createTime := uint64(time.Now().Unix())

	newToken := &Token{
		ID:         GenerateTokenID(),
		CreateTime: createTime,
		ExpireTime: createTime + uint64(expireSeconds),
		Info:       info,
	}

	newToken.Sign(key)

	return newToken
}

func GenerateTokenNoCopyInfo(info map[string]string, expireSeconds uint32, key []byte) *Token {
	return generateTokenNoCopyInfo(info, expireSeconds, key)
}

func GenerateToken(info map[string]string, expireSeconds uint32, key []byte) *Token {
	return generateTokenNoCopyInfo(copyInfo(info), expireSeconds, key)
}

func Sign(token *Token, key []byte) string {
	tokenStr := fmt.Sprintf("CreateTime=%d&ExpireTime=%d\n", token.CreateTime, token.ExpireTime)

	var keys []string
	for k := range token.Info {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tokenStr += fmt.Sprintf("&%s=%s", k, token.Info[k])
	}

	signature := crypto.Hmac([]byte(tokenStr), key)

	return hex.EncodeToString(signature)
}

func (token *Token) Sign(key []byte) {
	token.Signature = Sign(token, key)
}

func (token *Token) CheckSign(key []byte) bool {
	return token.Signature == Sign(token, key)
}

func (token *Token) CheckTime() bool {
	return token.ExpireTime > uint64(time.Now().Unix())
}

func (token *Token) Check(key []byte) bool {
	return token.CheckTime() && token.CheckSign(key)
}
