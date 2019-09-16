package main

import (
	"crypto/rand"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/qwwqe/VocabKing/pkg/errors"
	"github.com/qwwqe/VocabKing/pkg/requests"
)

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const (
	HeaderClientName    = "X-Client-Name"
	HeaderClientVersion = "X-Client-Version"
	HeaderServerVersion = "X-Server-Version"
)

var (
	version = "unset"

	jwtSigningKey = flag.String("jwt-signing-key", "", "JWT signing key")
	isDebugMode   = flag.Bool("debug", false, "run in debug mode")
)

func main() {
	flag.Parse()

	if *isDebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if *jwtSigningKey == "" {
		if !*isDebugMode {
			log.Fatal("JWT signing key was not provided. Exiting.")
		}

		log.Println("JWT signing key is empty. Generating random key.")
		key := make([]byte, 64)
		rand.Read(key)
		*jwtSigningKey = string(key)
	}

	r := gin.Default()

	r.Use(
		versionMiddleware(version),
		enforceContentTypeMiddleware(gin.MIMEJSON),
	)

	if *isDebugMode {
		r.Use(cors.Default())
	} else {
		r.Use(cors.New(cors.Config{
			AllowOrigins:  strings.Fields(os.Getenv("ALLOWED_ORIGINS")),
			AllowMethods:  []string{"GET", "POST"},
			AllowHeaders:  []string{HeaderClientName, HeaderClientVersion},
			ExposeHeaders: []string{HeaderServerVersion},
		}))
	}

	auth := r.Group("/auth")

	auth.POST("/login", func(c *gin.Context) {
		const op errors.Op = "login"

		meta := errors.Meta{
			HeaderServerVersion: version,
			HeaderClientName:    c.Request.Header.Get(HeaderClientName),
			HeaderClientVersion: c.Request.Header.Get(HeaderClientVersion),
		}

		f := &requests.LoginForm{}

		if err := c.ShouldBindJSON(f); err != nil {
			err := errors.NilOrNew(op, errors.KindBadForm, err, meta)
			c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
			return
		}

		// TODO(dario) implement authentication

		expireAt := time.Now().Add(5 * time.Minute).Unix()

		claim := &claims{
			Username: f.Data.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expireAt,
			},
		}

		t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).
			SignedString([]byte(*jwtSigningKey))

		if err != nil {
			err := errors.New(op, errors.KindInternalError, err, meta)
			c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, requests.NewLoginResponse(expireAt, t))
	})

	auth.POST("/refresh", func(c *gin.Context) {
		// TODO(dario) implement token refresh endpoint
		c.Status(http.StatusOK)
	})

	api := r.Group("/api")

	api.Use(authorizationMiddleware)

	api.POST("/word", func(c *gin.Context) {
		const op errors.Op = "save.word"

		meta := errors.Meta{
			HeaderServerVersion: version,
			HeaderClientName:    c.Request.Header.Get(HeaderClientName),
			HeaderClientVersion: c.Request.Header.Get(HeaderClientVersion),
		}

		f := &requests.SaveWordForm{}

		if err := c.ShouldBindJSON(f); err != nil {
			err := errors.NilOrNew(op, errors.KindBadForm, err, meta)
			c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, requests.NewSaveWordResponse())
	})

	api.POST("/picture", func(c *gin.Context) {
		const op errors.Op = "save.picture"

		meta := errors.Meta{
			HeaderServerVersion: version,
			HeaderClientName:    c.Request.Header.Get(HeaderClientName),
			HeaderClientVersion: c.Request.Header.Get(HeaderClientVersion),
		}

		f := &requests.SavePictureForm{}

		if err := c.ShouldBindJSON(f); err != nil {
			err := errors.NilOrNew(op, errors.KindBadForm, err, meta)
			c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, requests.NewSavePictureResponse())
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
	if len(f) != 2 || f[0] != "Bearer" {
		return "", false
	}
	return f[1], true
}

func enforceContentTypeMiddleware(t string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op errors.Op = "middleware.enforce-content-type"

		// TODO(dario) this should probably be more lenient
		if c.Request.Header.Get("Accept") != t {
			err := errors.NewFromString(
				op,
				errors.KindBadRequest,
				"Accept header must be set to "+t,
				nil,
			)
			c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
			c.Abort()
			return
		}

		c.Header("Content-Type", t)
		c.Next()
	}
}

func versionMiddleware(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(HeaderServerVersion, version)
		c.Next()
	}
}

func authorizationMiddleware(c *gin.Context) {
	const op errors.Op = "middleware.authorization"

	meta := errors.Meta{
		HeaderServerVersion: version,
		HeaderClientName:    c.Request.Header.Get(HeaderClientName),
		HeaderClientVersion: c.Request.Header.Get(HeaderClientVersion),
	}

	if isPreflightRequest(c) {
		c.Next()
		return
	}

	if t, ok := getAuthorizationBearerToken(c); ok {
		claim := &claims{}
		t, err := jwt.ParseWithClaims(t, claim, func(_ *jwt.Token) (interface{}, error) {
			return []byte(*jwtSigningKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				err := errors.New(op, errors.KindInvalidToken, err, meta)
				c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
				c.Abort()
				return
			}

			err := errors.New(op, errors.KindBadRequest, err, meta)
			c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
			c.Abort()
			return
		}

		if !t.Valid {
			err := errors.New(op, errors.KindInvalidToken, nil, meta)
			c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
			c.Abort()
			return
		}

		c.Next()
		return
	}

	err := errors.NewFromString(op, errors.KindBadRequest, "authorization header is missing", meta)
	c.JSON(err.Kind().StatusCode(), requests.NewErrorResponse(err))
	c.Abort()
}
