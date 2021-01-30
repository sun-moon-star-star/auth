package token

import (
    "testing"
    "time" 
)

func TestBase(t *testing.T) {
    key := "bajiuwenqingtian"

    if len(key) % 16 != 0 {
        t.Fatal("key len error")
    }

    rawTokenContent := &TokenContent{
        CreateTime : uint64(time.Now().Unix()),
        ID : "sun-moon-star-star",
        Message : "天若不爱酒，酒星不在天。地若不爱酒，地应无酒泉。",
    }

    token, err := GenerateToken(rawTokenContent, []byte(key))

    if err != nil {
        t.Fatal(err.Error())
	}

    t.Log(rawTokenContent.CreateTime, rawTokenContent.ID, rawTokenContent.Message)

    var tokenContent *TokenContent
    tokenContent, err = CheckToken(token, []byte(key))
    if err != nil {
        t.Fatal(err.Error())
    }

    if tokenContent != tokenContent{
        t.Fatal("not useful")
    }

    t.Log(tokenContent.CreateTime, tokenContent.ID, tokenContent.Message)
}