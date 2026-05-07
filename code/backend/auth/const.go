// Package auth provides authentication and authorization functionality
package auth

// Access and refresh token expiration times
const (
	MaxCookieAge       = 30 * 60           // 30 minutes
	MaxRefreshTokenAge = 30 * 24 * 60 * 60 // 30 days

	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
	XSRFTokenCookieName    = "XSRF-TOKEN" // #nosec G101
)
