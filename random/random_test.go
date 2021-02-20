package random

import (
	"testing"
)

func TestRandomBase(t *testing.T) {
	uint64_id := RandomUint64()
	t.Log(uint64_id)

	int64_id := RandomInt64()
	t.Log(int64_id)

	uuid, err := RandomUUID()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(uuid)
}
