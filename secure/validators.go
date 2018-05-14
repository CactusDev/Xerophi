package secure

import (
	"errors"

	"github.com/gbrlsnchs/jwt"
)

// Custom validator functions for JWT validation

// TokenValidator validates the JWT token's token field matches our current one
func TokenValidator(token string) jwt.ValidatorFunc {
	return func(j *jwt.JWT) error {
		if j.Public()["token"] != token {
			return errors.New("jwt: JWT Token field doesn't match request token")
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
		jwtScopes, ok := j.Public()["scopes"].([]string)
		if !ok {
			// Unable to convert scopes to the slice we need
			return errors.New("jwt: Invalid scopes key, should be list of strings")
		}
		// Add all the scopes from the JWT into the map with an empty struct
		for _, scope := range jwtScopes {
			scopes[scope] = struct{}{}
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
