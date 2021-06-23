package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/btowers/blog-go/pkg/auth"
	"github.com/gin-gonic/gin"
)

var identityKey = "email"

func JWT(aut auth.Service) *jwt.GinJWTMiddleware {

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "test zone",
		Key:            []byte("secret key"),
		Timeout:        time.Hour,
		MaxRefresh:     time.Hour,
		IdentityKey:    identityKey,
		SendCookie:     true,
		SecureCookie:   false, // non HTTPS dev environments
		CookieHTTPOnly: true,  // JS can't modify
		CookieDomain:   "localhost:8080",
		CookieName:     "token", // default jwt
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*auth.User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &auth.User{
				Email: claims["email"].(string),
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals auth.User
			if err := c.ShouldBind(&loginVals); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			user, err := aut.Login(loginVals)
			if err != nil {
				return nil, err
			}

			c.Set("user", user)
			return &loginVals, nil
		},

		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			user, _ := c.Get("user")

			c.JSON(http.StatusOK, gin.H{
				"authorized": true,
				"user":       user,
			})
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {

			payload, ok := c.Get("JWT_PAYLOAD")
			if !ok {
				return false
			}
			userClaim, _ := json.Marshal(payload)

			var userr auth.User
			json.Unmarshal(userClaim, &userr)

			user, err := aut.GetUser(userr.Email)
			if err != nil {
				return false
			}

			c.Set("user", user)
			return true

		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusOK, gin.H{
				"authorized": false,
				"user":       nil,
			})
		},

		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
