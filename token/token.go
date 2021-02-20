package token

import (
	"auth/crypto"
	"encoding/hex"
	"fmt"
	"sort"
)

type Token struct {
	ID            uint64            `json:"ID"`
	CreateTime    uint64            `json:"CreateTime"`
	ExpireSeconds uint32            `json:"ExpireSeconds"`
	Info          map[string]string `json:"Info"`
	Signature     string            `json:"Signature"`
}

func signature(token *Token, key []byte) string {
	tokenStr := fmt.Sprintf("CreateTime=%d&ExpireSeconds=%d\n", token.CreateTime, token.ExpireSeconds)

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

func (token *Token) signature(key []byte) {
	token.Signature = signature(token, key)
}

func (token *Token) check(key []byte) bool {
	return token.Signature == signature(token, key)
}
