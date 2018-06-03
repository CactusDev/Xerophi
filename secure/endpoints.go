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
)

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
func validateUser(token string, password string, ctx *gin.Context) (bool, error) {
	filter := map[string]interface{}{"token": token}

	// Check if this requested user (by token) exists
	fromDB, err := rethink.RethinkConn.GetSingle(filter, "users")
	if err != nil {
		return false, err
	}

	// There were no values returned for that token
	if fromDB == nil {
		return false, AuthError{"Non-existant user"}
	}

	// Put the response from the DB into an object
	userVals := user.Database{}
	mapstruct.Decode(fromDB, &userVals)

	// Verify the password against the retrieved hash
	if !VerifyHash([]byte(userVals.Hash), password) {
		// Failed hash verification
		return false, AuthError{"Invalid username or password"}
	}

	// Passed all checks and no errors, valid username/password
	return true, nil
}

// Login handles requests to /user/:token/login endpoint for authentication
func Login(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))

	// Get the request values
	requestVals := AuthRequestVals{}
	if err := ctx.BindJSON(&requestVals); err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
	}

	if valid, err := validateUser(token, requestVals.Password, ctx); err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	} else if !valid {
		util.NiceError(
			ctx, AuthError{"Invalid username or password"}, http.StatusBadRequest)
		return
	}

	// Split the request scopes up into a slice of valid scope strings
	scopes := ReadScope(strings.Join(requestVals.Scopes, ", "))

	// Generate the new JWT token
	jwtToken, err := GenToken(scopes, token, "testfoobar")
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Store the newly generated JWT token
	err = redis.CacheUserLogin(jwtToken)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Return the JWT token to the user
	ctx.Header("X-Auth-New-Token", jwtToken)
	ctx.JSON(http.StatusOK, nil)
}
