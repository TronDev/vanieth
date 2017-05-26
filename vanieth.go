package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto" // you need to `go get` this, as is not in the stdlib
)

// Example run:
//
// $ go run vanieth.go 1234
//
// Address found:
// addr: 123411cc4a2e2e3238ee8e22d0d7b3cf2c8add9c
// pvt: 208439bf49edbc236bcffaa831e32006b91e6251150992fe5e704a3c3870415d
//
// https://github.com/ethereum/go-ethereum
//

// "main" method, generates a public key,  address
//
func addrGen(toMatch string, mode uint8) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	if mode == 0 {
		addrStr := hex.EncodeToString(addr[:])
		addrMatch(addrStr, toMatch, key)
	} else if mode == 1 {
		numOnly(addr[:], key)
	}
}

//Check if the address is only numbers
func numOnly(addr []byte, key *ecdsa.PrivateKey) {
	var byteA uint8
	var byteB uint8
	var success bool
	success = true
	for _, element := range addr {
		byteA = element>>4
		byteB = element<<4>>4
		if byteA > 9 {
			success = false
			break
		}
		if byteB > 9 {
			success = false
			break
		}

	}
	if success {
		keyStr := hex.EncodeToString(crypto.FromECDSA(key))
		addrFound(hex.EncodeToString(addr), keyStr)
		found = true
		os.Exit(0)
	}
}
// tries to match the address with the string provided by the user, exits if successful
//
func addrMatch(addrStr string, toMatch string, key *ecdsa.PrivateKey) {
	toMatch = strings.ToLower(toMatch)
	addrStrMatch := strings.TrimPrefix(addrStr, toMatch)
	found := addrStrMatch != addrStr
	if found {
		// fmt.Println("pub:", hex.EncodeToString(crypto.FromECDSAPub(&key.PublicKey))) // uncomment if you want the public key
		keyStr := hex.EncodeToString(crypto.FromECDSA(key))
		addrFound(addrStr, keyStr)
		found = true
		os.Exit(0) // here the program exits when it found a match
	}
}

//Keeps track of addresses per second
func watchman() {
	for !found {
		time.Sleep(15000 * time.Millisecond)
		fmt.Printf("Generating addresses at %d A/s\n", addressPerSecond/15)
		addressPerSecond = 0
	}
}
func worker(toMatch string, mode uint8) {
	for !found {
		addrGen(toMatch, mode)
		addressPerSecond += 1
	}
}
// main, executes addrGen ad-infinitum, until a match is found
//
var addressPerSecond int
var found bool
func main() {
	runtime.GOMAXPROCS(8)
	found = false
	var toMatch string
	if len(os.Args) == 1 {
		errNoArg()
		os.Exit(1)
	} else {
		toMatch = os.Args[1]
		// errWrongMatch(toMatch)
	}
	go watchman()
	//TODO create workers up to GOMAXPROCS
	if toMatch != "num" {
		fmt.Printf("Looking for 0x%s...", toMatch)
		worker(toMatch, 0)
	} else {
		fmt.Printf("Looking for numbers-only address")
		worker(toMatch, 1)
	}
}

// non-interesting functions follow...

func addrFound(addrStr string, keyStr string) {
	println("Address found:")
	fmt.Printf("addr: 0x%s\n", addrStr)
	fmt.Printf("pvt: 0x%s\n", keyStr)
	println("\nexiting...")
}

func errNoArg() {
	println("You need to pass a vanity match, retry with an extra agrument like: 42")
	println("\nexample: go run vanieth.go 42")
	println("\nexiting...")
}

// func errWrongMatch(match string) {
// 	strings.ContainsAny(match, "")
// 	if (wrongMatch) {
// 		println("You need to pass a findable address")
// 	}
// }
