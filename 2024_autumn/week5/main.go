package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Pong: %s\n", pong)

	panicIfErr(rdb.Set(ctx, "key", "value", 5*time.Second).Err())

	val, err := rdb.Get(ctx, "key").Result()
	panicIfErr(err)
	fmt.Println(val)

	deleted, err := rdb.Del(ctx, "key", "key2").Result()
	panicIfErr(err)
	fmt.Printf("Deleted number: %d\n", deleted)

	val, err = rdb.Get(ctx, "key").Result()
	if err == redis.Nil {
		fmt.Println("No value")
	} else if err != nil {
		panicIfErr(err)
	} else {
		// Не должны сюда попасть
		fmt.Println(val)
	}

	// clear
	panicIfErr(rdb.Del(ctx, "fruits", "fruits2").Err())

	fruits := []string{"apple", "banana", "cherry", "date"}
	added, err := rdb.SAdd(ctx, "fruits", fruits).Result()
	panicIfErr(err)
	fmt.Printf("Added number: %d\n", added)

	members, err := rdb.SMembers(ctx, "fruits").Result()
	panicIfErr(err)
	fmt.Printf("Members: %v\n", members)

	// isMember, err := rdb.SIsMember(ctx, "fruits", "banana").Result() // true
	isMember, err := rdb.SIsMember(ctx, "fruits", "banana2").Result() // false
	panicIfErr(err)
	fmt.Printf("Is member: %v", isMember)

	fruits2 := []string{"apple", "banana2", "cherry2"}
	added, err = rdb.SAdd(ctx, "fruits2", fruits2).Result()
	panicIfErr(err)
	fmt.Printf("Added number: %d\n", added)

	unionCount, err := rdb.SUnionStore(ctx, "union_set", "fruits", "fruits2").Result()
	panicIfErr(err)
	fmt.Printf("Unioned number: %d\n", unionCount)

	unionMembers, err := rdb.SMembers(ctx, "union_set").Result()
	panicIfErr(err)
	fmt.Printf(" Union members: %v\n", unionMembers)
}
