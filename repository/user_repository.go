package repository

import (
	"database/sql"
	"errors"
	"project_rentalmobil/model"
)

type userRepository struct {
	db *sql.DB
}

// Interface UserRepository
type UserRepository interface {
	Create(user model.UserCredential) (model.UserCredential, error)
	List() ([]model.UserCredential, error)
	FindById(id uint32) (model.UserCredential, error)
	FindByUsernamePassword(username, password string) (model.UserCredential, error)
}

// Constructor untuk UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Implementasi metode Create
func (u *userRepository) Create(user model.UserCredential) (model.UserCredential, error) {
	var userId uint32
	err := u.db.QueryRow(
		"INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Password, user.Role,
	).Scan(&userId)

	if err != nil {
		return model.UserCredential{}, err
	}

	user.Id = userId
	return user, nil
}

// Implementasi metode List
func (u *userRepository) List() ([]model.UserCredential, error) {
	rows, err := u.db.Query("SELECT id, username, password, role FROM users")
	if err != nil {
		return nil, err
	}

	var users []model.UserCredential
	for rows.Next() {
		var user model.UserCredential
		if err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Implementasi metode FindById
func (u *userRepository) FindById(id uint32) (model.UserCredential, error) {
	var user model.UserCredential
	err := u.db.QueryRow(
		"SELECT id, username, password, role FROM users WHERE id = $1",
		id,
	).Scan(&user.Id, &user.Username, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserCredential{}, errors.New("user not found")
		}
		return model.UserCredential{}, err
	}

	return user, nil
}

// Implementasi metode FindByUsernamePassword
func (u *userRepository) FindByUsernamePassword(username, password string) (model.UserCredential, error) {
	var user model.UserCredential
	err := u.db.QueryRow(
		"SELECT id, username, password, role FROM users WHERE username = $1 AND password = $2",
		username, password,
	).Scan(&user.Id, &user.Username, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserCredential{}, errors.New("invalid username or password")
		}
		return model.UserCredential{}, err
	}

	return user, nil
}
