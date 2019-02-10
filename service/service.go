package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/RichardKnop/voucher/pb"
	"github.com/go-redis/redis"
)

const (
	defaultPageSize = int64(10)
	prefix          = "voucher_"
	index           = "__index__"
)

var (
	isAlpha = regexp.MustCompile(`^[A-Za-z123]+$`).MatchString
)

type impl struct {
	redisClient *redis.Client
}

// IFace ...
type IFace interface {
	Create(data []byte) (*pb.Voucher, int64, error)
	FindByID(voucherID string) (*pb.Voucher, int64, error)
	FindAll(offset, count int64) ([]*pb.Voucher, int64, int64, error)
	DeleteByID(voucherID string) (int64, error)
}

// New ...
func New(redisClient *redis.Client) IFace {
	return &impl{redisClient: redisClient}
}

// ValidateVoucherID ...
func ValidateVoucherID(voucherID string) error {
	if !isAlpha(voucherID) {
		return errors.New("voucher ID must be alphanumeric")
	}
	return nil
}

// Create ...
func (svc *impl) Create(data []byte) (*pb.Voucher, int64, error) {
	voucher := new(pb.Voucher)
	if err := json.Unmarshal(data, voucher); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("failed to unmarshal data into voucher entity: %s", err)
	}

	// validate voucher ID
	if voucher.Id == "" {
		return nil, http.StatusBadRequest, errors.New("voucher ID cannot be empty")
	}
	voucher.Id = fmt.Sprintf("%s%s", prefix, voucher.Id)

	// validade createdAt and expiresAt
	now := time.Now()

	createdAt := parseTime(voucher.CreatedAt, now)
	voucher.CreatedAt = createdAt.Format(time.RFC3339)

	expiresAt := parseTime(voucher.ExpiresAt, time.Time{})
	expires := time.Duration(0)
	if !expiresAt.IsZero() {
		if expiresAt.Before(now) {
			return nil, http.StatusBadRequest, errors.New("voucher expiresAt must be in the future")
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
		return nil, http.StatusInternalServerError, errors.New("voucher ID already used")
	}

	// Store the voucher and add the key to the index
	_, err = svc.redisClient.Pipelined(func(pipe redis.Pipeliner) error {
		if setErr := svc.redisClient.Set(voucher.Id, data, expires).Err(); setErr != nil {
			return setErr
		}
		return svc.redisClient.ZAdd(index, redis.Z{Score: 0, Member: voucher.Id}).Err()
	})
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("redis error: %s", err)
	}

	log.Printf("Created a new voucher \"%s\"", voucher.Id)

	return voucher, 0, nil
}

// FindByID ...
func (svc *impl) FindByID(voucherID string) (*pb.Voucher, int64, error) {
	voucherID = fmt.Sprintf("%s%s", prefix, voucherID)

	data, err := svc.redisClient.Get(voucherID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, http.StatusNotFound, errors.New("voucher ID not found")
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
func (svc *impl) FindAll(offset, count int64) ([]*pb.Voucher, int64, int64, error) {
	if count <= 0 {
		count = defaultPageSize
	}

	log.Printf("Listing vouchers offset=%d, count=%d", offset, count)

	total := svc.redisClient.ZCount(index, "-inf", "+inf").Val()
	if offset >= total {
		return []*pb.Voucher{}, 0, 0, nil
	}

	keys, err := svc.redisClient.ZRange(index, offset, offset+count-1).Result()
	if err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}

	if len(keys) < 1 {
		return []*pb.Voucher{}, 0, 0, nil
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

	nextOffset := offset + count
	if nextOffset >= total {
		nextOffset = -1
	}

	return vouchers, nextOffset, 0, nil
}

// DeleteByID ...
func (svc *impl) DeleteByID(voucherID string) (int64, error) {
	log.Printf("Deleting voucher \"%s\"", voucherID)

	voucherID = fmt.Sprintf("%s%s", prefix, voucherID)

	// Delete the voucher and the key from the index
	_, err := svc.redisClient.Pipelined(func(pipe redis.Pipeliner) error {
		if delErr := svc.redisClient.Del(voucherID).Err(); delErr != nil {
			return delErr
		}
		return svc.redisClient.ZRem(index, voucherID).Err()
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("redis error: %s", err)
	}
	return 0, nil
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
