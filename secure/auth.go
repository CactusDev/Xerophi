package secure

import (
	"html"
	"net/http"
	"strings"

	"github.com/CactusDev/Xerophi/redis"
	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/user"
	"github.com/CactusDev/Xerophi/util"

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
// request password is correct. Return false and nil means an internal error
func validateUser(token string, password string) (bool, error) {
	filter := map[string]interface{}{"token": token}

	// Check if this requested user (by token) exists
	fromDB, err := rethink.RethinkConn.GetSingle("users", filter)
	if err != nil {
		log.Error(err)
		return false, nil
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
		return false, nil
	}
	// Failed hash verification
	if !success {
		return false, AuthError{"Invalid username or password"}
	}

	// Passed all checks and no errors, valid username/password
	return true, nil
}

// Authenticator handles authentication requests for JWT auth
func Authenticator(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))

	// Get the request values
	requestVals := AuthRequestVals{}
	if err := ctx.BindJSON(&requestVals); err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
	}

	valid, err := validateUser(token, requestVals.Password)
	if err == nil && !valid {
		util.NiceError(
			ctx, AuthError{"Internal server error"}, http.StatusInternalServerError)
		return
	} else if err != nil && !valid {
		util.NiceError(ctx, err, http.StatusForbidden)
		return
	}

	// Split the request scopes up into a slice of valid scope strings
	scopes := ReadScope(strings.Join(requestVals.Scopes, ", "))

	// Generate the new token
	newToken, expiration, err := GenToken(scopes, token)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Authorize this new token
	err = redis.CacheUserLogin(newToken)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Return the success and the new token
	response := map[string]interface{}{
		"expiresAt": expiration,
		"token":     newToken,
	}
	ctx.JSON(http.StatusOK, response)
}
