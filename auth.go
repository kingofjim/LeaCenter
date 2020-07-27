package main

import (
	"github.com/gin-gonic/gin"
)

func AuthToken(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Query("token")
		if clientToken != token {
			c.AbortWithStatus(401)
		}
	}
}