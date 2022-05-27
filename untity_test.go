package main

import (
	"fmt"
	"testing"
)

func TestXxx(t *testing.T) {
	un := IntToBytes(1)
	fmt.Println(len(un))
	fmt.Println(BytesToInt(un))
}
