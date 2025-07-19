package model

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID          []byte                `json:"id"`
	Name        string                `json:"name"`
	DisplayName string                `json:"displayName"`
	Credentials []webauthn.Credential `json:"credentials"`
}

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	return u.ID
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.Name
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnIcon returns the user's icon URL (optional)
func (u User) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
