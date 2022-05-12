package cache

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

func init() {
	pool = newPool()
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "redis:6379")
			if err != nil {
				log.Printf("unipic: failed to dial redis: %v\n", err)
			}
			return c, err
		},
	}
}

func ping(c redis.Conn) error {
	_, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	return nil
}

// GetConn return connection
func GetConn() (redis.Conn, error) {
	conn := pool.Get()
	err := ping(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Set stores to redis, only string:string now
func Set(c redis.Conn, key string, value string) error {
	_, err := c.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

// Get executes redis GET command
func Get(c redis.Conn, key string) (string, error) {
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}

	return s, nil
}
