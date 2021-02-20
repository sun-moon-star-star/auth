package token_pool

import (
	token "auth/token"
	list "container/list"
	"time"
)

type TokenFlag struct {
	ID            uint64
	ExpireSeconds uint32
}

type TokenPool struct {
	DefaultExpireSeconds uint32
	DefaultKey           []byte
	TokenFlags           *list.List
	IndexID              map[uint64]*list.Element
}

func New(DefaultExpireSeconds uint32, DefaultKey []byte) *TokenPool {
	return &TokenPool{
		DefaultExpireSeconds: DefaultExpireSeconds,
		DefaultKey:           DefaultKey,
		TokenFlags:           list.New(),
		IndexID:              make(map[uint64]*list.Element),
	}
}

func (t TokenPool) Push(token *token.Token) {
	tf := &TokenFlag{
		ID:            token.ID,
		ExpireSeconds: token.ExpireSeconds,
	}
	element := t.TokenFlags.PushBack(tf)
	t.IndexID[tf.ID] = element
}

func generateTokenNoCopyInfo(info map[string]string, expireSeconds uint32, key []byte) *token.Token {
	token := &token.Token{
		ID:            token.RandomUint64(),
		CreateTime:    uint64(time.Now().Unix()),
		ExpireSeconds: expireSeconds,
		Info:          info,
	}
	token.Sign(key)
	return token
}

func (t *TokenPool) GenerateTokenNoCopyInfo(info map[string]string) *token.Token {
	token := generateTokenNoCopyInfo(info, t.DefaultExpireSeconds, t.DefaultKey)
	t.Push(token)
	return token
}

func copyInfo(info map[string]string) map[string]string {
	infoNew := make(map[string]string)

	for k, v := range info {
		infoNew[k] = v
	}

	return infoNew
}

func (t *TokenPool) GenerateToken(info map[string]string) *token.Token {
	return t.GenerateTokenNoCopyInfo(copyInfo(info))
}

func (t *TokenPool) GenerateTokenNoCopyInfoWithKey(info map[string]string, key []byte) *token.Token {
	token := generateTokenNoCopyInfo(info, t.DefaultExpireSeconds, key)
	t.Push(token)
	return token
}

func (t *TokenPool) GenerateTokenWithKey(info map[string]string, key []byte) *token.Token {
	return t.GenerateTokenNoCopyInfoWithKey(copyInfo(info), key)
}

func (t *TokenPool) GenerateTokenNoCopyInfoWithExpireSeconds(info map[string]string, expireSeconds uint32) *token.Token {
	token := generateTokenNoCopyInfo(info, expireSeconds, t.DefaultKey)
	t.Push(token)
	return token
}

func (t *TokenPool) GenerateTokenWithExpireSeconds(info map[string]string, expireSeconds uint32) *token.Token {
	return t.GenerateTokenNoCopyInfoWithExpireSeconds(copyInfo(info), expireSeconds)
}

func (t *TokenPool) GenerateTokenNoCopyInfoWithExpireSecondsAndKey(info map[string]string, expireSeconds uint32, key []byte) *token.Token {
	token := generateTokenNoCopyInfo(info, expireSeconds, key)
	t.Push(token)
	return token
}

func (t *TokenPool) GenerateTokenWithExpireSecondsAndKey(info map[string]string, expireSeconds uint32, key []byte) *token.Token {
	return t.GenerateTokenNoCopyInfoWithExpireSecondsAndKey(copyInfo(info), expireSeconds, key)
}
