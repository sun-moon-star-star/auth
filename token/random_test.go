package token

import (
	"testing"
)

func TestRandomBase(t *testing.T) {
	uint64_id := GetUint64()
	t.Log(uint64_id)

	int64_id := GetInt64()
	t.Log(int64_id)

	uuid, err := GetUUID()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(uuid)
}
