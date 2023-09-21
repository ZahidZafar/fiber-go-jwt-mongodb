package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
)

const (
	salt = "0vslM5fwJBSytKuZgqaRxK6PmPzLcYsa"
)

func EncryptWithSHA512(v string) string {
	// Data to be hashed
	data := []byte(v + salt)

	// Create a new SHA-256 hash object
	hasher := sha512.New()

	// Write the data to the hasher
	hasher.Write(data)

	// Get the final hash sum as a byte slice
	hashBytes := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hashBytes)
}

func RandomNumber() int {
	// Seed the random number generator with the current time
	//rand.Seed(time.Now().UnixNano())

	// Generate a random 6-digit number
	min := 100000 // Minimum 6-digit number
	max := 999999 // Maximum 6-digit number
	randomNum := rand.Intn(max-min+1) + min

	return randomNum
}
