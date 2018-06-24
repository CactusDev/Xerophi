package secure

import (
	"errors"
	"testing"
	"time"

	"github.com/CactusDev/Xerophi/redis"
	"github.com/gbrlsnchs/jwt"
	goRedis "github.com/go-redis/redis"
)

func TestTokenValidator(t *testing.T) {
	testCases := []struct {
		token    string
		genToken string
		err      error
	}{
		{"token", "token", nil},
		{"token", "wrong", ErrInvalidToken},
		{"token", "tokenButLonger", ErrInvalidToken},
	}

	for _, test := range testCases {
		validateFunc := TokenValidator(test.token)
		jwtStr, _ := jwt.Sign(jwt.HS256("secret"),
			&jwt.Options{Public: map[string]interface{}{"token": test.genToken}})
		jot, _ := jwt.FromString(jwtStr)

		if err := validateFunc(jot); err != test.err {
			t.Errorf("Unexpected error returned. Expected %v, got %v", test.err, err)
		}
	}
}

func TestScopeValidator(t *testing.T) {
	testCases := []struct {
		token     string
		reqScopes []string
		jwtScopes []string
		err       error
	}{
		{
			"token", []string{"command:create"}, []string{"command:create"}, nil},
		{
			"token", []string{"command:create", "command:manage"},
			[]string{"command:create", "command:manage"}, nil,
		},
		{
			"token", []string{"command:create", "command:manage"},
			[]string{"command:create"}, ErrWrongNumScopes,
		},
		{
			"token", []string{"command:create", "command:manage"},
			[]string{"command:create", "foo:manage"}, errors.New("jwt: Missing required scope(s): [ command:manage ]"),
		},
	}

	for _, test := range testCases {
		validateFunc := ScopeValidator(test.reqScopes)
		jwtStr, _ := jwt.Sign(jwt.HS256("secret"),
			&jwt.Options{Public: map[string]interface{}{
				"token": test.token, "scopes": test.jwtScopes}})
		jot, _ := jwt.FromString(jwtStr)

		if err := validateFunc(jot); err != nil && err.Error() != test.err.Error() {
			t.Errorf("Unexpected error returned. Expected \"%v\", got \"%v\"",
				test.err, err)
		}
	}
	// Separate test for passing a string instead of a slice of strings for scopes
	validateFunc := ScopeValidator([]string{"command:manage"})
	jwtStr, _ := jwt.Sign(jwt.HS256("secret"),
		&jwt.Options{Public: map[string]interface{}{
			"token": "token", "scopes": "command:manage"}})
	jot, _ := jwt.FromString(jwtStr)

	if err := validateFunc(jot); err != ErrFailedScopesConversion {
		t.Errorf("Unexpected error returned. Expected \"%v\", got \"%v\"",
			ErrFailedScopesConversion, err)
	}
}

func TestActiveValidator(t *testing.T) {
	testCases := []struct {
		doAdd         bool
		token         string
		scopes        []string
		expectedError error
	}{
		{
			doAdd: true,
			token: "testToken",
			scopes: []string{
				"test:test",
			},
			expectedError: nil,
		},
	}

	// Instantiate redis server connection since validators use that
	redis.RedisConn = &redis.Connection{
		DB: 0,
		Opts: redis.ConnectionOpts{
			Host:     "localhost",
			Port:     6379,
			User:     "",
			Password: "",
		},
	}

	if err := redis.RedisConn.Connect(); err != nil {
		// Can't continue without redis running
		t.Fatal("Redis connection failed - ", err)
	}

	for _, test := range testCases {
		// Generate the JWT token
		jwtToken, err := jwt.Sign(
			jwt.HS256("secret"),
			&jwt.Options{
				ExpirationTime: time.Now(),
				Timestamp:      true,
				Public: map[string]interface{}{
					"token":  test.token,
					"scopes": test.scopes,
				},
			})
		if err != nil {
			t.Error(err)
		}

		// Convert the JWT token string to a JWT object again
		jwtTokenObj, err := jwt.FromString(jwtToken)
		if err != nil {
			t.Error(err)
		}

		if test.doAdd {
			// Add the token to redis
			err = redis.RedisConn.Session.Set(
				jwtTokenObj.Public()["token"].(string),
				jwtToken,
				-time.Since(jwtTokenObj.ExpirationTime()),
			).Err()
			if err != nil && err != goRedis.Nil {
				t.Error(err)
			}
		}

		validateFunc := ActiveValidator(test.token)
		err = validateFunc(jwtTokenObj)
		if err != nil && err != test.expectedError {
			t.Error(err)
		}

		// Only attempt to remove if we added it earlier
		if test.doAdd {
			// Current test complete, remove the token added
			err = redis.RedisConn.Session.Del(jwtTokenObj.Public()["token"].(string)).Err()
			if err != nil && err != goRedis.Nil {
				t.Error(err)
			}
		}
	}
}
