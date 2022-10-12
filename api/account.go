package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	db "github.com/bank/simple-bank/db/sqlc"
	"github.com/bank/simple-bank/model"
	"github.com/bank/simple-bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// A logged-in user can only create an account for him/herself
func (server *Server) createAccount(ctx *gin.Context) {
	var req model.Account
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	/*
		The login user store in the authorization payload. Therefore in this createAccount handler
		we will get the authorization payload by calling ctx.MustGet() by pass Key and note that this
		call will return a general interface, so we have to cast it to token.Payload object.
	*/
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violoation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

/* A user should only be able to get the account that he or she owns */
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

// A logged-in user can only list accounts that belong to him/here
func (sever *Server) listAccounts(ctx *gin.Context) {
	var req model.PageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := sever.store.ListAccounts(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse((err)))
		return
	}

	var res model.Response
	res.HasNext = len(accounts) == int(req.PageSize)
	res.Page = req.PageID
	res.PageSize = req.PageSize
	res.Total = len(accounts)
	res.Result = accounts
	ctx.JSON(http.StatusOK, res)
}
