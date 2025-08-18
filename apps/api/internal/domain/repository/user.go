package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/db"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/logger"
)

type User interface {
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Create(ctx context.Context, user *model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type userRepository struct {
	db *db.Client
}

func NewUser(db *db.Client) User {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name, display_name FROM users WHERE name = $1")
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return true, err
	}
	defer stmt.Close()

	var user model.User
	if err = stmt.QueryRowContext(ctx, username).Scan(&user.ID, &user.Name, &user.DisplayName); err != nil {

		//Not found
		if err == sql.ErrNoRows {
			logger.Info(ctx, fmt.Sprintf("repo: No Exists username: %s", username))
			return false, nil
		}
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return true, err
	}
	logger.Info(ctx, fmt.Sprintf("Exists username: %s", username))
	return true, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	// users table
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (id, name, display_name) VALUES ($1, $2, $3)")
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, user.ID, user.Name, user.DisplayName)
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}

	logger.Debug(ctx, fmt.Sprintf("Last Insert user id: %v", row))

	// credentials table
	jsonData, err := json.Marshal(user.Credentials[0])
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}

	stmt, err = r.db.PrepareContext(ctx, "INSERT INTO credentials (user_id, metadata) VALUES ($1, $2)")
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}
	defer stmt.Close()

	res, err = stmt.ExecContext(ctx, user.ID, jsonData)
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}

	row, err = res.RowsAffected()
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}

	logger.Debug(ctx, fmt.Sprintf("Last Insert session id: %v", row))

	return nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	stmt, err := r.db.PrepareContext(ctx, "SELECT id, name, display_name FROM users WHERE name = $1")
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	if err = stmt.QueryRowContext(ctx, username).Scan(&user.ID, &user.Name, &user.DisplayName); err != nil {

		//Not found
		if err == sql.ErrNoRows {
			logger.Info(ctx, fmt.Sprintf("repo: No Exists username: %s", username))
			return nil, nil
		}
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return nil, err
	}
	logger.Info(ctx, fmt.Sprintf("Exists username: %s", username))

	// credentials table select
	stmt, err = r.db.PrepareContext(ctx, "SELECT metadata FROM credentials WHERE user_id = $1")
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, user.ID)
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		logger.Info(ctx, fmt.Sprintf("repo: No Exists credential: %s", username))
		return &user, nil
	}

	logger.Info(ctx, "repo: Fetching credential")

	var credentials []webauthn.Credential
	for rows.Next() {
		var metadata []byte
		if err := rows.Scan(&metadata); err != nil {
			logger.Error(ctx, "Database Error", logger.WithError(err))
			return nil, err
		}

		var credential webauthn.Credential
		if err := json.Unmarshal(metadata, &credential); err != nil {
			logger.Error(ctx, "JSON Unmarshal Error", logger.WithError(err))
			continue // Skip this credential but continue with others
		}

		logger.Info(ctx, fmt.Sprintf("repo: Found credential: %s", credential.ID))

		credentials = append(credentials, credential)
	}

	if err := rows.Err(); err != nil { // ADDED: check for iteration errors
		logger.Error(ctx, "Database Rows Error", logger.WithError(err))
		return nil, err
	}

	user.Credentials = credentials

	logger.Info(ctx, fmt.Sprintf("Repo: Found user %s with %d credentials", username, len(credentials)))

	return &user, nil
}
