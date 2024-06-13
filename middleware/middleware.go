package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loyyal/k8s-deployer-go/common"
)

// func JWTAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		err := token.TokenValid(c)
// 		if err != nil {
// 			// c.String(http.StatusUnauthorized, "")
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"message": err.Error(),
// 			})
// 			c.Abort()
// 			return
// 		}

// 		c.Next()
// 	}
// }

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusAccepted)
			return
		}

		c.Next()
	}
}

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse basic auth credentials
		username, password, ok := ctx.Request.BasicAuth()

		if !ok {
			ctx.Header("WWW-Authenticate", `Basic realm="Please enter your username and password for access"`)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if username != "username" || password != "password" {
			// The user is not authorized.
			common.PrepareCustomError(ctx, 401, "", "No access key and secret found in the request", "basic authorisation failed ")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
