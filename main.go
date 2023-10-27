package main

import (
	"keychain/keychain"
	"log"
)

var (
	keychainFilePath = "/Users/ycb/Desktop/chainbreaker/login.keychain-db"
)

func main() {

	_, err := keychain.NewKeychain(keychainFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}

	println("*******")
}
