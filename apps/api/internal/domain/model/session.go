package model

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

type Session struct {
	ID                 string
	Username           string                `json:"username"`
	Authenticated      bool                  `json:"authenticated"`
	RegistrationData   *webauthn.SessionData `json:"registration_data,omitempty"`
	AuthenticationData *webauthn.SessionData `json:"authentication_data,omitempty"`
	ExpiresAt          time.Time             `json:"expires_at"`
}
