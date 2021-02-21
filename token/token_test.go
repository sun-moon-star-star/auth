package token

import (
	"encoding/json"
	"testing"
)

func TestTokenBase(t *testing.T) {
	key := "bajiuwenqingtian"

	token := &Token{
		ID:         GenerateTokenID(),
		CreateTime: 1612276579,
		ExpireTime: 1612276580,
		Info:       make(map[string]string),
	}

	token.Info["version"] = "1.0"
	token.Info["age"] = "21"
	token.Info["name"] = "zhao"

	token.Sign([]byte(key))

	t.Log(token.Signature)

	tokenJson, err := json.Marshal(token)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(tokenJson))
}

func BenchmarkTokenBase(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "bajiuwenqingtian"
		token := &Token{
			ID:         GenerateTokenID(),
			CreateTime: 1612276579,
			ExpireTime: 1612276580,
			Info:       make(map[string]string),
		}

		token.Info["version"] = "1.0"
		token.Info["age"] = "21"
		token.Info["name"] = "zhao"
		token.Info["sex"] = "female"
		token.Info["github"] = "https://github.com/sun-moon-star-star"

		token.Sign([]byte(key))

		tokenJson, err := json.Marshal(token)
		if err != nil {
			b.Fatal(err)
		}

		b.Log(string(tokenJson))
	}
}
