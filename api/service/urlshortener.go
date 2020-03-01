package service

import (
	"math"
	"strings"
)

var (
	// CharacterSet consists of 62 characters [0-9][A-Z][a-z].
	CharacterSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Base         = 62
)

func ShortenURL(url string) string {
	redisClient := GetRedisClient()

	// incr id
	encodedID := Encode(redisClient.id)
	newID := redisClient.client.HIncrBy("urls", encodedID, 1)
	redisClient.client.HSet("urls", newID.Val(), url)

	shortUrl := "localhost:8080/" + encodedID
	return shortUrl
}

func Encode(num int) string {
	b := make([]byte, 0)

	for num > 0 {
		r := math.Mod(float64(num), float64(Base))

		num /= Base

		b = append([]byte{CharacterSet[int(r)]}, b...)
	}

	return string(b)
}

func Decode(s string) int {
	var id, pow int

	for index, val := range s {
		// convert position to power
		pow = len(s) - (index + 1)

		// IndexRune returns the index of the first instance of the Unicode code point
		pos := strings.IndexRune(CharacterSet, val)

		// calculate
		id += pos * int(math.Pow(float64(Base), float64(pow)))
	}

	return int(id)
}
