package math

import (
	"crypto/rand"
	"fmt"
	xRand "math/rand"
	"time"
)

// Source String used when generating a random identifier.
const idSource = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const idSourceUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const idSourceNumber = "0123456789"

// Save the length in a constant so we don't look it up each time.
const idSourceLen = byte(len(idSource))
const idSourceUpperLen = byte(len(idSourceUpper))
const idSourceNumberLen = byte(len(idSourceNumber))

// GenerateID creates a prefixed random identifier.
func RandString(prefix string, length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = idSource[b%idSourceLen]
	}

	// Return the formatted id
	return fmt.Sprintf("%s_%s", prefix, string(id))
}

func RandStringNew(length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = idSource[b%idSourceLen]
	}

	// Return the formatted id
	return string(id)
}
func RandStringUpper(prefix string, length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = idSourceUpper[b%idSourceUpperLen]
	}

	// Return the formatted id
	return fmt.Sprintf("%s_%s", prefix, string(id))
}

func RandStringNumber(prefix string, length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = idSourceNumber[b%idSourceNumberLen]
	}

	// Return the formatted id
	return fmt.Sprintf("%s_%s", prefix, string(id))
}

type RandStringMaker struct {
	Prefix string
	Length int
}

func (m *RandStringMaker) Next() string {
	return RandString(m.Prefix, m.Length)
}

var numbers = "0123456789"

func RandNumString(length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = numbers[b%10]
	}
	return string(id)
}

func RandNumStringUpper(length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = numbers[b%10]
	}
	return string(id)
}

func RandStringToUpper(length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = idSourceUpper[b%idSourceUpperLen]
	}
	return string(id)
}

// Source String used when generating a random identifier.
const idSourceNew = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
const idSourceUpperNew = "ABCDEFGHIJKLMNPQRSTUVWXYZ"
const idSourceNumberNew = "123456789"

// Save the length in a constant so we don't look it up each time.
const idSourceLenNew = byte(len(idSourceNew))
const idSourceUpperLenNew = byte(len(idSourceUpperNew))
const idSourceNumberLenNEw = byte(len(idSourceNumberNew))

// GenerateID creates a prefixed random identifier.
func RandStringQR(length int) string {
	// Create an array with the correct capacity
	id := make([]byte, length)
	// Fill our array with random numbers
	rand.Read(id)

	// Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = idSourceNew[b%idSourceLenNew]
	}

	// Return the formatted id
	return fmt.Sprintf("%s", string(id))
}

func GenCodes(quantity int64, lenght int) []string {
	var res = make([]string, quantity)
	for i := 0; int64(i) < quantity; i++ {
		var val = RandStringQR(lenght)
		res[i] = val
	}
	return res
}

func RandIntn(r *xRand.Rand) int {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	randomIndex := r.Intn(len(arr))
	pick := arr[randomIndex]
	return pick
}

func GetRandSource() *xRand.Rand {
	s := xRand.NewSource(time.Now().Unix())
	r := xRand.New(s) // initialize local pseudorandom generator
	return r
}

func GetIndexs(index string) []int {
	var mapIndex = make(map[string][]int)
	mapIndex["1"] = []int{2, 5, 8}
	mapIndex["2"] = []int{3, 9, 6}
	mapIndex["3"] = []int{3, 8, 5}
	mapIndex["4"] = []int{5, 2, 9}
	mapIndex["5"] = []int{2, 6, 8}
	mapIndex["6"] = []int{9, 2, 5}
	mapIndex["7"] = []int{7, 3, 9}
	mapIndex["8"] = []int{8, 5, 3}
	mapIndex["9"] = []int{6, 3, 8}
	return mapIndex[index]
}
