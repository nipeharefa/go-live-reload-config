package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/nipeharefa/go-live-reload-config/model"
)

var ctx = context.Background()

func main() {
	var err error

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)

	config := model.Config{
		DBURL: "postgres://user=nipeharefa password=password dbname=pulang port=5432",
	}

	b, _ := json.Marshal(config)
	// Publish a message.
	err = rdb.Publish(ctx, "mychannel1", string(b)).Err()
	if err != nil {
		panic(err)
	}
}
