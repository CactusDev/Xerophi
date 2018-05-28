package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

var (
	// RedisConn is the shareable Redis connection resource
	RedisConn *Connection
)

// ConnectionOpts is the data required to make a connection to redis
type ConnectionOpts struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Connection defines a connection to a Redis instance
type Connection struct {
	DB      int            // The redis server to use
	Opts    ConnectionOpts // Connection options for the redis server
	Session *redis.Client
}

// Connect connects you to Redis
func (c *Connection) Connect() error {
	client := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", c.Opts.Host, c.Opts.Port),
			Password: c.Opts.Password,
			DB:       c.Opts.DB,
		})

	// Ping the server to make sure we had a succesful connection
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	// Connection succeeded
	c.Session = client
	return nil
}
