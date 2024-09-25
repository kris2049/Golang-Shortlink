package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alextanhongpin/base62"
	"github.com/go-redis/redis/v8"
)

const (
	// URLIDKEY is global counter
	URLIDKEY = "next.url.id"

	// ShortlinkKey mapping the shortlink to the url
	ShortlinkKey = "shortlink:%s:url"

	// URLHashKey mapping the hash of the url to the shortlink
	URLHashKey = "urlhash:%s:url"

	//ShortlinkDetailKey mapping the shortlink to the detail of url
	ShortlinkDetailKey = "shortlink%s:detail"
)

type RedisClient struct {
	cli *redis.Client
}

// URLDetail contains the detail of the shortlink
type URLDetail struct {
	URL                 string        `json:"url"`
	CreateAt            string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

// NewRedisClient create a redis client
func NewRedisClient(add string, passwd string, db int) *RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     add,
		Password: passwd,
		DB:       db,
	})

	if pong, err := c.Ping(context.Background()).Result(); err != nil {
		panic(err)
	} else {
		log.Printf("Redis connected: %s", pong)
	}

	return &RedisClient{cli: c}
}

// Shorten convert url to shortlink
func (r *RedisClient) Shorten(url string, exp int64) (string, error) {
	ctx := context.Background()
	// conver url to sha1 hash
	h := toSha1(url)

	// fetch it if the url is cached
	s := fmt.Sprintf(URLHashKey, h)
	fmt.Print(s)
	d, err := r.cli.Get(ctx, fmt.Sprintf(URLHashKey, h)).Result()
	if err == redis.Nil {
		// not existed, do nothing
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			// expiration, do nothing
		} else {
			return d, nil
		}
	}

	// increase the global counter
	err = r.cli.Incr(ctx, URLIDKEY).Err()
	if err != nil {
		return "", err
	}

	// encode global counter to base62
	id, err := r.cli.Get(ctx, URLIDKEY).Int64()
	if err != nil {
		return "", err
	}
	eid := base62.Encode(uint64(id))

	// store the url against this encoded id
	err = r.cli.Set(ctx, fmt.Sprintf(ShortlinkKey, eid), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	// store the url against the hash of it
	err = r.cli.Set(ctx, fmt.Sprintf(URLHashKey, h), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	detail, err := json.Marshal(
		&URLDetail{
			URL:                 url,
			CreateAt:            time.Now().String(),
			ExpirationInMinutes: time.Duration(exp)},
	)

	if err != nil {
		return "", err
	}

	// store the url detail against this encoded id
	err = r.cli.Set(ctx, fmt.Sprintf(ShortlinkDetailKey, eid), string(detail), time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	return eid, nil
}

func toSha1(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// ShortlinkInfo returns the details of the shortlink
func (r *RedisClient) ShortlinkInfo(eid string) (interface{}, error) {
	d, err := r.cli.Get(context.Background(), fmt.Sprintf(ShortlinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{404, errors.New("unknow short URL")}
	} else if err != nil {
		return "", err
	} else {
		return d, nil
	}
}

// Unshorten convert shortlink to url
func (r *RedisClient) Unshorten(eid string) (string, error) {
	url, err := r.cli.Get(context.Background(), fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{404, err}
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}
