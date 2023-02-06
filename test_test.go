package main

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestSendBills(t *testing.T) {
	hexStr := hex.EncodeToString([]byte{0, 0, 0, 2, 1, 3})
	log.Println(hexStr)
}
