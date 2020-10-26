package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/nipeharefa/go-live-reload-config/model"
)

var config *model.Config
var ctx = context.Background()
var err error

const defaultDB string = "defaultDB"

func configHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(config)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func watchConfig(pubsub *redis.PubSub) {
	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()
	// Consume messages.
	for msg := range ch {
		rawIn := json.RawMessage(msg.Payload)
		bytes, err := rawIn.MarshalJSON()
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bytes, config)
		if err != nil {
			panic(err)
		}
	}
}

func main() {

	config = &model.Config{}
	r := mux.NewRouter()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()

	pubsub := rdb.Subscribe(ctx, "mychannel1")

	// define routes
	r.HandleFunc("/config", configHandler)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// watch on background
	go watchConfig(pubsub)

	log.Fatal(srv.ListenAndServe())
}
