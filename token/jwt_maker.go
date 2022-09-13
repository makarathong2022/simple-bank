package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token make
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be lease %d characters", minSecretKeySize)
	}
	/*
		Alright, now you can see that the red line is gone cuz our JWT struct has all 2 required function of the interface
	*/
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

/*
VerifyToken checks if the token is valid or not. If it is valid,
the method will return the payload data store inside the body of the token
*/
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrExpiredToken
		}
		// return convert key after converting it to byte slice and a nil error
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		// Convert the return error of the ParseWithClaims func to jwt.ValidationError
		verr, ok := err.(*jwt.ValidationError)
		/* Converted error to the verr variable if the convertion is ok, we use the errors.Is function
		   to check if the verr.Inner is actually the ErrExpiredToken. If it is, we just return a nil payload and the ErrEpiredToken  */
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	/*
		In case everything is good, and the token is successfully parsed and verified
	*/

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
