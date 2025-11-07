package repository

import (
	"context"
	"debez/internal/models"
	"debez/pkg/logger"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *UserRepository) Select(ctx context.Context, offset, limit int) ([]models.User, error) {
	sql, args, err := r.builder.Select("id", "email", "name", "last_name", "role").
		From("users").
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	defer func(){
		rows.Close()
	}()
	if err != nil {
		return nil, fmt.Errorf("select, query: %w", err)
	}
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.LastName, &user.Role); err != nil {
			return nil, fmt.Errorf("select, Scan: %w", err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("select, rows: %w", err)
	}
	return users, nil
}
func (r *UserRepository) SelectByID(ctx context.Context, id int64) (models.User, error) {
	sql, args, err := r.builder.Select("id", "email", "name", "last_name", "role").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return models.User{}, fmt.Errorf("selectByID: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	defer func(){
		rows.Close()
	}()
	if !rows.Next() {
		return models.User{}, fmt.Errorf("selectByID: user with id %d not found", id)
	}
	if err != nil {
		return models.User{}, fmt.Errorf("selectByID, query: %w", err)
	}
	var user models.User
	if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.LastName, &user.Role); err != nil {
		return models.User{}, fmt.Errorf("selectByID, Scan: %w", err)
	}
	if rows.Next() {
		return models.User{}, fmt.Errorf("selectByID: multiple users found with id %d", id)
	}
	return user, nil
}
func (r *UserRepository) Insert(ctx context.Context, user models.User) error {
	sql, args, err := r.builder.Insert("users").
		Columns("email", "name", "last_name", "role").
		Values(user.Email, user.Name, user.LastName, user.Role).
		ToSql()
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("insert, exec: %w", err)
	}
	return nil
}
func (r *UserRepository) Update(ctx context.Context, user models.User) error {
	sql, args, err := r.builder.Update("users").
		Set("email", user.Email).
		Set("name", user.Name).
		Set("last_name", user.LastName).
		Set("role", user.Role).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "failed to build update sql", zap.Error(err))
		return fmt.Errorf("update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "failed to exec update sql", zap.Error(err))
		return fmt.Errorf("update, exec: %w", err)
	}
	return nil
}
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	sql, args, err := r.builder.Delete("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("delete, exec: %w", err)
	}
	return nil
}
