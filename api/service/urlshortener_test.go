package service

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testURL = "http://www.test.com"
var testEncodedID = "1C"

func TestShortenNewURL(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// setup redis addr
	InitializeRedisClient(s.Addr())

	shortenedURL := ShortenURL(testURL)
	fmt.Println(s.Addr())
	assert.Equal(t, "localhost:8080/1C", shortenedURL)
}

func TestShortenExistingURL(t *testing.T) {
	GetLogger()
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// setup redis addr
	InitializeRedisClient(s.Addr())

	shortenedURL := ShortenURL(testURL)
	shortenedURL = ShortenURL(testURL)

	assert.Equal(t, "localhost:8080/1C", shortenedURL)
}

func TestEncodeAndDecode(t *testing.T) {
	id := 100

	encodedID := Encode(id)
	decodedID := Decode(encodedID)

	assert.Equal(t, encodedID, testEncodedID)
	assert.Equal(t, decodedID, id)
}

func TestGetOriginalURL(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// setup redis addr
	InitializeRedisClient(s.Addr())

	s.HSet(testEncodedID, "url", testURL)
	url, err := GetOriginalURL(testEncodedID)

	assert.NoError(t, err)
	assert.Equal(t, testURL, url)
}

func TestGetEncodedID(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// setup redis addr
	InitializeRedisClient(s.Addr())

	s.HSet(testURL, "id", testEncodedID)
	id := getEncodedID(testURL)

	assert.Equal(t, id, testEncodedID)
}
