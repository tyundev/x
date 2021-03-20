package math

import (
	"crypto/rand"
	"fmt"
	mRand "math/rand"
)

// Source String used when generating a random identifier.
const idSourceNew = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
const idSourceUpperNew = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
const idSourceNumberNew = "123456789"

// Save the length in a constant so we don't look it up each time.
const idSourceLenNew = byte(len(idSource))
const idSourceUpperLenNew = byte(len(idSourceUpper))
const idSourceNumberLenNEw = byte(len(idSourceNumber))

// GenerateID creates a prefixed random identifier.
func RandStringQR(length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = idSource[b%idSourceLen]
	}

	// Return the formatted id
	return fmt.Sprintf("%s", string(id))
}

func GenCodes(cycle, lenght int) []string {
	var res = make([]string, cycle)
	for i := 0; i < cycle; i++ {
		var val = RandStringQR(lenght)
		res[i] = val
	}
	return res
}

func RandIntn() int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	randomIndex := mRand.Intn(9)
	pick := arr[randomIndex]
	return pick
}
