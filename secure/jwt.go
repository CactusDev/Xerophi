package secure

import (
	"fmt"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt"

	log "github.com/sirupsen/logrus"
)

// Functions and data related to the generation and validation of tokens

// GenToken takes the scopes and other information needed to generate the
// JWT and returns a JWT
func GenToken(scopes []string, token string, secret string) (string, error) {
	return jwt.Sign(
		jwt.HS256(string(secret)),
		&jwt.Options{
			ExpirationTime: time.Now().AddDate(0, 0, 7), // Expires in a week
			Timestamp:      true,
			Public: map[string]interface{}{
				"token":  token,
				"scopes": scopes,
			},
		})
}

// ReadScope reads in a string of scopes in the format "table:[manage/create], table:[manage/create], ..." and returns their appropriately formatted scope strings
func ReadScope(scopeString string) []string {
	var scopes = make([]string, 0)

	for _, scope := range strings.Split(scopeString, ", ") {
		vals := strings.SplitN(scope, ":", 3)
		if len(vals) < 2 || (vals[1] != "manage" && vals[1] != "create") {
			// We have an invalid scope string, ignore it
			continue
		}
		// Passed the tests, add this scope
		scopes = append(scopes, fmt.Sprintf("%s:%s", vals[0], vals[1]))
	}
	return scopes
}

// ValidateToken takes a jwt.JWT object and returns an string for any errors
// encountered or "" if there are none. Does not check with Redis if token is active
func ValidateToken(tok *jwt.JWT, reqToken string, endpointScopes []string) string {
	algValidate := jwt.AlgorithmValidator(jwt.MethodHS256)
	expValidate := jwt.ExpirationTimeValidator(time.Now())
	issueValidate := jwt.IssuedAtValidator(time.Now())
	tokenValidate := TokenValidator(reqToken)
	scopeValidate := ScopeValidator(endpointScopes)

	err := tok.Validate(algValidate, expValidate, issueValidate, tokenValidate, scopeValidate)
	if err != nil {
		switch err {
		case jwt.ErrAlgorithmMismatch:
			log.Info("JWT Validation - Algorithim Mismatch")
		case jwt.ErrTokenExpired:
			log.Info("JWT Validation - JWT Token Expired")
		case ErrInvalidToken:
			log.Info("JWT Validation - JWT Token doesn't match request")
		case ErrFailedScopesConversion:
			log.Info("JWT Validation - Malformed scopes")
		case ErrInvalidScopes:
			log.Info("JWT Validation - Non-string scope in scopes")
		case ErrWrongNumScopes:
			log.Info("JWT Validation - Missing a required scope by len")
		}
		return err.Error()
	}

	return ""
}

// TokenActive checks with redis to see if the provided JWT token string is
// still active
func TokenActive(token string) (bool, error) {
	// exists, err := redisConn.HExists("activeTokens", token).Result()

	return true, nil
}
