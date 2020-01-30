package redishelper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func newRedisPool() (redis.Conn, error) {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Delete(key string) error {
	c, err := newRedisPool()
	if err != nil {
		return err
	}
	defer c.Close()
	result, err := c.Do("del", key)
	if err != nil {
		return err
	}
	fmt.Println("delete key success", result)
	return nil
}

func SetString(key string, value string) error {
	c, err := newRedisPool()
	if err != nil {
		return err
	}
	defer c.Close()
	result, err := c.Do("SET", key, value)
	if err != nil {
		return err
	}
	fmt.Println("SetString success", result)

	return nil
}

func ReadString(key string) (*string, error) {
	c, err := newRedisPool()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	result, err := redis.String(c.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func SetMap(key string, value map[string]string) error {
	c, err := newRedisPool()
	if err != nil {
		return err
	}
	defer c.Close()

	v, _ := json.Marshal(value)
	result, err := c.Do("set", key, v)
	if err != nil {
		return err
	}
	fmt.Println("set map success", result)
	return nil
}

func ReadMap(key string) (*map[string]string, error) {
	c, err := newRedisPool()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	result, err := redis.Bytes(c.Do("get", key))
	if err != nil {
		return nil, err
	}

	var r map[string]string
	err = json.Unmarshal(result, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func push(direction string, key string, value string) error {
	c, err := newRedisPool()
	if err != nil {
		return err
	}
	defer c.Close()

	result, err := c.Do(direction, key, value)
	if err != nil {
		return err
	}
	fmt.Println("lpush success", result)
	return nil
}

func PushList(direction string, key string, value string) error {
	switch direction {
	case "lpush":
		return push(direction, key, value)
	case "rpush":
		return push(direction, key, value)
	default:
		return errors.New("not command")
	}
}

func ReadList(key string, start int, end int) (*[]string, error) {
	c, err := newRedisPool()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	result, err := redis.Strings(c.Do("lrange", key, start, end))
	if err != nil {
		return nil, err
	}
	return &result, nil
}
