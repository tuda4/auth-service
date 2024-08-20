package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	db "github.com/tuda4/mb-backend/db/sqlc"
	"github.com/tuda4/mb-backend/internal/errors"
	"github.com/tuda4/mb-backend/internal/response"
	"github.com/tuda4/mb-backend/internal/val"
	"github.com/tuda4/mb-backend/util"
	"github.com/tuda4/mb-backend/worker"
)

type CreateAccountRequest struct {
	Email       string    `json:"email" validate:"required,email,max=256"`
	Password    string    `json:"password" validate:"required,alphanum,min=8,max=32"`
	FirstName   string    `json:"first_name" validate:"required,alphanumunicode,max=200"`
	LastName    string    `json:"last_name" validate:"required,alphanumunicode,max=200"`
	PhoneNumber string    `json:"phone_number" validate:"required,alphanum,min=8,max=20"`
	Birthday    time.Time `json:"birthday" validate:"required"`
	Address     string    `json:"address" validate:"alphanum,max=200"`
}

// CreateAccount godoc
// @Summary      Create an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        param   body      CreateAccountRequest  true  "params create account"
// @Success      200  {object}  bool
// @Router       /signup [post]
func (server *Server) createAccount(ctx echo.Context) (err error) {
	var account CreateAccountRequest

	if err = val.ValidateRequest(ctx, &account); err != nil {
		return err
	}

	arg := db.AccountTxParams{
		Email:       account.Email,
		Password:    account.Password,
		FirstName:   account.FirstName,
		LastName:    account.LastName,
		PhoneNumber: account.PhoneNumber,
		Address:     account.Address,
		BirthDay:    account.Birthday,
		AfterAccount: func(account db.Account) error {
			taskPayload := &worker.PayloadVerifyEmail{
				Email: account.Email,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskVerifyEmail(ctx.Request().Context(), taskPayload, opts...)
		},
	}

	success, err := server.store.CreateAccountTx(ctx.Request().Context(), arg)
	if err != nil {
		return response.ResponseError(ctx, http.StatusForbidden, errors.ErrorFailCreateAccount)
	}

	return response.ResponseSuccess(ctx, http.StatusOK, success)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanumunicode,min=8,max=32"`
}

type TokenResponse struct {
	Session       int64     `json:"session"`
	AccountID     uuid.UUID `json:"account_id"`
	AccessToken   string    `json:"access_token"`
	ExpiredAccess time.Time `json:"expired_token"`
	RefreshToken  string    `json:"refresh_token"`
}

func (server *Server) createToken(ctx echo.Context, accountID uuid.UUID) (rsp TokenResponse, err error) {
	accessToken, payloadAccessToken, err := server.token.CreateToken(accountID, server.config.DurationAccessToken)
	if err != nil {
		return rsp, response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorInvalidToken)
	}

	refreshToken, payloadRefreshToken, err := server.token.CreateToken(accountID, server.config.DurationRefreshToken)
	if err != nil {

		return rsp, response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorInvalidToken)
	}

	session, err := server.store.CreateSession(ctx.Request().Context(), db.CreateSessionParams{
		AccountID:    accountID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request().UserAgent(),
		ClientID:     ctx.RealIP(),
		IsBlocked:    false,
		ExpiredAt:    time.Now().Add(server.config.DurationRefreshToken),
	})
	if err != nil {

		return rsp, response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorInvalidToken)
	}

	rsp = TokenResponse{
		Session:       session.ID,
		AccountID:     payloadRefreshToken.AccountID,
		AccessToken:   accessToken,
		ExpiredAccess: payloadAccessToken.ExpiredAt,
		RefreshToken:  refreshToken,
	}
	return
}

// Login godoc
// @Summary     Login an account
// @Description  Login
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        param   body      LoginRequest  true  "params login account"
// @Success      200  {object}  response.SuccessResponse[TokenResponse] "ok"
// @Router       /login [post]
func (server *Server) login(ctx echo.Context) (err error) {
	var req LoginRequest

	if err = val.ValidateRequest(ctx, &req); err != nil {
		return err
	}

	account, err := server.store.GetAccountByEmail(ctx.Request().Context(), req.Email)
	if err != nil {
		return response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorFailEmailPassword)
	}

	err = util.CheckHashPassword(req.Password, account.HashPassword)
	if err != nil {
		return response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorFailEmailPassword)
	}

	rsp, err := server.createToken(ctx, account.AccountID)
	if err != nil {
		return
	}

	return response.ResponseSuccess(ctx, http.StatusOK, rsp)
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type RefreshTokenResponse struct {
	AccessToken   string    `json:"access_token"`
	ExpiredAccess time.Time `json:"expired_token"`
}

func (server *Server) refreshToken(ctx echo.Context) (err error) {
	var req RefreshTokenRequest
	if err = val.ValidateRequest(ctx, &req); err != nil {
		return err
	}

	session, err := server.store.GetOneSession(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		return response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorInvalidRefreshToken)
	}
	if session.ClientID != ctx.RealIP() {
		return response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorInvalidClientIP)
	}
	if session.UserAgent != ctx.Request().UserAgent() {
		return response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorInvalidUserAgent)
	}
	accessToken, payloadAccessToken, err := server.token.CreateToken(session.AccountID, server.config.DurationAccessToken)
	if err != nil {
		return response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorCreateNewAccessToken)
	}

	rsp := RefreshTokenResponse{
		AccessToken:   accessToken,
		ExpiredAccess: payloadAccessToken.ExpiredAt,
	}

	return response.ResponseSuccess(ctx, http.StatusOK, rsp)
}

type VerifyEmailRequest struct {
	AccountID  uuid.UUID `json:"account_id" validate:"required"`
	SecretCode string    `json:"secret_code" validate:"required"`
}

func (server *Server) verifyEmail(ctx echo.Context) error {
	var req VerifyEmailRequest
	if err := val.ValidateRequest(ctx, &req); err != nil {
		return err
	}

	err := server.store.VerifyEmailTx(ctx.Request().Context(), db.VerifyEmailTxParams{
		AccountID:  req.AccountID,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		return response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorInvalidSecretCode)
	}
	rsp, err := server.createToken(ctx, req.AccountID)
	if err != nil {
		return response.ResponseError(ctx, http.StatusBadRequest, errors.ErrorInvalidSecretCode)
	}

	return response.ResponseSuccess(ctx, http.StatusOK, rsp)
}
