package secure

import (
	"testing"
	"time"

	"github.com/gbrlsnchs/jwt"
)

func TestGenToken(t *testing.T) {
	// TODO: Add more tests for this
	now := time.Now().Format("2006-01-02 15:04:05 -0700 MST")
	expiration := time.Now().AddDate(0, 0, 7).Format("2006-01-02 15:04:05 -0700 MST")
	token := "testing"
	secret := "testSecret"
	scopes := []string{"commands:add", "quotes:manage", "user:manage"}
	jwtTok, err := GenToken(scopes, token, secret)
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
		res := ReadScope(test.input)
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
