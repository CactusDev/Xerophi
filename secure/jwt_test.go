package secure_test

import (
	"testing"
	"time"

	"github.com/CactusDev/Xerophi/redis"
	"github.com/CactusDev/Xerophi/secure"

	"github.com/gbrlsnchs/jwt"
	goRedis "github.com/go-redis/redis"
)

func TestGenToken(t *testing.T) {
	secure.SetSecret("secret")

	now := time.Now().Format("2006-01-02 15:04:05 -0700 MST")
	expiration := time.Now().AddDate(0, 0, 7).Format("2006-01-02 15:04:05 -0700 MST")
	token := "testing"
	scopes := []string{"commands:add", "quotes:manage", "user:manage"}
	jwtTok, expr, err := secure.GenToken(scopes, token)
	if err != nil {
		t.Error("Unexpected error in token generation")
	}
	tok, err := jwt.FromString(jwtTok)
	if err != nil {
		t.Error("Unexpected error in reading token")
	}
	if tok.Algorithm() != "HS256" {
		t.Errorf("Invalid algorithim. Expecting HS256, got %s", tok.Algorithm())
	}
	if tok.ExpirationTime().String() != expiration {
		t.Errorf("Invalid expiration date. Expecting %sm got %s",
			expiration, tok.ExpirationTime())
	}
	if tok.ExpirationTime().String() != expr {
		t.Errorf("Invalid expiration date. Expecting %s got %s",
			expr, tok.ExpirationTime())
	}
	if tok.IssuedAt().String() != now {
		t.Errorf("Invalid timestamp. Expecting %s, got %s", now, tok.IssuedAt())
	}
	pub := tok.Public()
	if pub["token"] != token {
		t.Errorf("Invalid token. Expecting %s, got %s", token, pub["token"])
	}
	pubScopes, ok := pub["scopes"].([]interface{})
	if !ok {
		t.Error("Invalid scopes. Unable to convert to slice")
	}
	for pos, val := range pubScopes {
		if scopes[pos] != val {
			t.Errorf("Invalid scopes. Expecting %s at pos %d, got %s",
				scopes[pos], pos, val)
		}
	}
}

func TestReadScope(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"command:create", []string{"command:create"}},
		{
			input:    "command:create, quote:create, quote:manage",
			expected: []string{"command:create", "quote:create", "quote:manage"},
		},
		{
			input:    "command:create, quote:create, foo:manage:spam",
			expected: []string{"command:create", "quote:create", "foo:manage"},
		},
		{
			input:    "invalidscope",
			expected: []string{},
		},
		{
			input:    "user:manage, weird:create:but:sorta:valid, has:invalidRight totesInvalid",
			expected: []string{"user:manage", "weird:create"},
		},
	}

	for _, test := range testCases {
		res := secure.ReadScope(test.input)
		if len(test.expected) != len(res) {
			t.Errorf("Invalid output from ReadScope. Expected length %d, got length %d",
				len(test.expected), len(res))
		}
		// Have to have this massive ugliness because we have to just check if the
		// required scopes exist, don't have to be the exact same position
		for _, expected := range test.expected {
			exists := false
			for _, scope := range res {
				if scope == expected {
					exists = true
					break
				}
			}
			if !exists {
				t.Errorf("Missing required value %s", expected)
			}
		}
	}
}

func TestValidateToken(t *testing.T) {
	testCases := []struct {
		scopes     []string
		algorithim string
		expiration time.Time
		issuedAt   time.Time
		reqToken   string
		err        error
	}{
		{
			algorithim: jwt.MethodHS256,
			expiration: time.Now().AddDate(0, 0, 7),
			issuedAt:   time.Now(),
			scopes:     []string{"test:test"},
			reqToken:   "testToken",
			err:        nil,
		},
		{
			algorithim: jwt.MethodHS256,
			expiration: time.Now().AddDate(0, 0, 7),
			issuedAt:   time.Now(),
			scopes:     []string{"test:test", "second:test"},
			reqToken:   "multiScopesToken",
			err:        nil,
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
				ExpirationTime: test.expiration,
				Timestamp:      true,
				Public: map[string]interface{}{
					"token":  test.reqToken,
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

		// Add the token to redis
		err = redis.RedisConn.Session.Set(
			jwtTokenObj.Public()["token"].(string),
			jwtToken,
			-time.Since(jwtTokenObj.ExpirationTime()),
		).Err()
		if err != nil && err != goRedis.Nil {
			t.Error(err)
		}

		// Use secure.ValidateToken
		validateErr := secure.ValidateToken(jwtTokenObj, test.reqToken, test.scopes)
		if validateErr != test.err {
			t.Error(validateErr)
		}

		// Current test complete, remove the token added
		err = redis.RedisConn.Session.Del(jwtTokenObj.Public()["token"].(string)).Err()
		if err != nil && err != goRedis.Nil {
			t.Error(err)
		}
	}
}
