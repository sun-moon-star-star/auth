package token

import (
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/sun-moon-star-star/auth/crypto"
	"github.com/sun-moon-star-star/auth/random"
)

type TokenID uint64

var defaultFillChar byte = '-'

type TokenInfo map[string]interface{}

type Token struct {
	ID         TokenID   `json:"ID"`
	CreateTime uint64    `json:"CreateTime"` // unixnano
	ExpireTime uint64    `json:"ExpireTime"` // unixnano
	Info       TokenInfo `json:"Info"`       // _token_* keeps
	Signature  string    `json:"Signature"`
}

func GenerateTokenID() TokenID {
	return TokenID(random.RandomUint64())
}

func copyInfo(info TokenInfo) TokenInfo {
	if info == nil {
		return nil
	}

	infoNew := make(TokenInfo)

	for k, v := range info {
		infoNew[k] = v
	}

	return infoNew
}

func generateTokenNoCopyInfo(info TokenInfo, expireUnixNano uint64, key []byte) *Token {
	createTime := uint64(time.Now().UnixNano())

	newToken := &Token{
		ID:         GenerateTokenID(),
		CreateTime: createTime,
		ExpireTime: createTime + expireUnixNano,
		Info:       info,
	}

	newToken.Sign(key)

	return newToken
}

func GenerateTokenNoCopyInfo(info TokenInfo, expireUnixNano uint64, key []byte) *Token {
	return generateTokenNoCopyInfo(info, expireUnixNano, key)
}

func GenerateToken(info TokenInfo, expireUnixNano uint64, key []byte) *Token {
	return generateTokenNoCopyInfo(copyInfo(info), expireUnixNano, key)
}

func TokenString(token *Token) string {
	tokenStr := fmt.Sprintf("CreateTime=%d&ExpireTime=%d\n", token.CreateTime, token.ExpireTime)
	var keys []string
	for k := range token.Info {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tokenStr += fmt.Sprintf("&%s=%+v", k, token.Info[k])
	}

	return tokenStr
}

func Sign(token *Token, key []byte) string {
	tokenStr := fmt.Sprintf("CreateTime=%d&ExpireTime=%d\n", token.CreateTime, token.ExpireTime)

	var keys []string
	for k := range token.Info {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tokenStr += fmt.Sprintf("&%s=%+v", k, token.Info[k])
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
	return token.ExpireTime >= uint64(time.Now().UnixNano())
}

func (token *Token) Check(key []byte) bool {
	return token.CheckTime() && token.CheckSign(key)
}
