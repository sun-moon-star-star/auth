package token_pool

import (
	"auth/token"
	"strings"
	"testing"
	"time"
)

func TestTokenPoolBase(t *testing.T) {
	pool := New()
	pool.DefaultExpireSeconds = 2

	token := pool.GenerateTokenID()
	if err := pool.Check(token); err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)
	if err := pool.Check(token); err == nil {
		t.Fatal("check strategy failed")
	}
}

func TestTokenPoolExpired(t *testing.T) {
	pool := New()
	pool.DefaultExpireSeconds = 2

	token := pool.GenerateTokenID()
	if err := pool.Check(token); err != nil {
		t.Fatal(err)
	}

	pool.Expired(TokenFlag{ID: token.ID, ExpireTime: token.ExpireTime})
	if err := pool.Check(token); err == nil {
		t.Fatal("check strategy failed")
	}
}

func TestTokenPoolCheckStrategy(t *testing.T) {
	pool := New()
	pool.DefaultExpireSeconds = 2
	pool.PushStrategy = func(token *token.Token, pool *TokenPool) bool {
		return false
	}
	pool.CheckStrategy = func(token *token.Token, pool *TokenPool) bool {
		value, ok := token.Info["version"]
		if !ok {
			return false
		}

		if strings.Compare(value, "1.0") == 0 {
			return false
		}

		return token.Check(pool.DefaultKey)
	}

	info := map[string]string{
		"version": "1.0",
		"name":    "sun-moon-star-star",
	}

	token := pool.GenerateTokenNoCopyInfo(info)
	if err := pool.Check(token); err == nil {
		t.Fatal(err)
	}

	info["version"] = "2.0"
	t.Log(token)
	if err := pool.Check(token); err == nil {
		t.Fatal(err)
	}
}
