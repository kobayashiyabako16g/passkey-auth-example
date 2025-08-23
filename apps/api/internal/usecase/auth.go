package usecase

import (
	"context"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/repository"
	dtos "github.com/kobayashiyabako16g/passkey-auth-example/internal/usecase/dto/auth"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/logger"
)

type Auth interface {
	BeginRegistration(ctx context.Context, dto dtos.BeginRegistrationRequest) (*dtos.BeginRegistrationResponse, error)
	FinishRegistration(ctx context.Context, dto dtos.FinishRegistrationRequest) error
	BeginLogin(ctx context.Context, dto dtos.BeginLoginRequest) (*dtos.BeginLoginResponse, error)
	FinishLogin(ctx context.Context, dto dtos.FinishLoginRequest) error
}

type auth struct {
	sr       repository.Session
	ur       repository.User
	webAuthn *webauthn.WebAuthn
}

func NewAuth(sr repository.Session, ur repository.User, webAuthn *webauthn.WebAuthn) Auth {
	return &auth{
		sr:       sr,
		ur:       ur,
		webAuthn: webAuthn,
	}
}

func (a *auth) BeginRegistration(ctx context.Context, dto dtos.BeginRegistrationRequest) (*dtos.BeginRegistrationResponse, error) {
	// ユーザー確認
	exists, err := a.ur.ExistsByUsername(ctx, dto.Username)
	if err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		return nil, err
	}
	if exists {
		logger.Info(ctx, fmt.Sprintf("Exists User name: %s", dto.Username))
		return nil, dtos.ErrUserExists
	}

	// ユーザ作成
	var user model.User
	if err = user.GenerateID(); err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		return nil, err
	}
	user.Name = dto.Username
	user.DisplayName = dto.Username

	// チャレンジ生成
	options, sessionData, err := a.webAuthn.BeginRegistration(&user)
	if err != nil {
		logger.Error(ctx, "Error beginning registration", logger.WithError(err))
		return nil, err
	}

	// セッション作成
	session, err := a.sr.Create(ctx, user.ID)
	if err != nil {
		logger.Error(ctx, "Failed to create session", logger.WithError(err))
		return nil, err
	}

	session.Username = dto.Username
	session.RegistrationData = sessionData

	// Store に保存
	err = a.sr.Save(ctx, session)
	if err != nil {
		logger.Error(ctx, "Failed to store challenge", logger.WithError(err))
		return nil, err
	}
	return &dtos.BeginRegistrationResponse{Cred: options, Session: session}, nil
}

func (a *auth) FinishRegistration(ctx context.Context, dto dtos.FinishRegistrationRequest) error {
	session, err := a.sr.Get(ctx, dto.Session)
	if err != nil {
		logger.Error(ctx, "can't get session", logger.WithError(err))
		return err
	}
	if session == nil || session.RegistrationData == nil {
		logger.Info(ctx, "session is nil")
		return dtos.ErrSessionNotFound
	}

	// ユーザー確認
	exists, err := a.ur.ExistsByUsername(ctx, session.Username)
	if err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		return err
	}
	if exists {
		logger.Info(ctx, fmt.Sprintf("Exists User name: %s", session.Username))
		return dtos.ErrUserExists
	}

	var user model.User
	user.ID = dto.Session
	user.Name = session.Username
	user.DisplayName = session.Username

	credential, err := a.webAuthn.FinishRegistration(&user, *session.RegistrationData, dto.Request)
	if err != nil {
		logger.Error(ctx, "can't finish registration", logger.WithError(err))
		return dtos.ErrFinishRegistration
	}

	user.AddCredential(*credential)
	a.ur.Create(ctx, &user)
	// Delete the session data
	a.sr.Delete(ctx, session)

	return nil
}

func (a *auth) BeginLogin(ctx context.Context, dto dtos.BeginLoginRequest) (*dtos.BeginLoginResponse, error) {

	// ユーザー確認
	user, err := a.ur.FindByUsername(ctx, dto.Username)
	if err != nil {
		logger.Error(ctx, "can't get user", logger.WithError(err))
		return nil, err
	}
	if user == nil {
		logger.Error(ctx, "user not found")
		return nil, dtos.ErrUserNotFound
	}
	logger.Info(ctx, fmt.Sprintf("user credential: %v", len(user.Credentials)))

	// webauthn
	options, sessionData, err := a.webAuthn.BeginLogin(user)
	if err != nil {
		logger.Error(ctx, "can't begin login", logger.WithError(err))
		return nil, err
	}

	// Session確認
	session, err := a.sr.Get(ctx, dto.Session)
	if err != nil {
		logger.Error(ctx, "can't get session", logger.WithError(err))
		return nil, err
	}
	if session == nil {
		session, err = a.sr.Create(ctx, user.ID)
		if err != nil {
			logger.Error(ctx, "can't create session", logger.WithError(err))
			return nil, err
		}
	}

	session.Username = dto.Username
	session.AuthenticationData = sessionData

	if err := a.sr.Save(ctx, session); err != nil {
		logger.Error(ctx, "can't save session", logger.WithError(err))
		return nil, err
	}

	return &dtos.BeginLoginResponse{
		Cred:    options,
		Session: session,
	}, nil
}

func (a *auth) FinishLogin(ctx context.Context, dto dtos.FinishLoginRequest) error {

	// Session確認
	session, err := a.sr.Get(ctx, dto.Session)
	if err != nil {
		logger.Error(ctx, "can't get session", logger.WithError(err))
		return err
	}
	if session == nil || session.AuthenticationData == nil {
		logger.Error(ctx, "usecase: session not found")
		return dtos.ErrSessionNotFound
	}

	// user確認
	user, err := a.ur.FindByUsername(ctx, session.Username)
	if err != nil {
		logger.Error(ctx, "can't find user", logger.WithError(err))
		return err
	}

	_, err = a.webAuthn.FinishLogin(user, *session.AuthenticationData, dto.Request)
	if err != nil {
		logger.Error(ctx, "can't finish login", logger.WithError(err))
		return err
	}

	// success
	session.Authenticated = true
	session.AuthenticationData = nil

	if err := a.sr.Save(ctx, session); err != nil {
		logger.Error(ctx, "can't save session", logger.WithError(err))
		return err
	}

	return nil
}
