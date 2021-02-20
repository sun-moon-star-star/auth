package token_pool

import (
	"testing"
	"time"
)

func TestTokenPoolBase(t *testing.T) {
	pool := New()
	pool.DefaultExpireSeconds = 2

	token := pool.GenerateTokenID()
	t.Log(token)
	if err := pool.Check(token); err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)
	t.Log(token)
	if err := pool.Check(token); err == nil {
		t.Fatal("check strategy failed")
	}
}
