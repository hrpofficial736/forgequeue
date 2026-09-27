package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	// sorted set functions
	DestroySortedSet(context.Context, string) error
	AddToSortedSet(context.Context, string, float64, any) error
	GetMembersWithASpecificScore(context.Context, string, float64, float64, float64) ([]string, error)
	GetAndRemoveNextMemberFromSortedSet(context.Context, string) (string, error)
	RemoveASpecificMemberFromSortedSet(context.Context, string, string) error

	// hash functions
	CreateHash(context.Context, string, any, time.Duration) error
	GetHashValue(context.Context, string) (*redis.MapStringStringCmd, error)
}

type redisStruct struct {
	client *redis.Client
}

func NewRedisInterface(client *redis.Client) RedisInterface {
	return &redisStruct{
		client: client,
	}
}

func (r *redisStruct) DestroySortedSet(ctx context.Context, sortedSetName string) error {
	err := r.client.Del(ctx, sortedSetName).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *redisStruct) AddToSortedSet(
	ctx context.Context,
	sortedSetName string,
	score float64,
	data any,
) error {
	err := r.client.ZAdd(ctx, sortedSetName, redis.Z{Score: score, Member: data}).Err()
	if err != nil {
		return err
	}
	return nil
}

func(r *redisStruct) GetMembersWithASpecificScore(
	ctx context.Context,
	sortedSetName string,
	score float64,
	start float64,
	end float64,
) ([]string, error) {
	list, err := r.client.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key: sortedSetName,
		Start: score,
		Stop: score,
		ByScore: true,
		Offset: 0,
		Count: -1,
	}).Result()

	if err != nil {
		return nil, err
	}

	return list, nil
}


func (r *redisStruct) GetAndRemoveNextMemberFromSortedSet(
	ctx context.Context,
	sortedSetName string,
) (string, error) {

	d, err := r.client.ZPopMax(ctx, sortedSetName, 1).Result()

	if err != nil {
		return "", err
	}

	return d[0].Member.(string), nil
}

func (r *redisStruct) RemoveASpecificMemberFromSortedSet(
	ctx context.Context,
	sortedSetName,
	member string,
) error {
	err := r.client.ZRem(ctx, sortedSetName).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisStruct) CreateHash(
	ctx context.Context,
	hash string,
	data any,
	expiration time.Duration,
) error {
	err := r.client.HSet(ctx, hash, data).Err()
	if err != nil {
		return err
	}

	err = r.client.Expire(ctx, hash, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisStruct) GetHashValue(ctx context.Context, hash string) (*redis.MapStringStringCmd, error) {
	result := r.client.HGetAll(ctx, hash)

	if result.Err() != nil {
		fmt.Println(fmt.Errorf("error occured: %s", result.Err()))
		return nil, result.Err()
	}

	return result, nil
}
