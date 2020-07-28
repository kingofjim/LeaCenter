package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AuthToken(token string, originToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Query("token")
		fmt.Println(clientToken)
		fmt.Println(originToken)
		if clientToken != token {
			if clientToken !=  originToken {
				c.AbortWithStatus(401)
			} else {
				fmt.Println(c.ClientIP())
			}
		}
	}
}