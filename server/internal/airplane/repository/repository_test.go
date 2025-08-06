package repository_test

import (
	"context"
	"fmt"
	"testing"
	"tiket_test/server/internal/airplane/repository"

	"github.com/redis/go-redis/v9"

	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	redisConn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r := repository.NewRepository(redisConn, context.TODO())
	res, err := r.GetAirplanes("1754425626816")
	fmt.Println(err)
	fmt.Print(res)
	assert.Nil(t, err)
}
