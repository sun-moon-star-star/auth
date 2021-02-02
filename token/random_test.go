package token

import (
	"testing"
)

func TestRandomBase(t *testing.T) {
	RandSeedTime()
	uuid, err := GetUUID()
	if err != nil {
        t.Fatal(err.Error())
	}
	
	t.Log(uuid)
}