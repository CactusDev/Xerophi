package secure

import (
	"net/http"
	"strings"

	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/user"
	"github.com/CactusDev/Xerophi/util"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"

	mapstruct "github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// Claims is an alias to make code clearer/more understandable
type Claims = map[string]interface{}

// AuthError is an error object used for raising errors during authentication
type AuthError struct {
	Message string
}

// Error allows AuthError to implement the error interface
func (a AuthError) Error() string {
	return a.Message
}

// AuthRequestVals is the values that are required in the login request
type AuthRequestVals struct {
	Password string   `json:"password" binding:"required"`
	Scopes   []string `json:"scopes" binding:"required"`
}

// Non-exported helper function for determining if a user exists and if the
// request password is correct
func validateUser(token string, password string) (bool, error) {
	filter := map[string]interface{}{"token": token}

	// Check if this requested user (by token) exists
	fromDB, err := rethink.RethinkConn.GetSingle("users", filter)
	if err != nil {
		return false, err
	}

	// There were no values returned for that token
	if fromDB == nil {
		return false, AuthError{"Non-existant user"}
	}

	// Put the response from the DB into an object
	var userVals user.Database
	mapstruct.Decode(fromDB, &userVals)

	// Verify the password against the retrieved hash
	success, err := VerifyHash(userVals.Hash, password, userVals.Salt)
	// An error occured, log it but don't show the users
	if err != nil {
		log.Error(err)
		return false, AuthError{"Internal server error"}
	}
	// Failed hash verification
	if !success {
		return false, AuthError{"Invalid username or password"}
	}

	// Passed all checks and no errors, valid username/password
	return true, nil
}

// Authenticator handles authentication requests for JWT auth
func Authenticator(username string, password string, ctx *gin.Context) (interface{}, bool) {
	if valid, err := validateUser(username, password); err != nil || !valid {
		log.Error(err)
		return nil, false
	}

	// Get the request values
	requestVals := AuthRequestVals{}
	if err := ctx.BindJSON(&requestVals); err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
	}

	// Split the request scopes up into a slice of valid scope strings
	scopes := ReadScope(strings.Join(requestVals.Scopes, ", "))

	// Success, return the valid scopes for use in the token
	return Claims{"scopes": scopes, "token": username}, true
}

// Authorizator handles the authorization of a request
func Authorizator(vals interface{}, ctx *gin.Context) bool {
	// Our Identity function should be returning a map
	authVals, ok := vals.(Claims)
	if !ok {
		log.Error("Non-map vals for authorizator: ", vals)
		return false
	}

	// Make sure this data is valid

	// Attempt to retrieve username from request
	usernameVal, ok := authVals["username"]
	if !ok {
		// Token isn't a key in the auth data
		log.Error("Missing username key")
		return false
	}
	username, ok := usernameVal.(string)
	if !ok {
		// Token isn't a string
		log.Error("Failed to convert username to string")
		return false
	}

	// Attempt to retrieve scopes from request
	scopesVal, ok := authVals["scopes"]
	if !ok {
		// Token isn't a key in the auth data
		log.Error("Missing scopes key")
		return false
	}
	scopes, ok := scopesVal.([]string)
	if !ok {
		// Token isn't a string
		log.Error("Failed to convert scopes to list of strings")
		return false
	}

	// Verify this is a currently active token

	// Verify the scopes in the token match the required ones for the endpoint

	return true
}

// Payload handles the filling of the token's claims, in this case the scopes
func Payload(data interface{}) jwt.MapClaims {
	// Should always be in this form
	claimsData, ok := data.(Claims)
	if !ok {
		// Don't want authentication to complete if the data is somehow invalid
		// So Fatal to panic and then have our middleware recover it for us
		log.WithField("data", data).Fatal("Failed to assert correct type during payload prep")
		// Well that's ... odd
		return jwt.MapClaims{}
	}

	scopes, ok := claimsData["scopes"].([]string)
	if !ok {
		// Don't want authentication to complete if the data is somehow invalid
		// So Fatal to panic and then have our middleware recover it for us
		log.WithField("scopes", claimsData).Fatal("Failed to assert correct type during payload prep")
		// Well that's ... odd
		return jwt.MapClaims{}
	}

	username, ok := claimsData["username"].(string)
	if !ok {
		// Don't want authentication to complete if the data is somehow invalid
		// So Fatal to panic and then have our middleware recover it for us
		log.WithField("scopes", claimsData).Fatal("Failed to assert correct type during payload prep")
		// Well that's ... odd
		return jwt.MapClaims{}
	}

	return jwt.MapClaims{"scopes": scopes, "username": username}
}

// Identity handles parsing the claims and setting the values we'll need
// for confirming if the user has the correct scopes to access the endpoint
func Identity(claims jwt.MapClaims) interface{} {
	// Attempt to get username and scopes from claim
	username, ok := claims["username"]
	if !ok {
		log.Error("Missing username in JWT token")
		// No values will cause validation failure
		return Claims{}
	}
	scopes, ok := claims["scopes"]
	if !ok {
		log.Error("Missing scopes in JWT token")
		// No values will cause validation failure
		return Claims{}
	}

	return Claims{"username": username, "scopes": scopes}
}
