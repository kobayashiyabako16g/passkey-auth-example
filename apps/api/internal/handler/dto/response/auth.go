package response

type BeginRegistrationResponse struct {
	Token string `json:"token"`
}

type PublicKeyCredentialCreationOptions struct {
	RP                     RelyingParty               `json:"rp"`
	User                   PublicKeyUserEntity        `json:"user"`
	Challenge              string                     `json:"challenge"` // base64url
	PubKeyCredParams       []PublicKeyCredentialParam `json:"pubKeyCredParams"`
	Timeout                int                        `json:"timeout,omitempty"`
	Attestation            string                     `json:"attestation,omitempty"`
	AuthenticatorSelection AuthenticatorSelection     `json:"authenticatorSelection,omitempty"`
}

type RelyingParty struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type PublicKeyUserEntity struct {
	ID          string `json:"id"` // base64url (should be []byte encoded)
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type PublicKeyCredentialParam struct {
	Type string `json:"type"` // always "public-key"
	Alg  int    `json:"alg"`  // e.g., -7 for ES256
}

type AuthenticatorSelection struct {
	ResidentKey      string `json:"residentKey,omitempty"`      // "required", "preferred", "discouraged"
	UserVerification string `json:"userVerification,omitempty"` // "required", "preferred", "discouraged"
}
