package middleware

import "github.com/gin-gonic/gin"

var Mode = "dev"

func AllowAllOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Origin, Content-Type, Accept")
			c.Header("Access-Control-Allow-Credentials", "true")
			// c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(200)
		}
	}
}

func CheckOptions() gin.HandlerFunc {
	if Mode == "dev" {
		return AllowAllOptions()
	}
	return func(c *gin.Context) {
		//获取配置，限制访问  TODO:

	}
}
