package token_pool

import (
	"auth/token"
	"container/list"
	"fmt"
	"time"
)

type TokenFlag struct {
	ID         token.TokenID
	ExpireTime uint64
}

type TokenPool struct {
	DefaultExpireSeconds uint32
	DefaultKey           []byte
	TokenFlags           *list.List
	IndexID              map[token.TokenID]*list.Element
}

func New(DefaultExpireSeconds uint32, DefaultKey []byte) *TokenPool {
	return &TokenPool{
		DefaultExpireSeconds: DefaultExpireSeconds,
		DefaultKey:           DefaultKey,
		TokenFlags:           list.New(),
		IndexID:              make(map[token.TokenID]*list.Element),
	}
}

// 定时任务清理
func (t *TokenPool) RemoveExpire() {
	for t.TokenFlags.Len() > 0 {
		element := t.TokenFlags.Front()
		tokenFlag := element.Value.(*TokenFlag)

		if tokenFlag.ExpireTime <= uint64(time.Now().Unix()) {
			delete(t.IndexID, tokenFlag.ID)
			t.TokenFlags.Remove(element)
		} else {
			return
		}
	}
}

func (t *TokenPool) push(token *token.Token) error {
	t.RemoveExpire()

	_, ok := t.IndexID[token.ID]
	if ok {
		return fmt.Errorf("token(ID: %d) exists in token pool", token.ID)
	}

	tokenFlag := &TokenFlag{
		ID:         token.ID,
		ExpireTime: token.ExpireTime,
	}

	var element *list.Element

	if t.TokenFlags.Len() == 0 {
		element = t.TokenFlags.PushBack(tokenFlag)
	} else {
		element = t.TokenFlags.Back()

		for {
			if element == nil || element.Value.(*TokenFlag).ExpireTime <= tokenFlag.ExpireTime {
				break
			}
			element = element.Prev()
		}

		if element == nil {
			element = t.TokenFlags.PushFront(tokenFlag)
		} else {
			element = t.TokenFlags.InsertAfter(tokenFlag, element)
		}
	}

	t.IndexID[tokenFlag.ID] = element

	return nil
}

func (t *TokenPool) PushWithSelfGenerate(token *token.Token) error {
	return t.push(token)
}

func (t *TokenPool) Remove(ID token.TokenID) error {
	t.RemoveExpire()

	element, ok := t.IndexID[ID]
	if !ok {
		return fmt.Errorf("token(ID: %d) is expired or not exists, cannot be remove", ID)
	}

	t.TokenFlags.Remove(element)
	delete(t.IndexID, ID)

	return nil
}

func (t *TokenPool) Check(token *token.Token) error {
	t.RemoveExpire()

	_, ok := t.IndexID[token.ID]
	if !ok {
		return fmt.Errorf("token(ID: %d) is expired or not exists, please try to acquire new token", token.ID)
	}

	if !token.CheckSign(t.DefaultKey) {
		return fmt.Errorf("token(ID: %d) is not valid", token.ID)
	}

	return nil
}

func (t *TokenPool) generateTokenNoCopyInfo(info map[string]string, expireSeconds uint32, key []byte) *token.Token {
	newToken := token.GenerateTokenNoCopyInfo(info, expireSeconds, key)

	for {
		_, ok := t.IndexID[newToken.ID]
		if !ok {
			break
		}
		newToken.ID = token.GenerateTokenID()
	}

	newToken.Sign(key)

	t.push(newToken)

	return newToken
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

func (t *TokenPool) GenerateTokenNoCopyInfo(info map[string]string) *token.Token {
	return t.generateTokenNoCopyInfo(info, t.DefaultExpireSeconds, t.DefaultKey)
}

func (t *TokenPool) GenerateToken(info map[string]string) *token.Token {
	return t.generateTokenNoCopyInfo(copyInfo(info), t.DefaultExpireSeconds, t.DefaultKey)
}
