package main

import (
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/qwwqe/VocabKing/pkg/errors"
	"github.com/qwwqe/VocabKing/pkg/requests"
)

const (
	HeaderClientName    = "X-Client-Name"
	HeaderClientVersion = "X-Client-Version"
	HeaderServerVersion = "X-Server-Version"
)

var (
	version = "unset"

	isDebugMode = flag.Bool("debug", false, "run in debug mode")
	disableCORS = flag.Bool("disable-cors", false, "disable CORS enforcement")
)

func main() {
	flag.Parse()

	if *isDebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(
		versionMiddleware(version),
		enforceContentTypeMiddleware,
	)

	if *isDebugMode || *disableCORS {
		r.Use(cors.Default())
	} else {
		r.Use(cors.New(cors.Config{
			AllowOrigins:  strings.Fields(os.Getenv("ALLOWED_ORIGINS")),
			AllowMethods:  []string{"GET", "POST"},
			AllowHeaders:  []string{HeaderClientName, HeaderClientVersion},
			ExposeHeaders: []string{HeaderServerVersion},
		}))
	}

	r.POST("/login", func(c *gin.Context) {
		const op errors.Op = "login"

		f := &requests.LoginForm{}
		err := errors.NilOrNew(op, errors.KindBadForm, c.ShouldBindJSON(f))

		err.Meta().Set(HeaderServerVersion, version)
		err.Meta().Set(HeaderClientName, c.Request.Header.Get(HeaderClientName))
		err.Meta().Set(HeaderClientVersion, c.Request.Header.Get(HeaderClientVersion))

		if err != nil {
			// TODO(dario) add logger
			c.JSON(err.Kind().StatusCode(), err)
			return
		}

		c.JSON(http.StatusOK, requests.LoginResponse{
			Expiry: time.Now().Add(time.Hour).UnixNano(),
			Token:  f.Username,
		})
	})

	api := r.Group("/api")

	api.Use(authorizationMiddleware)

	api.POST("/word", func(c *gin.Context) {
		const op errors.Op = "save.word"

		c.Status(http.StatusOK)
	})

	api.POST("/picture", func(c *gin.Context) {
		const op errors.Op = "save.picture"

		c.Status(http.StatusOK)
	})

	r.Run()
}

func isPreflightRequest(c *gin.Context) bool {
	return c.Request.Method == http.MethodOptions &&
		c.Request.Header.Get("Access-Control-Request-Method") != "" &&
		c.Request.Header.Get("Access-Control-Request-Headers") != ""
}

func getAuthorizationBearerToken(c *gin.Context) (string, bool) {
	f := strings.Fields(c.Request.Header.Get("Authorization"))
	if len(f) != 2 || f[0] != "bearer" {
		return "", false
	}
	return f[1], true
}

func enforceContentTypeMiddleware(c *gin.Context) {
	const op errors.Op = "middleware.enforce-content-type"

	// TODO(dario) this should probably be more lenient
	if c.Request.Header.Get("Accept") != gin.MIMEJSON {
		err := errors.NewFromString(
			op,
			errors.KindBadRequest,
			"Accept header must be set to "+gin.MIMEJSON,
		)
		c.JSON(err.Kind().StatusCode(), err)
		c.Abort()
		return
	}
	c.Header("Content-Type", gin.MIMEJSON)
	c.Next()
}

func versionMiddleware(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(HeaderServerVersion, version)
		c.Next()
	}
}

func authorizationMiddleware(c *gin.Context) {
	const op errors.Op = "middleware.authorization"

	if isPreflightRequest(c) {
		c.Next()
		return
	}

	if _, ok := getAuthorizationBearerToken(c); ok {
		// TODO(dario) implement authorization
		c.Next()
		return
	}

	err := errors.NewFromString(
		op,
		errors.KindBadRequest,
		"Authorization header is missing",
	)

	c.JSON(err.Kind().StatusCode(), err)
	c.Abort()
}
