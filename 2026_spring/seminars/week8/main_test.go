package main

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetUserByID(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	app := App{
		db:    gormDB,
		redis: nil,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "John", "a@a.ru")
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	resultUser, err := app.GetUserByID(1)
	require.NoError(t, err)
	require.Equal(
		t,
		&User{
			ID:    1,
			Name:  "John",
			Email: "a@a.ru",
		},
		resultUser,
	)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetValueFromRedis(t *testing.T) {
	redis, mockRedis := redismock.NewClientMock()
	defer redis.Close()

	mockRedis.ExpectGet("a").SetVal("b")

	app := App{
		db:    nil,
		redis: redis,
	}

	result, err := app.GetValueFromRedis(context.Background(), "a")
	require.NoError(t, err)
	require.Equal(t, "b", result)

	require.NoError(t, mockRedis.ExpectationsWereMet())
}

func TestGetValueFromRedisNoValue(t *testing.T) {
	redisDB, mockRedis := redismock.NewClientMock()
	defer redisDB.Close()

	mockRedis.ExpectGet("a").SetErr(redis.Nil)

	app := App{
		db:    nil,
		redis: redisDB,
	}

	result, err := app.GetValueFromRedis(context.Background(), "a")
	require.NoError(t, err)
	require.Equal(t, "", result)

	require.NoError(t, mockRedis.ExpectationsWereMet())
}
