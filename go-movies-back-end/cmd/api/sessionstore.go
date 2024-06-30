package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func (app *application) connectSessionStore() (*redis.Client, error) {
	// Connect to redis
	client := redis.NewClient(&redis.Options{
		Addr:     app.SessionStoreAdr,
		Password: app.SessionStorePass,
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to SessionStore")

	return client, nil
}
