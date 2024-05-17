package middleware

import "github.com/gin-gonic/gin"

func InitMiddleWare(r *gin.Engine) {

	r.Use(RequestId(XRequestId))
	r.Use(gin.Recovery())
	r.Use(CheckOptions())
}
