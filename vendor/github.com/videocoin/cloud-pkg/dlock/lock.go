package dlock

import (
	"errors"

	lock "github.com/bsm/redis-lock"
	"github.com/go-redis/redis"
)

var (
	ErrObtainLock = errors.New("obtain lock")
)

type Locker struct {
	cli *redis.Client
}

func New(uri string) (*Locker, error) {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	opts.MaxRetries = 3
	opts.PoolSize = 50

	cli := redis.NewClient(opts)
	err = cli.Ping().Err()
	if err != nil {
		return nil, err
	}

	return &Locker{
		cli: cli,
	}, nil
}

func (l *Locker) Obtain(key string) (*lock.Locker, error) {
	return lock.Obtain(l.cli, key, nil)
}
