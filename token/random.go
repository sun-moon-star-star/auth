package token

import (
	"fmt"
	"math/rand"
	"time"
)

func RandSeedTime() {
	rand.Seed(time.Now().Unix())
}

func GetUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
    	b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}