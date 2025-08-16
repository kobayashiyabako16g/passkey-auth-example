package auth

import (
	"net/http"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
)

type BeginRegistrationRequest struct {
	Username string
}

type BeginRegistrationResponse struct {
	Cred    *protocol.CredentialCreation
	Session *model.Session
}

type FinishRegistrationRequest struct {
	Username string
	Session  string
	Request  *http.Request
}

type BeginLoginRequest struct {
	Username string
	Session  string
}

type BeginLoginResponse struct {
	Cred    *protocol.CredentialAssertion
	Session *model.Session
}
