package model

import (
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

// implements: https://pkg.go.dev/github.com/go-webauthn/webauthn/webauthn#User
type User struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	DisplayName string                `json:"displayName"`
	Credentials []webauthn.Credential `json:"credentials"`
}

func NewUser(name string, displayName string) *User {
	user := &User{}
	user.GenerateID()
	user.Name = name
	user.DisplayName = displayName
	return user
}

func (u *User) GenerateID() error {
	if u.ID != "" {
		return fmt.Errorf("Exists ID")
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = uid.String()
	return nil
}

// WebAuthnID returns the user's ID
func (u *User) WebAuthnID() []byte {
	buf := []byte(u.ID)
	return buf
}

// WebAuthnName returns the user's username
func (u *User) WebAuthnName() string {
	return u.Name
}

// WebAuthnDisplayName returns the user's display name
func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnIcon is not (yet) implemented
func (u *User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.Credentials = append(u.Credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u *User) UpdateCredential(credential *webauthn.Credential) {
	for i, c := range u.Credentials {
		if string(c.ID) == string(credential.ID) {
			u.Credentials[i] = *credential
		}
	}
}
