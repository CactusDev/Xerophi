package redis

import (
	"time"

	"github.com/gbrlsnchs/jwt"

	log "github.com/sirupsen/logrus"
)

// CacheUserLogin handles adding a login record to Redis
func CacheUserLogin(keyString string) error {
	// Try to convert the string provided to a JWT object
	key, err := jwt.FromString(keyString)
	if err != nil {
		log.Error(err)
		return err
	}

	// Try to set the redis key
	log.Debugf("Cacheing login for %s: %s", key.Public()["token"], key.String())
	// Have to negate the time.Since result because it's a future time wooo spooky
	log.Debugf("Expiration: %s", -time.Since(key.ExpirationTime()))

	// Add the JWT token to Redis under the token key, expiring in the future
	err = RedisConn.Session.Set(
		key.Public()["token"].(string),
		key.String(),
		-time.Since(key.ExpirationTime()),
	).Err()

	// Redis goofed
	if err != nil {
		log.Error(err)
		return err
	}

	// Nothing failed
	return nil
}
