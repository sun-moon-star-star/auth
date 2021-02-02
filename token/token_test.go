package token

import (
    "testing"
    "encoding/json"
)

func TestTokenBase(t *testing.T) {
    key := "bajiuwenqingtian"

    token := &Token{
        CreateTime : 1612276579,
        ExpireSeconds : 3600,
        Info: make(map[string]string),
    }

    token.Info["age"] = "21"
    token.Info["name"] = "zhao"
    
    token.Signature = Signature(token, []byte(key))

    t.Log(token.Signature)

    tokenJson, err := json.Marshal(token)
    if err != nil {
        t.Fatal(err)
    }
    
    t.Log(string(tokenJson))
}

func BenchmarkTokenBase(b *testing.B){
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        key := "bajiuwenqingtian"
        token := &Token{
            CreateTime : 1612276579,
            ExpireSeconds : 3600,
            Info: make(map[string]string),
        }

        token.Info["age"] = "21"
        token.Info["name"] = "zhao"
        token.Info["sex"] = "female"
        token.Info["github"] = "https://github.com/sun-moon-star-star"
        
        token.Signature = Signature(token, []byte(key))

        tokenJson, err := json.Marshal(token)
        if err != nil {
            b.Fatal(err)
        }
        
        b.Log(string(tokenJson))
    }
}