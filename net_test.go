package util

import (
	"testing"
)

func TestGetLocalAddr(t *testing.T) {
	addrs, err := GetLocalAddr()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("addrs: %+v", addrs)
	}
}
