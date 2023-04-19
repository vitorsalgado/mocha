package mnet

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestMNet(t *testing.T) {
	ctx := context.Background()

	mnet := New(":26379")
	mnet.Listen()

	rdb := redis.NewClient(&redis.Options{
		Addr:     mnet.listener.Addr().String(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value3", 0).Err()
	if err != nil {
		fmt.Println(err.Error())
		// panic(err)
	}
}
