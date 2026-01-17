package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/senatroxx/filmix-backend/internal/database/entities"
)

type IUserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetRoleByName(ctx context.Context, name string) (*entities.Role, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, name, email, password, role_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Name, user.Email, user.Password, user.RoleID,
	)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, name, email, password, role_id FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, name, email, password, role_id FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.RoleID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetRoleByName(ctx context.Context, name string) (*entities.Role, error) {
	role := &entities.Role{}
	query := `SELECT id, name FROM roles WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return role, nil
}
