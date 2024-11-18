package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents a user in the database
type User struct {
	ID    uint
	Name  string
	Email string
}

// App holds the database and Redis clients
type App struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewApp initializes the App with PostgreSQL and Redis connections
func NewApp(dsn string, redisOptions *redis.Options) (*App, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the User struct
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	rdb := redis.NewClient(redisOptions)

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &App{
		DB:    db,
		Redis: rdb,
	}, nil
}

// GetUserByID retrieves a user by ID using GORM
func (app *App) GetUserByID(id uint) (*User, error) {
	var user User
	result := app.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetValueFromRedis retrieves a value from Redis by key
func (app *App) GetValueFromRedis(ctx context.Context, key string) (string, error) {
	val, err := app.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func main() {
	// Example DSN for PostgreSQL
	dsn := "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable TimeZone=UTC"

	// Redis options
	redisOptions := &redis.Options{
		Addr: "localhost:6379",
	}

	app, err := NewApp(dsn, redisOptions)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Example usage
	user, err := app.GetUserByID(1)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
	} else {
		fmt.Printf("User: %+v\n", user)
	}

	ctx := context.Background()
	value, err := app.GetValueFromRedis(ctx, "mykey")
	if err != nil {
		log.Printf("Error fetching from Redis: %v", err)
	} else {
		fmt.Printf("Redis Value: %s\n", value)
	}
}
