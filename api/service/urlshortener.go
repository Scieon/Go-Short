package service

import (
	"errors"
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

	// need to get by URL
	existingID := getEncodedID(url)

 	// we already have this URL stored
	if existingID != "" {
		logger.Infof("%s was already stored", url)
		return "localhost:8080/" + existingID
	}

	// storing new url means incrementing id
	newEncodedID := Encode(redisClient.id)

	redisClient.id++

	redisClient.Client.HSet(newEncodedID, "url", url)
	redisClient.Client.HSet(url, "id", newEncodedID)
	return "localhost:8080/" + newEncodedID
}

func GetOriginalURL(encodedID string) (string, error) {
	redisClient := GetRedisClient()

	originalURL := redisClient.Client.HGet(encodedID, "url")

	if originalURL.Val() == "" {
		return "", errors.New("invalid id")
	}
	return originalURL.Val(), nil
}

func getEncodedID(storedURL string) string {
	redisClient := GetRedisClient()

	encodedID := redisClient.Client.HGet(storedURL, "id")

	return encodedID.Val()
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
