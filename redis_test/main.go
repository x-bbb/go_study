package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

// redis 连接池
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:   100,
		MaxActive: 25,
	}
}

func main() {
	pool = newPool("172.16.0.20:6379", "123")
	conn := pool.Get()
	r, err := redis.String(conn.Do("GET", "name"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(r)
}
