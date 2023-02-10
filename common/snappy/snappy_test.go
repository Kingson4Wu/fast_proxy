package snappy

import (
	"fmt"
	"testing"
)

func TestCompressPanic(t *testing.T) {

	_, ok := compressPanic([]byte(""), false)

	fmt.Println(ok)

	if !ok {
		t.Fatal("success failure")
	}

	_, ok = compressPanic([]byte(""), true)

	if ok {
		t.Fatal("panic failure")
	}

}
