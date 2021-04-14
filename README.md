# auth

#### token
```go
package main

import (
	"fmt"
	"time"

	"github.com/sun-moon-star-star/auth/token"
)

func main() {
	key := "bajiuwenqingtian"

	token := &token.Token{
		ID:         token.GenerateTokenID(),
		CreateTime: 1612276579,
		ExpireTime: 1612276580,
		Info:       make(token.TokenInfo),
	}

	token.Info["version"] = "1.0"
	token.Info["age"] = "22"
	token.Info["name"] = "sun-moon-star-star"

	token.Sign([]byte(key))

	if !token.CheckSign([]byte(key)) {
		fmt.Errorf("check strategy failed")
	}

	time.Sleep(time.Second)
	// timeout failed token become invalid
	if token.Check([]byte(key)) {
		fmt.Errorf("check strategy failed")
	}
}
```