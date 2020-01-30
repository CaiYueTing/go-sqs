package redishelper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	c redis.Conn
}

func NewRedisPool(url string) (*Redis, error) {
	c, err := redis.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	return &Redis{c}, nil
}

func (r *Redis) Close() error {
	return r.c.Close()
}

func (r *Redis) Delete(key string) error {
	result, err := r.c.Do("del", key)
	if err != nil {
		return err
	}
	fmt.Println("delete key success", result)
	return nil
}

func (r *Redis) SetString(key string, value string) error {
	result, err := r.c.Do("SET", key, value)
	if err != nil {
		return err
	}
	fmt.Println("SetString success", result)

	return nil
}

func (r *Redis) ReadString(key string) (*string, error) {
	result, err := redis.String(r.c.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Redis) SetMap(key string, value map[string]string) error {
	v, _ := json.Marshal(value)
	result, err := r.c.Do("set", key, v)
	if err != nil {
		return err
	}
	fmt.Println("set map success", result)
	return nil
}

func (r *Redis) ReadMap(key string) (*map[string]string, error) {
	result, err := redis.Bytes(r.c.Do("get", key))
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
	result, err := r.c.Do(direction, key, value)
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
	result, err := redis.Strings(r.c.Do("lrange", key, start, end))
	if err != nil {
		return nil, err
	}
	return &result, nil
}
