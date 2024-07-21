package cache

import (
	"errors"

	"github.com/go-redis/redis/v8"
)

var (
	ErrNoRecord  = redis.Nil
	ErrUnMarshal = errors.New("failed to unmarshal data")
)
