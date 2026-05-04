package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	err = rdb.Set(ctx, "var", "val", 5*time.Second).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "var").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)

	deleted, err := rdb.Del(ctx, "var", "var2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted number: %d\n", deleted)

	val, err = rdb.Get(ctx, "var").Result()
	if err == redis.Nil {
		fmt.Println("Not found")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println(val)
	}

	_, err = rdb.Del(ctx, "fruits").Result()
	if err != nil {
		panic(err)
	}

	added, err := rdb.SAdd(ctx, "fruits", "apple", "banana").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Added1 %d\n", added)

	added, err = rdb.SAdd(ctx, "fruits", "apple", "cherry").Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Added2 %d\n", added)

	members, err := rdb.SMembers(ctx, "fruits").Result()
	fmt.Printf("Members: %v\n", members)

	isMember, err := rdb.SIsMember(ctx, "fruits", "apple").Result()
	fmt.Printf("IsMember: %v\n", isMember)

	count, err := rdb.SCard(ctx, "fruits").Result()
	fmt.Printf("Count: %d\n", count)

	added, err = rdb.SAdd(ctx, "fruits2", "apple", "pineapple").Result()

	unionCount, err := rdb.SUnionStore(ctx, "unioned", "fruits", "fruits2").Result()
	fmt.Printf("Unioned %d\n", unionCount)
}
