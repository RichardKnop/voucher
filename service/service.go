package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RichardKnop/voucher/pb"
	"github.com/go-redis/redis"
)

const countPerPage = int64(10)

type impl struct {
	redisClient *redis.Client
}

// IFace ...
type IFace interface {
	Create(data []byte) (*pb.Voucher, int64, error)
	FindByID(voucherID string) (*pb.Voucher, int64, error)
	FindAll(cursor uint64) ([]*pb.Voucher, uint64, int64, error)
}

// New ...
func New(redisClient *redis.Client) IFace {
	return &impl{redisClient: redisClient}
}

// Create ...
func (svc *impl) Create(data []byte) (*pb.Voucher, int64, error) {
	voucher := new(pb.Voucher)
	if err := json.Unmarshal(data, voucher); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to unmarshal data into voucher entity: %s", err)
	}

	// validate voucher ID
	if voucher.Id == "" {
		return nil, http.StatusBadRequest, errors.New("voucher.Id cannot be empty")
	}

	// validade createdAt and expiresAt
	now := time.Now()

	createdAt := parseTime(voucher.CreatedAt, now)
	voucher.CreatedAt = createdAt.Format(time.RFC3339)

	expiresAt := parseTime(voucher.CreatedAt, time.Time{})
	expires := time.Duration(0)
	if !expiresAt.IsZero() {
		if expiresAt.Before(now) {
			return nil, http.StatusBadRequest, errors.New("voucher.ExpiresAt must be in the future")
		}
		expires = expiresAt.Sub(now)
	}
	voucher.ExpiresAt = expiresAt.Format(time.RFC3339)

	// check if voucher already exists
	exists, err := svc.redisClient.Exists(voucher.Id).Result()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("redis error: %s", err)
	}
	if exists == 1 {
		return nil, http.StatusInternalServerError, errors.New("voucher.Id already used")
	}

	// 0 means no expiration
	if err := svc.redisClient.Set(voucher.Id, data, expires).Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return voucher, 0, nil
}

// FindByID ...
func (svc *impl) FindByID(voucherID string) (*pb.Voucher, int64, error) {
	data, err := svc.redisClient.Get(voucherID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, http.StatusNotFound, errors.New("voucher.Id not found")
		} else {
			return nil, http.StatusInternalServerError, fmt.Errorf("redis error: %s", err)
		}
	}

	voucher := new(pb.Voucher)
	if err := json.Unmarshal(data, voucher); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to unmarshal data into voucher entity: %s", err)
	}

	return voucher, 0, nil
}

// FindAll ...
func (svc *impl) FindAll(cursor uint64) ([]*pb.Voucher, uint64, int64, error) {
	var (
		n    int
		keys []string
		err  error
	)
	for {
		keys, cursor, err = svc.redisClient.Scan(cursor, "*", countPerPage).Result()
		if err != nil {
			return nil, 0, http.StatusInternalServerError, err
		}
		n += len(keys)
		if cursor == 0 {
			break
		}
	}

	datas, err := svc.redisClient.MGet(keys...).Result()
	if err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}

	vouchers := make([]*pb.Voucher, 0)
	for _, dataInterface := range datas {
		data, ok := dataInterface.(string)
		if !ok {
			continue
		}

		voucher := new(pb.Voucher)
		if err := json.Unmarshal([]byte(data), voucher); err != nil {
			return nil, 0, http.StatusInternalServerError, fmt.Errorf("failed to unmarshal data into voucher entity: %s", err)
		}
		vouchers = append(vouchers, voucher)
	}

	return vouchers, cursor, 0, nil
}

func parseTime(val string, defaultVal time.Time) time.Time {
	if val == "" {
		return defaultVal
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return defaultVal
	}
	return t
}
