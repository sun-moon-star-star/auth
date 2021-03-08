package token_scope

import "auth/token"

const (
	READ = 1 << iota
	WRITE
	EXEC
)

func New(topic string, authority uint8, expireSeconds uint32, key []byte) *token.Token {
	info := map[string]interface{}{
		"_token_topic":     topic,
		"_token_authority": authority,
	}
	return token.GenerateTokenNoCopyInfo(info, expireSeconds, key)
}
