package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	// sorted set functions
	DeleteSortedSet(context.Context, string) error
	AddMemberToSortedSet(context.Context, string, float64, string) error
	GetARangeOfMembers(context.Context, string, int64, int64) ([]string, error)
	GetAndRemoveNextMemberFromSortedSet(context.Context, string) (string, error)
	RemoveASpecificMemberFromSortedSet(context.Context, string, string) error

	// set functions
	GetValue(context.Context, string) (string, error)
	SetKeyValue(context.Context, string, any, time.Duration) error

	// hash functions
	CreateHash(context.Context, string, any) error
	GetHashValue(context.Context, string) (map[string]string, error)
}

type redisStruct struct {
	client *redis.Client
}

func NewRedisInterface(client *redis.Client) RedisInterface {
	return &redisStruct{
		client: client,
	}
}

func (r *redisStruct) DeleteSortedSet(ctx context.Context, sortedSetName string) error {
	err := r.client.Del(ctx, sortedSetName).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *redisStruct) AddMemberToSortedSet(
	ctx context.Context,
	sortedSetName string,
	score float64,
	data string,
) error {
	err := r.client.ZAdd(ctx, sortedSetName, redis.Z{Score: score, Member: data}).Err()
	if err != nil {
		return err
	}
	return err
}

func (r *redisStruct) GetARangeOfMembers(
	ctx context.Context,
	sortedSetName string,
	start,
	end int64,
) ([]string, error) {
	result, err := r.client.ZRangeArgs(
		ctx,
		redis.ZRangeArgs{
			Key:     sortedSetName,
			Start:   start,
			Stop:    end,
			ByScore: true,
			Offset:  0,
		},
	).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *redisStruct) GetAndRemoveNextMemberFromSortedSet(
	ctx context.Context,
	sortedSetName string,
) (string, error) {

	d, err := r.client.ZPopMax(ctx, sortedSetName, 1).Result()

	if err != nil || len(d) == 0 {
		return "", err
	}

	return d[0].Member.(string), nil
}

func (r *redisStruct) RemoveASpecificMemberFromSortedSet(
	ctx context.Context,
	sortedSetName,
	member string,
) error {
	err := r.client.ZRem(ctx, sortedSetName, member).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisStruct) GetValue(
	ctx context.Context,
	key string,
) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *redisStruct) SetKeyValue(
	ctx context.Context,
	key string,
	value any,
	expirationTime time.Duration,
) error {
	err := r.client.Set(ctx, key, value, expirationTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisStruct) CreateHash(
	ctx context.Context,
	hash string,
	data any,
) error {
	err := r.client.HSet(ctx, hash, data).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisStruct) GetHashValue(ctx context.Context, hash string) (map[string]string, error) {
	result, err := r.client.HGetAll(ctx, hash).Result()

	if err != nil {
		fmt.Println(fmt.Errorf("error occured: %s", err))
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("hash not found!")
	}

	return result, nil
}
