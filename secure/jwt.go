package secure

import (
	"fmt"
	"strings"
	"time"

	"github.com/CactusDev/Xerophi/redis"
	"github.com/gbrlsnchs/jwt"

	log "github.com/sirupsen/logrus"
)

// Functions and data related to the generation and validation of tokens
var secret string

// SetSecret sets the module-wide local variable for keeping track of the secret
// which is used for JWT token generation
func SetSecret(setSecret string) {
	secret = setSecret
}

// GenToken takes the scopes and other information needed to generate the
// JWT and returns a JWT and when it expires
func GenToken(scopes []string, token string) (string, string, error) {
	expr := time.Now().AddDate(0, 0, 7) // Expires in a week
	token, err := jwt.Sign(
		jwt.HS256(secret),
		&jwt.Options{
			ExpirationTime: expr,
			Timestamp:      true,
			Public: map[string]interface{}{
				"token":  token,
				"scopes": scopes,
			},
		})
	return token, expr.String(), err
}

// ReadScope reads in a string of scopes in the format "table:[manage/create],
// table:[manage/create], ..." and returns their appropriately formatted scope
// strings
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
// encountered or "" if there are none. Does not check with Redis if token is
// active
func ValidateToken(tok *jwt.JWT, reqToken string, scopes []string) string {
	// The algorithim matches HS256
	algValidate := jwt.AlgorithmValidator(jwt.MethodHS256)
	// The token hasn't expired
	expValidate := jwt.ExpirationTimeValidator(time.Now())
	// The token was issued at a valid time
	issueValidate := jwt.IssuedAtValidator(time.Now())
	// The token in the request matches the token in the JWT token
	tokenValidate := TokenValidator(reqToken)
	// The scopes in the JWT token match the required scopes for the endpoint
	scopeValidate := ScopeValidator(scopes)
	// The token is currently active
	activeValidate := ActiveValidator()

	err := tok.Validate(
		algValidate, expValidate, issueValidate, // Lib validators
		tokenValidate, scopeValidate, activeValidate, // Local validators
	)
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
		case ErrMissingToken:
			log.Info("JWT Validation - Token doesn't exist or has expired")
		case ErrInternalError:
			log.Error("Internal error happened. Fun stuff.")
		}
		return err.Error()
	}

	return ""
}

// TokenActive checks with redis to see if the provided JWT token string is
// still active for the provided user
func TokenActive(user string, token string) bool {
	redisToken, err := redis.RedisConn.Session.Get(user).Result()
	if err != nil {
		log.Error(err)
		return false
	}

	return redisToken == token
}
