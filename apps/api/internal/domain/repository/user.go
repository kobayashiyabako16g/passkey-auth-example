package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kobayashiyabako16g/passkey-auth-example/internal/domain/model"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/db"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/logger"
)

type User interface {
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Create(ctx context.Context, user *model.User) error
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
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users (id, name, email, displayName) VALUES ($1, $2, $3, $4)")
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, user.ID, user.Name, user.Email, user.DisplayName)
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}

	row, err := res.LastInsertId()
	if err != nil {
		logger.Error(ctx, "Database Error", logger.WithError(err))
		return err
	}

	logger.Debug(ctx, fmt.Sprintf("Last Insert id: %v", row))

	return nil
}
