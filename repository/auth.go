package repository

import (
	"context"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/acoshift/acourse/appctx"
)

// StoreMagicLink stores magic link to redis
func StoreMagicLink(ctx context.Context, linkID string, userID string) error {
	pool, prefix := appctx.GetRedisPool(ctx)
	db := pool.Get()
	defer db.Close()

	_, err := db.Do("SETEX", prefix+"magic:"+linkID, int64(time.Hour/time.Second), userID)
	if err != nil {
		return err
	}
	return nil
}

// FindMagicLink finds magic link from redis
func FindMagicLink(ctx context.Context, linkID string) (string, error) {
	pool, prefix := appctx.GetRedisPool(ctx)
	db := pool.Get()
	defer db.Close()

	key := prefix + "magic:" + linkID
	userID, err := redis.String(db.Do("GET", key))
	if err == redis.ErrNil {
		return "", appctx.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	db.Do("DEL", key)
	return userID, nil
}

// CanAcquireMagicLink checks rate limit to acquire magic link
func CanAcquireMagicLink(ctx context.Context, email string) (bool, error) {
	pool, prefix := appctx.GetRedisPool(ctx)
	db := pool.Get()
	defer db.Close()

	key := prefix + "magic-rate:" + email
	current, err := redis.Int(db.Do("INCR", key))
	if err != nil {
		return false, err
	}
	if current > 1 {
		return false, nil
	}
	_, err = db.Do("EXPIRE", key, int64(5*time.Minute/time.Second))
	if err != nil {
		return false, err
	}
	return true, nil
}
