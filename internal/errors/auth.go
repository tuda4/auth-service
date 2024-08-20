package errors

var (
	ErrorFailEmailPassword = &BaseError{
		Code:    1001,
		Type:    "unauthorized",
		Key:     "fail_email_password",
		Message: "Email or password fail",
	}

	ErrorUnauthorized = &BaseError{
		Code:    1002,
		Type:    "unauthorized",
		Key:     "unauthorized",
		Message: "Unauthorized",
	}

	ErrorExpiredToken = &BaseError{
		Code:    1003,
		Type:    "unauthorized",
		Key:     "expired_token",
		Message: "Token is expired",
	}

	ErrorInvalidToken = &BaseError{
		Code:    1004,
		Type:    "unauthorized",
		Key:     "invalid_token",
		Message: "Token is invalid",
	}

	ErrorInvalidRefreshToken = &BaseError{
		Code:    1005,
		Type:    "unauthorized",
		Key:     "invalid_refresh_token",
		Message: "Refresh token is invalid",
	}

	ErrorFailCreateAccount = &BaseError{
		Code:    1006,
		Type:    "unauthorized",
		Key:     "fail_create_account",
		Message: "Fail create account",
	}

	ErrorInvalidClientIP = &BaseError{
		Code:    1007,
		Type:    "unauthorized",
		Key:     "invalid_client_ip",
		Message: "Invalid client ip",
	}

	ErrorInvalidUserAgent = &BaseError{
		Code:    1008,
		Type:    "unauthorized",
		Key:     "invalid_user_agent",
		Message: "Invalid user agent",
	}

	ErrorCreateNewAccessToken = &BaseError{
		Code:    1009,
		Type:    "unauthorized",
		Key:     "fail_create_new_access_token",
		Message: "Fail create new access token",
	}

	ErrorInvalidSecretCode = &BaseError{
		Code:    1010,
		Type:    "unauthorized",
		Key:     "invalid_secret_code",
		Message: "Invalid secret code",
	}
)
