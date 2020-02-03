package redishelper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	pool *redis.Pool
}

func NewRedisPool(url string) (*Redis, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", url)
		},
	}
	c := pool.Get()
	defer c.Close()
	return &Redis{pool: pool}, nil
}

func (r *Redis) Delete(key string) error {
	conn := r.pool.Get()
	result, err := conn.Do("del", key)
	if err != nil {
		return err
	}
	fmt.Println("delete key success", result)
	return nil
}

func (r *Redis) SetString(key string, value string) error {
	conn := r.pool.Get()
	result, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	fmt.Println("SetString success", result)

	return nil
}

func (r *Redis) ReadString(key string) (*string, error) {
	conn := r.pool.Get()
	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Redis) SetMap(key string, value map[string]string) error {
	conn := r.pool.Get()
	v, _ := json.Marshal(value)
	result, err := conn.Do("set", key, v)
	if err != nil {
		return err
	}
	fmt.Println("set map success", result)
	return nil
}

func (r *Redis) ReadMap(key string) (*map[string]string, error) {
	conn := r.pool.Get()
	result, err := redis.Bytes(conn.Do("get", key))
	if err != nil {
		return nil, err
	}

	var m map[string]string
	err = json.Unmarshal(result, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Redis) push(direction string, key string, value string) error {
	conn := r.pool.Get()
	result, err := conn.Do(direction, key, value)
	if err != nil {
		return err
	}
	fmt.Println("lpush success", result)
	return nil
}

func (r *Redis) PushList(direction string, key string, value string) error {
	switch direction {
	case "lpush":
		return r.push(direction, key, value)
	case "rpush":
		return r.push(direction, key, value)
	default:
		return errors.New("not command")
	}
}

func (r *Redis) ReadList(key string, start int, end int) (*[]string, error) {
	conn := r.pool.Get()
	result, err := redis.Strings(conn.Do("lrange", key, start, end))
	if err != nil {
		return nil, err
	}
	return &result, nil
}
