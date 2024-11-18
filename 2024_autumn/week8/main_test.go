package main

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetUserByID retrieves a user by ID using GORM
func TestGetUserByID(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	app := App{DB: gormDB}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "John", "a@b.c")
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	resultUser, err := app.GetUserByID(1)
	require.NoError(t, err)
	require.Equal(t, &User{ID: 1, Name: "John", Email: "a@b.c"}, resultUser)

	fmt.Println("HERE")

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetValueFromRedis(t *testing.T) {
	redis, mockdb := redismock.NewClientMock()

	app := App{Redis: redis}

	mockdb.ExpectGet("k").SetVal("v")

	ctx := context.Background()
	value, err := app.GetValueFromRedis(ctx, "k")
	require.NoError(t, err)
	require.Equal(t, "v", value)

	require.NoError(t, mockdb.ExpectationsWereMet())
}

func TestGetValueFromRedisNoValue(t *testing.T) {

	redisDB, mockdb := redismock.NewClientMock()

	app := App{Redis: redisDB}

	mockdb.ExpectGet("k").RedisNil()
	ctx := context.Background()
	_, err := app.GetValueFromRedis(ctx, "k")
	fmt.Println(err)
	require.ErrorIs(t, redis.Nil, err)

	require.NoError(t, mockdb.ExpectationsWereMet())
}
