package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/acoshift/configfile"
	"github.com/acoshift/hime"
	"github.com/anthoz69/beam/app"
	"github.com/garyburd/redigo/redis"

	_ "github.com/lib/pq"
)

func main() {
	config := configfile.NewReader("config")

	db, err := sql.Open("postgres", config.String("db"))
	if err != nil {
		log.Fatal(err)
	}

	sessionHost := config.String("session_host")

	redisPool := &redis.Pool{
		IdleTimeout: 60 * time.Minute,
		MaxIdle:     2,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", sessionHost)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	appFactory := app.New(&app.Config{
		DB:            db,
		SessionKey:    config.Bytes("session_key"),
		SessionSecret: config.Bytes("session_secret"),
		RedisPrefix:   config.String("session_prefix"),
		RedisPool:     redisPool,
	})

	err = hime.New().
		TemplateDir("template").
		TemplateRoot("root").
		Minify().
		Handler(appFactory).
		GracefulShutdown().
		ListenAndServe(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
