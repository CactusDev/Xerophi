package secure

import (
	"errors"
	"testing"

	"github.com/gbrlsnchs/jwt"
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
