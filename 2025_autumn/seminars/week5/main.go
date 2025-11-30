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
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	panicIfErr(err)
	fmt.Println("HERE")

	err = rdb.Set(ctx, "var", "val", 5*time.Second).Err()
	panicIfErr(err)

	val, err := rdb.Get(ctx, "var").Result()
	panicIfErr(err)
	fmt.Println(val)

	deleted, err := rdb.Del(ctx, "var", "key2").Result()
	panicIfErr(err)
	fmt.Println(deleted)

	val, err = rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("Not found")
	} else if err != nil {
		panicIfErr(err)
	} else {
		fmt.Println(val)
	}

	deleted, err = rdb.Del(ctx, "fruits", "fruits2", "unioned").Result()
	panicIfErr(err)
	fmt.Printf("DELETED: %d\n", deleted)

	added, err := rdb.SAdd(ctx, "fruits", "apple", "banana").Result()
	panicIfErr(err)
	fmt.Printf("ADDED1: %d\n", added)

	added, err = rdb.SAdd(ctx, "fruits", "apple", "cherry").Result()
	panicIfErr(err)
	fmt.Printf("ADDED2: %d\n", added)

	added, err = rdb.SAdd(ctx, "fruits2", "apple2", "cherry").Result()
	panicIfErr(err)
	fmt.Printf("ADDED2: %d\n", added)

	members, err := rdb.SMembers(ctx, "fruits").Result()
	panicIfErr(err)
	fmt.Printf("MEMBERS: %v\n", members)

	isMember, err := rdb.SIsMember(ctx, "fruits", "cherry").Result()
	panicIfErr(err)
	fmt.Printf("IS_MEMBER: %v\n", isMember)

	unionCount, err := rdb.SUnionStore(ctx, "unioned", "fruits", "fruits2").Result()
	panicIfErr(err)
	fmt.Printf("UNION COUNT: %d\n", unionCount)
}
