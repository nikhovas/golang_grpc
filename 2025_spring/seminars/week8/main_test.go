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
		db: gormDB,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email"}).AddRow(1, "John", "a@a.a")
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
			Email: "a@a.a",
		},
		resultUser,
	)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetValueFromRedis(t *testing.T) {
	redis, mockDB := redismock.NewClientMock()
	defer redis.Close()

	app := App{
		redis: redis,
	}

	mockDB.ExpectGet("abc").SetVal("ttt")

	ctx := context.Background()
	res, err := app.GetValueFromRedis(ctx, "abc")
	require.NoError(t, err)
	require.Equal(t, "ttt", res)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetValueFromRedisNoValue(t *testing.T) {
	redisDB, mockDB := redismock.NewClientMock()
	defer redisDB.Close()

	app := App{
		redis: redisDB,
	}

	mockDB.ExpectGet("abc").SetErr(redis.Nil)

	ctx := context.Background()
	res, err := app.GetValueFromRedis(ctx, "abc")
	require.NoError(t, err)
	require.Equal(t, "", res)

	require.NoError(t, mockDB.ExpectationsWereMet())
}
