package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var version string = os.Getenv("BUILD_VER")

func main() {
	isDebugMode := flag.Bool("debug", false, "run in debug mode")
	flag.Parse()

	if *isDebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(versionMiddleware(version))

	if *isDebugMode {
		r.Use(cors.Default())
	} else {
		r.Use(cors.New(cors.Config{
			AllowOrigins:  []string{},
			AllowMethods:  []string{"GET", "POST"},
			AllowHeaders:  []string{"X-Client-Name", "X-Client-Version"},
			ExposeHeaders: []string{"X-API-Version"},
		}))
	}

	r.POST("/login", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	api := r.Group("/api")

	// Authorization
	api.Use(authorizationMiddleware)

	api.POST("/word", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	api.POST("/picture", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.Run()
}

func isPreflightRequest(c *gin.Context) bool {
	return c.Request.Method == http.MethodOptions &&
		c.Request.Header.Get("Access-Control-Request-Method") != "" &&
		c.Request.Header.Get("Access-Control-Request-Headers") != ""
}

func validAuthorization(c *gin.Context) bool {
	return c.Request.Header.Get("Authorization") != ""
}

func versionMiddleware(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-API-Version", version)
		c.Next()
	}
}

func authorizationMiddleware(c *gin.Context) {
	// Do not enforce authorization on preflight requests
	if isPreflightRequest(c) {
		c.Next()
		return
	}

	// Enforce authorization header
	if validAuthorization(c) {
		c.Next()
		return
	}

	// Abort all requests
	c.Status(http.StatusUnauthorized)
	c.Abort()
}
