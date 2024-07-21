package utils

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	alphabets = "abcdefghijklmnopqrstuvwxyz"
	numbers   = "0123456789"
	symbols   = "_%#:>,<@!`$*()"
)

var (
	randSource                        rand.Source
	alphabetLen, numberLen, symbolLen int
)

func init() {
	randSource = rand.NewSource(time.Now().UnixNano())
	alphabetLen = len(alphabets)
	numberLen = len(numbers)
	symbolLen = len(symbols)
}

// RandomUUID generates a random UUID
func RandomUUID() uuid.UUID {
	return uuid.New()
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	r := rand.New(randSource)
	return min + r.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	r := rand.New(randSource)
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = alphabets[r.Intn(alphabetLen)]
	}
	return string(result)
}

// RandomNumberAsString generates a random number of length digits
func RandomNumberAsString(digits int) string {
	r := rand.New(randSource)
	result := make([]byte, digits)
	for i := 0; i < digits; i++ {
		result[i] = numbers[r.Intn(numberLen)]
	}
	return string(result)
}

// RandomSymbols generates a random string of length n
func RandomSymbols(n int) string {
	r := rand.New(randSource)
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = symbols[r.Intn(symbolLen)]
	}
	return string(result)
}

// RandomEmail generates a random email
func RandomEmail() string {
	randomString := RandomString(5)
	randomInt := strconv.Itoa(int(RandomInt(1, 20)))
	return randomString + randomInt + "@gmail.com"
}

// RandomUsername generates a random username
func RandomUsername() string {
	randomStringLength := int(RandomInt(5, 30))
	randomInt := strconv.Itoa(int(RandomInt(1, 20)))
	return RandomString(randomStringLength) + randomInt
}

// RandomPassword generates a random password
func RandomPassword() string {
	randomStringLength := int(RandomInt(8, 30))
	symbols := RandomSymbols(5)
	randomInt := strconv.Itoa(int(RandomInt(1, 20)))
	return RandomString(randomStringLength) + symbols + randomInt
}

func RandomSalt(size ...int) string {
	saltSize := 16
	if len(size) > 0 {
		value := size[0]
		if value > 0 {
			saltSize = value
		}
	}
	var buf bytes.Buffer
	for i := 0; i < saltSize; i++ {
		switch rand.Intn(3) {
		case 0:
			buf.WriteByte(alphabets[rand.Intn(alphabetLen)])
		case 1:
			buf.WriteByte(numbers[rand.Intn(numberLen)])
		case 2:
			buf.WriteByte(symbols[rand.Intn(symbolLen)])
		}
	}
	runes := []rune(buf.String())
	r := rand.New(randSource)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}
