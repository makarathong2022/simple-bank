package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bank/simple-bank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

/*
  - gin.HandlerFunc type is a function that takes a context as input
    so, I' am gonna to return an anonymous function with the same required signature here.
    This anonymous function is in fact the authentication middleware function we want to implement
  - The purpose of this Bearer prefix is to let the server know the type of authorization, because in the reality,
    the server might support multiple of authorization schemes. Such as Oauth, Digest, AWS signature, or many more.
*/
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// authorize the user to perform the request we first check the authorization header from request
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			// mean that the client doesnot provider this header
			err := errors.New("authorization header is not provided")
			// allow us to abort the request and send the json response to the client with a specific status code
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// authorization header is provided
		// call strings.Fields func to split header by space
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// convert to ToLower it make it easier to compare
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		// store payload to context by key-pair value
		ctx.Set(authorizationPayloadKey, payload)
		// Then the last stop is to call ctx.Next() to forward the request to the next handler.
		ctx.Next()
	}
}
