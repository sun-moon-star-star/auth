package token_pool

import (
	"auth/token"
	"container/list"
	"fmt"
	"time"
)

type TokenFlag struct {
	ID            uint64
	CreateTime    uint64
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

func (t *TokenPool) RemoveExpire() {
	for t.TokenFlags.Len() > 0 {
		element := t.TokenFlags.Front()
		tokenFlag := element.Value.(*TokenFlag)

		if tokenFlag.CreateTime+uint64(tokenFlag.ExpireSeconds) >= uint64(time.Now().Unix()) {
			delete(t.IndexID, tokenFlag.ID)
			t.TokenFlags.Remove(element)
		} else {
			return
		}
	}
}

func (t *TokenPool) Push(token *token.Token) error {
	t.RemoveExpire()

	_, ok := t.IndexID[token.ID]
	if ok {
		return fmt.Errorf("token(ID: %d) exists in token pool", token.ID)
	}

	tokenFlag := &TokenFlag{
		ID:            token.ID,
		CreateTime:    token.CreateTime,
		ExpireSeconds: token.ExpireSeconds,
	}

	element := t.TokenFlags.PushBack(tokenFlag)
	t.IndexID[tokenFlag.ID] = element

	return nil
}

func (t *TokenPool) Remove(ID uint64) error {
	t.RemoveExpire()

	element, ok := t.IndexID[ID]
	if !ok {
		return fmt.Errorf("token(ID: %d) not exists in token pool, cannot be remove", ID)
	}

	t.TokenFlags.Remove(element)
	delete(t.IndexID, ID)

	return nil
}

func (t *TokenPool) Check(token *token.Token) error {
	t.RemoveExpire()

	element, ok := t.IndexID[token.ID]
	if !ok {
		return fmt.Errorf("token(ID: %d) not exists in token pool, please try to acquire new token", token.ID)
	}

	tokenFlag := element.Value.(*TokenFlag)
	if tokenFlag.CreateTime+uint64(tokenFlag.ExpireSeconds) >= uint64(time.Now().Unix()) {
		return fmt.Errorf("token(ID: %d) is expired, please try to acquire new token", token.ID)
	}

	if !token.Check(t.DefaultKey) {
		return fmt.Errorf("token(ID: %d) is not valid", token.ID)
	}

	return nil
}

func (t *TokenPool) generateTokenNoCopyInfo(info map[string]string, expireSeconds uint32, key []byte) *token.Token {
	newToken := &token.Token{
		ID:            token.RandomUint64(),
		CreateTime:    uint64(time.Now().Unix()),
		ExpireSeconds: expireSeconds,
		Info:          info,
	}

	for {
		_, ok := t.IndexID[newToken.ID]
		if !ok {
			break
		}
		newToken.ID = token.RandomUint64()
	}

	newToken.Sign(key)

	t.Push(newToken)

	return newToken
}

func copyInfo(info map[string]string) map[string]string {
	infoNew := make(map[string]string)

	for k, v := range info {
		infoNew[k] = v
	}

	return infoNew
}

func (t *TokenPool) GenerateTokenNoCopyInfo(info map[string]string) *token.Token {
	return t.generateTokenNoCopyInfo(info, t.DefaultExpireSeconds, t.DefaultKey)
}

func (t *TokenPool) GenerateToken(info map[string]string) *token.Token {
	return t.generateTokenNoCopyInfo(copyInfo(info), t.DefaultExpireSeconds, t.DefaultKey)
}
