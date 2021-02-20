package token_pool

import (
	"testing"
)

func TestTokenPoolBase(t *testing.T) {
	var expireSeconds uint32 = 1200
	key := []byte("DZFSJTDJQNYYJQRX")

	pool := New(expireSeconds, key)
	token := pool.GenerateToken(nil)
	t.Log(token)
}
