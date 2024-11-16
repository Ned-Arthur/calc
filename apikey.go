package main

import (
	"crypto/rand"
	"encoding/base64"
)

const keyLength = 32

func generateKey() string {
	bytes := make([]byte, keyLength)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(bytes)[:keyLength]
}
