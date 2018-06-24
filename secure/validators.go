package secure

import (
	"errors"
	"fmt"

	"github.com/CactusDev/Xerophi/redis"

	"github.com/gbrlsnchs/jwt"
	goRedis "github.com/go-redis/redis"
)

// Custom validator functions for JWT validation

var (
	// ErrInvalidToken - JWT token doesn't match request token
	ErrInvalidToken = errors.New(
		"jwt: JWT Token field doesn't match request token")
	// ErrFailedScopesConversion - Scopes passed should be a list of strings
	ErrFailedScopesConversion = errors.New(
		"jwt: Malformed scopes, should be list of strings")
	// ErrWrongNumScopes - Missing required scopes in request
	ErrWrongNumScopes = errors.New("jwt: Missing required scopes")
	// ErrInvalidScopes - Non-string scope was passed to ReadScope
	ErrInvalidScopes = errors.New("jwt: Non-string scope in scopes")
	// ErrMissingOrExpiredToken - Token doesn't exist/has expired
	ErrMissingOrExpiredToken = errors.New("jwt: Token doesn't exist or has expired")
	// ErrInternalError - Some error occured in the parsing
	ErrInternalError = errors.New("Internal error")
)

// TokenValidator validates the JWT token's token field matches our current one
func TokenValidator(token string) jwt.ValidatorFunc {
	return func(j *jwt.JWT) error {
		if j.Public()["token"] != token {
			return ErrInvalidToken
		}
		return nil
	}
}

// ScopeValidator validates the scopes claimed in the token against the provided ones that are required
func ScopeValidator(requiredScopes []string) jwt.ValidatorFunc {
	// Scope: empty struct map - allows for key lookup later
	scopes := make(map[string]struct{})
	missingScopes := ""

	return func(j *jwt.JWT) error {
		s := j.Public()["scopes"]
		jwtScopes, ok := s.([]interface{})
		if !ok {
			// Unable to convert scopes to the slice we need
			return ErrFailedScopesConversion
		}

		if len(requiredScopes) > len(jwtScopes) {
			// We have fewer scopes in the token than are required, obviously fail
			return ErrWrongNumScopes
		}

		// Add all the scopes from the JWT into the map with an empty struct
		for _, scope := range jwtScopes {
			scopeStr, ok := scope.(string)
			if !ok {
				// We have a non-string value somehow
				return ErrInvalidScopes
			}
			scopes[scopeStr] = struct{}{}
		}

		// Go through all the required scopes and check if they existed
		// within the JWT
		for _, scope := range requiredScopes {
			if _, ok := scopes[scope]; !ok {
				// We're missing a required scope, store it as "[scope] "
				missingScopes += scope + " "
			}
		}
		// There was at least one missing scope
		if len(missingScopes) > 0 {
			return errors.New("jwt: Missing required scope(s): [ " + missingScopes + "]")
		}

		return nil
	}
}

// ActiveValidator validates whether the token is currently active
func ActiveValidator(token string) jwt.ValidatorFunc {
	return func(j *jwt.JWT) error {
		// Check if the token exists in redis
		redisToken, err := redis.RedisConn.Session.Get(token).Result()

		if err != nil && err != goRedis.Nil {
			fmt.Println(err)
			return ErrInternalError
		}

		if redisToken != j.String() {
			return ErrMissingOrExpiredToken
		}

		return nil
	}
}
