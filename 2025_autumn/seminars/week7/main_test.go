package main

import (
	"context"
	"fmt"
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

	dialector := postgres.New(postgres.Config{Conn: db})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "aaa", "aaa@b.ru")
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	app := App{
		db: gormDB,
	}

	resultUser, err := app.GetUserByID(1)
	require.NoError(t, err)
	require.Equal(
		t,
		&User{
			ID:    1,
			Name:  "aaa",
			Email: "aaa@b.ru",
		},
		resultUser,
	)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetValueFromRedis(t *testing.T) {
	redisDB, mockDB := redismock.NewClientMock()
	defer redisDB.Close()

	mockDB.ExpectGet("asdf").SetVal("kkk")

	app := App{
		redis: redisDB,
	}

	result, err := app.GetValueFromRedis(context.Background(), "asdf")
	require.NoError(t, err)
	require.Equal(t, "kkk", result)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetValueFromRedisNoValue(t *testing.T) {
	redisDB, mockDB := redismock.NewClientMock()
	defer redisDB.Close()

	mockDB.ExpectGet("asdf").SetErr(redis.Nil)

	app := App{
		redis: redisDB,
	}

	fmt.Println("QQQQQQQQQQQ")
	println("IIIIII")
	t.Error("AHAHAHAHA")

	result, err := app.GetValueFromRedis(context.Background(), "asdf")
	require.NoError(t, err)
	require.Equal(t, "f", result)

	require.NoError(t, mockDB.ExpectationsWereMet())
}
