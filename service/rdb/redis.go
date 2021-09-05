package rdb

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/wuzehv/passport/util/config"
	"log"
	"time"
)

var Rdb *redis.Pool

func init() {
	Rdb = &redis.Pool{
		MaxIdle:     config.Redis.MaxIdleConn,
		MaxActive:   config.Redis.MaxActiveConn,
		IdleTimeout: config.Redis.MaxConnIdleTimeout * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Redis.Host)
			if err != nil {
				log.Fatalf("redis init error: %v\n", err)
				return nil, err
			}

			if config.Redis.Password != "" {
				if _, err := c.Do("AUTH", config.Redis.Password); err != nil {
					c.Close()
					log.Fatalf("redis auth error: %v\n", err)
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", config.Redis.DbNum); err != nil {
				c.Close()
				log.Fatalf("redis select db error: %v\n", err)
				return nil, err
			}

			if _, err := c.Do("PING"); err != nil {
				log.Fatalf("redis ping error: %v\n", err)
				return nil, err
			}

			return c, nil
		},
	}
}

func SetJson(k string, v interface{}, expiration int) (reply interface{}, err error) {
	conn := Rdb.Get()
	defer conn.Close()

	str, _ := json.Marshal(v)
	return conn.Do("SET", k, str, "EX", expiration)
}

func GetJson(k string, v interface{}) bool {
	conn := Rdb.Get()
	defer conn.Close()

	cache, err := conn.Do("GET", k)
	if err != nil {
		log.Printf("redis error: %v\n", err)
		return false
	}

	if cache != nil && cache != "" {
		json.Unmarshal(cache.([]byte), &v)
		return true
	}

	return false
}
