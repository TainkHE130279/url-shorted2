package utils

import (
	"context"
	cryptorand "crypto/rand"
	"errors"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrTimeout = errors.New("try lock timeout")

type Config struct {
	MaxLockTime time.Duration
	MaxTryTime  time.Duration
}

type IDLock interface {
	Lock(ctx context.Context, key string) (*LockData, error)
	Unlock(ctx context.Context, ld *LockData) error
}
type LockData struct {
	Key   string
	Value string
}

type redisLock struct {
	prefix string
	config *Config
	client *redis.Client
}

func NewRedisLock(rclient *redis.Client, cfg *Config) *redisLock {
	var rlock = &redisLock{
		prefix: "lock",
		config: cfg,
		client: rclient,
	}

	return rlock
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		// Use crypto/rand for secure random number generation
		num, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			// Fallback to uuid if crypto/rand fails
			id, _ := uuid.NewRandom()
			return strings.ReplaceAll(id.String(), "-", "")
		}
		b[i] = letterRunes[num.Int64()]
	}
	return string(b)
}

func mustId() string {
	var id, err = uuid.NewRandom()
	if nil != err {
		return randString(32)
	}
	return strings.ReplaceAll(id.String(), "-", "")
}

func (rlock *redisLock) Lock(ctx context.Context, key string) (*LockData, error) {
	val := mustId()
	var ld = &LockData{Key: key, Value: val}
	ctxEx, cancel := context.WithTimeout(context.TODO(), rlock.config.MaxTryTime)
	defer cancel()

	for {
		select {
		case <-ctxEx.Done():
			return nil, ErrTimeout

		case <-time.After(20 * time.Millisecond):
			success, err := rlock.tryLock(ctx, ld)
			if nil != err {
				return nil, err
			}

			if success {
				return ld, nil
			}
		}
	}
}

func (rlock *redisLock) Unlock(ctx context.Context, ld *LockData,
) error {
	var success, err = rlock.tryUnlock(ctx, ld)
	if nil != err {
		return err
	}
	if !success {
		log.Printf("Key %s is already unlocked or does not exist.\n", ld.Key)
	}
	return nil
}

func (rlock *redisLock) Close() error {
	return nil
}

func (rlock *redisLock) tryLock(ctx context.Context, ld *LockData) (bool, error) {
	var cmd = rlock.client.SetNX(ctx, ld.Key, ld.Value, rlock.config.MaxLockTime)
	return cmd.Result()
}

func (rlock *redisLock) tryUnlock(ctx context.Context, ld *LockData) (bool, error) {
	curVal, err := rlock.client.Get(ctx, ld.Key).Result()
	if nil != err {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if curVal == ld.Value {
		delCount, err := rlock.client.Del(ctx, ld.Key).Result()
		if nil != err && err != redis.Nil {
			return false, err
		}

		return delCount > 0, nil
	}

	if len(curVal) > 0 {
		log.Printf("Key %s is locked by another value: %s\n", ld.Key, curVal)
	}
	return false, nil
}

// MockLock là mock implementation cho test
type MockLock struct{}

// NewMockLock tạo mock lock mới
func NewMockLock() IDLock {
	return &MockLock{}
}

func (m *MockLock) Lock(ctx context.Context, key string) (*LockData, error) {
	return &LockData{Key: key, Value: "mock-lock"}, nil
}

func (m *MockLock) Unlock(ctx context.Context, ld *LockData) error {
	return nil
}
