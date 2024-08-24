package repository

import (
	mb "github.com/gflydev/db"
	"github.com/gflydev/modules/jwt/model"
)

// UserRepository struct for queries from a User model.
// The struct is an implementation of interface IUserRepository
type UserRepository struct {
}

// findOne query that getting one model by a specific field condition.
func (q *UserRepository) findOne(field string, value any) *model.User {
	// Create an instance of User
	m := model.User{}

	// Get model and assign into `m` struct
	err := mb.GetModel(&m, field, value)

	// Return an empty model
	if err != nil {
		return nil
	}

	return &m
}

// GetUserByID query for getting one User by given ID.
func (q *UserRepository) GetUserByID(id int) *model.User {
	return q.findOne("id", id)
}

// GetUserByEmail query for getting one User by given Email.
func (q *UserRepository) GetUserByEmail(email string) *model.User {
	return q.findOne("email", email)
}

// CreateUser a query for creating a new user by given user data.
func (q *UserRepository) CreateUser(u *model.User) error {
	return mb.CreateModel(u)
}

// UpdateUser a query for updating a user by given user data.
func (q *UserRepository) UpdateUser(u *model.User) error {
	return mb.UpdateModel(u)
}

// DeleteUser a query for updating a user by given user data.
func (q *UserRepository) DeleteUser(u *model.User) error {
	return mb.DeleteModel(u)
}

// GetUserByToken query for getting one User by given Token.
func (q *UserRepository) GetUserByToken(token string) *model.User {
	return q.findOne("token", token)
}
