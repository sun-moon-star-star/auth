package token_pool

import (
	"testing"
)

func TestTokenPoolBase(t *testing.T) {
	pool := New()
	token := pool.GenerateTokenID()
	t.Log(token)
}
