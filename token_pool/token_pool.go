package token_pool

import (
	"auth/random"
	"auth/token"
	"container/list"
	"fmt"
	"time"
)

type TokenFlag struct {
	ID         token.TokenID
	ExpireTime uint64
}

type Strategy func(*token.Token, *TokenPool) bool

type TokenPool struct {
	DefaultExpireSeconds uint32
	DefaultKey           []byte

	PushStrategy  Strategy
	CheckStrategy Strategy

	TokenFlags *list.List
	IndexID    map[token.TokenID]*list.Element

	ExpiredTokenFlags *list.List
	ExpiredIDs        map[token.TokenID]uint64
}

type TokenPoolOption struct {
	DefaultExpireSeconds uint32
	DefaultKey           []byte
	PushStrategy         Strategy
	CheckStrategy        Strategy
}

var defaultTokenPoolOption = TokenPoolOption{
	DefaultExpireSeconds: 1200, // 20min
	DefaultKey:           []byte(random.RandomString(16)),
	PushStrategy: func(*token.Token, *TokenPool) bool {
		return true
	},
	CheckStrategy: func(token *token.Token, pool *TokenPool) bool {
		_, ok := pool.ExpiredIDs[token.ID]
		if ok {
			return false
		}

		_, ok = pool.IndexID[token.ID]
		if !ok {
			return false
		}

		return token.Check(pool.DefaultKey)
	},
}

func NewWithOption(option TokenPoolOption) *TokenPool {
	return &TokenPool{
		DefaultExpireSeconds: option.DefaultExpireSeconds,
		DefaultKey:           option.DefaultKey,

		PushStrategy:  option.PushStrategy,
		CheckStrategy: option.CheckStrategy,

		TokenFlags: list.New(),
		IndexID:    make(map[token.TokenID]*list.Element),

		ExpiredTokenFlags: list.New(),
		ExpiredIDs:        make(map[token.TokenID]uint64),
	}
}

func New() *TokenPool {
	return NewWithOption(defaultTokenPoolOption)
}

// 定时任务清理
func (pool *TokenPool) ClearExpired() {
	for pool.TokenFlags.Len() > 0 {
		element := pool.TokenFlags.Front()
		tokenFlag := element.Value.(*TokenFlag)

		if tokenFlag.ExpireTime <= uint64(time.Now().Unix()) {
			pool.TokenFlags.Remove(element)
			delete(pool.IndexID, tokenFlag.ID)
		} else {
			break
		}
	}

	for pool.ExpiredTokenFlags.Len() > 0 {
		element := pool.TokenFlags.Front()
		tokenFlag := element.Value.(*TokenFlag)

		if tokenFlag.ExpireTime <= uint64(time.Now().Unix()) {
			pool.ExpiredTokenFlags.Remove(element)
			delete(pool.ExpiredIDs, tokenFlag.ID)
		} else {
			break
		}
	}
}

func pushSatisfied(value interface{}, list *list.List, condition func(*list.Element, interface{}) bool) *list.Element {
	if list.Len() == 0 {
		return list.PushFront(value)
	}

	element := list.Back()

	for {
		if condition(element, value) {
			break
		}
		element = element.Prev()
	}

	if element == nil {
		return list.PushFront(value)
	} else {
		return list.InsertAfter(value, element)
	}
}

func (pool *TokenPool) push(token *token.Token) error {
	pool.ClearExpired()

	if !pool.PushStrategy(token, pool) {
		return fmt.Errorf("push strategy refused token(id: %d)", token.ID)
	}

	tokenFlag := &TokenFlag{
		ID:         token.ID,
		ExpireTime: token.ExpireTime,
	}

	condition := func(element *list.Element, value interface{}) bool {
		return element == nil || element.Value.(*TokenFlag).ExpireTime <= value.(*TokenFlag).ExpireTime
	}

	element := pushSatisfied(tokenFlag, pool.TokenFlags, condition)
	pool.IndexID[tokenFlag.ID] = element

	return nil
}

func (pool *TokenPool) PushWithSelfGenerate(token *token.Token) error {
	return pool.push(token)
}

func (pool *TokenPool) GetTokenFlag(ID token.TokenID) *TokenFlag {
	element, ok := pool.IndexID[ID]
	if !ok {
		return nil
	}
	return element.Value.(*TokenFlag)
}

func (pool *TokenPool) RemoveToken(ID token.TokenID) error {
	pool.ClearExpired()

	element, ok := pool.IndexID[ID]
	if !ok {
		return fmt.Errorf("token(ID: %d) is expired or not exists, cannot be remove", ID)
	}

	pool.TokenFlags.Remove(element)
	delete(pool.IndexID, ID)

	return nil
}

func (pool *TokenPool) Expired(tokenFlag TokenFlag) {
	pool.ClearExpired()

	pool.ExpiredIDs[tokenFlag.ID] = tokenFlag.ExpireTime

	condition := func(element *list.Element, value interface{}) bool {
		return element == nil || element.Value.(*TokenFlag).ExpireTime <= value.(*TokenFlag).ExpireTime
	}
	pushSatisfied(tokenFlag, pool.ExpiredTokenFlags, condition)
}

func (pool *TokenPool) Check(token *token.Token) error {
	pool.ClearExpired()

	if !pool.CheckStrategy(token, pool) {
		return fmt.Errorf("check strategy check token(id: %d) failed", token.ID)
	}

	return nil
}

func (pool *TokenPool) generateToken(info token.TokenInfo, expireSeconds uint32, key []byte, copyInfo bool) *token.Token {
	var newToken *token.Token
	if copyInfo {
		newToken = token.GenerateToken(info, expireSeconds, key)
	} else {
		newToken = token.GenerateTokenNoCopyInfo(info, expireSeconds, key)
	}

	pool.push(newToken)
	return newToken
}

func (t *TokenPool) GenerateTokenNoCopyInfo(info token.TokenInfo) *token.Token {
	return t.generateToken(info, t.DefaultExpireSeconds, t.DefaultKey, false)
}

func (t *TokenPool) GenerateToken(info token.TokenInfo) *token.Token {
	return t.generateToken(info, t.DefaultExpireSeconds, t.DefaultKey, true)
}

func (t *TokenPool) GenerateTokenID() *token.Token {
	return t.generateToken(nil, t.DefaultExpireSeconds, t.DefaultKey, false)
}
