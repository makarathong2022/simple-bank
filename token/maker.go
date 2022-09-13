package token

import "time"

/*
The idea is to declare a general token make interface to manage the createion and verification of the tokens.
*/
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	/* VerifyToken checks if the token is valid or not. If it is valid,
	   the method will return the payload data store inside the body of the token */
	VerifyToken(token string) (*Payload, error)
}
