package middleware

import "github.com/gin-gonic/gin"

func InitMiddleWare(r *gin.Engine) {
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	r.Use(RequestId(XRequestId))
	r.Use(Options()).Use(Secure)

}
