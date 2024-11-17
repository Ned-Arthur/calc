package main

import (
	"crypto/rand"
	"encoding/base64"
)

const KEYLENGTH = 32

func generateKey() string {
	bytes := make([]byte, KEYLENGTH)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(bytes)[:KEYLENGTH]
}
