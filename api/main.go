package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Go-Short/api/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	defaultConfigPath = "./conf/conf.toml"
)

type Url struct {
	Url string      `json:"url"`
}

func (u *Url) Valid() bool {
	logger := service.GetLogger()

	// perform validation
	_, err := url.ParseRequestURI(u.Url)

	if err != nil {
		logger.Errorf(err.Error())
		return false
	}


	return true
}

func main() {
	err := readConfig(defaultConfigPath)
	logger := service.GetLogger()

	redisPort := viper.GetInt64("redis.port")
	redisHost := viper.GetString("redis.host")
	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	service.InitializeRedisClient(redisAddr)

	if err != nil {
		fmt.Printf("read config error: %s", err)
		return
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg" : "ok"})
	})

	r.GET("/:id", func(c *gin.Context) {
		logger := service.GetLogger()

		encodedID := c.Param("id")
		originalURL, err := service.GetOriginalURL(encodedID)

		if err != nil {
			logger.Info("could not locate original url")
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if strings.Index(originalURL, "http") == -1 {
			originalURL = "https://" + originalURL
		}

		logger.Infof("original URL: %s", originalURL)

		c.Redirect(http.StatusFound, originalURL)
		c.Abort()
	})

	r.POST("/", func(c *gin.Context) {
		// retrieve url
		rawRequestBody, _ := ioutil.ReadAll(c.Request.Body)

		var requestBody Url
		decoder := json.NewDecoder(bytes.NewReader(rawRequestBody))
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&requestBody)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		// validate
		if !requestBody.Valid() {
			c.JSON(400, gin.H{
				"error": "invalid url format",
			})
			return
		}

		// shorten
		shortUrl := service.ShortenURL(requestBody.Url)

		// return JSON
		c.JSON(200, gin.H{
			"url": shortUrl,
		})
		return
	})

	ginPort := fmt.Sprintf(":%d", viper.GetInt64("server.port"))
	logger.Info("Starting server")

	r.Run(ginPort)
}

func readConfig(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	viper.SetConfigType("toml")

	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("redis.port", "6379")
	return viper.ReadConfig(f)
}
