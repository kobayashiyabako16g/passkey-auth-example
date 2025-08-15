package auth

import (
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
)

type BeginRegistrationRequest struct {
	Username string `json:"username"`
}

type BeginRegistrationResponse struct {
	Cred    *protocol.CredentialCreation
	Session *model.Session
}

type FinishRegistrationRequest struct {
	Username string `json:"username"`
	Session  string
	Request  *http.Request
}
