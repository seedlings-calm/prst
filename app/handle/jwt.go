package handle

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/app/models"
	jwt "github.com/seedlings-calm/prst/middleware"
)

// AuthInit jwt验证new
func JWTInit() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "prst api",
		Key:             []byte("prst"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		PayloadFunc:     payloadFunc(),
		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Authorizator:    authorizator(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return jwt.MapClaims{
				"identityKey": v.UserName,
			}
		}
		return jwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &models.User{
			UserName: claims["identityKey"].(string),
		}
	}
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals models.User
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		username := loginVals.UserName
		password := loginVals.Password

		if (username == "admin" && password == "admin") || (username == "test" && password == "test") {
			return map[string]interface{}{
				"UserName":  username,
				"LastName":  "Bo-Yi",
				"FirstName": "Wu",
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*models.User); ok && v.UserName == "admin" {
			return true
		}
		return false
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}
