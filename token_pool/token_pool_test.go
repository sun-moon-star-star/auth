package token_pool

import (
	"testing"
)

func TestTokenPoolBase(t *testing.T) {
	pool := New(1200, []byte("DZFSJTDJQNYYJQRX"))
	token := pool.GenerateToken(nil)
	t.Log(token)
}
