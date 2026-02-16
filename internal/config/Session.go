package config

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

var Store *session.Store

func InitSession() {
	redisAddres := os.Getenv("REDIS_ADDRESS")
	host := "localhost"
	port := 6379
	if h, p, err := net.SplitHostPort(redisAddres); err == nil {
		host = h
		if portNum, err := strconv.Atoi(p); err == nil {
			port = portNum
		}
	} else {
		if redisAddres != "" {
			host = redisAddres
		}
	}
	Storage := redis.New(redis.Config{
		Host:     host,
		Port:     port,
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: 0,
		Reset:    false,
	})

	Store = session.New(session.Config{
		Storage:        Storage,
		Expiration:     72 * time.Hour,
		CookieHTTPOnly: true,
		CookieSecure:   false,
		KeyLookup:      "cookie:session_id",
	})
	log.Println("✅ Redis Session Store initialized")
}
