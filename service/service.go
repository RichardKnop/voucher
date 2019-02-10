package service

import (
	"github.com/go-redis/redis"
)

type impl struct {
	redisClient *redis.Client
}

// IFace ...
type IFace interface {
	Create(key string, score int, data []byte) error
}

// New ...
func New(redisClient *redis.Client) IFace {
	return &impl{redisClient: redisClient}
}

// Create ...
func (svc *impl) Create(key string, score int, data []byte) error {
	return nil
}
