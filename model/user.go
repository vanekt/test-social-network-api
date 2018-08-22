package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"github.com/vanekt/test-social-network-api/entity"
)

type UserModel struct {
	db     *sqlx.DB
	logger *logging.Logger
}

func NewUserModel(db *sqlx.DB, logger *logging.Logger) *UserModel {
	return &UserModel{db, logger}
}

func (m *UserModel) GetUserById(id int) (user entity.User, err error) {
	err = m.db.Get(&user, "select * from users where id = ?", id)
	return
}

func (m *UserModel) GetUserByLogin(login string) (user entity.User, err error) {
	err = m.db.Get(&user, "select * from users where login = ?", login)
	return
}

func (m *UserModel) GetAll() (users []entity.User, err error) {
	err = m.db.Select(&users, "select * from users")
	return
}
