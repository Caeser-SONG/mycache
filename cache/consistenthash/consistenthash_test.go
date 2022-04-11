package consistenthash

import (
	"strconv"
	"testing"
)

func Testconsistenthash(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":  "2",
		"11": "33",
		"32": "3211",
		"65": "66",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("asking for %s ,should haveyielded %s", k, v)
		}
	}
	hash.Add("8")
	testCases["65"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("asking for %s ,should haveyielded %s", k, v)
		}
	}
}
