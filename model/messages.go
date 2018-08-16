package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
	"time"
	"vanekt/test-social-network-api/entity"
)

type MessagesModel struct {
	db     *sqlx.DB
	logger *logging.Logger
}

func NewMessagesModel(db *sqlx.DB, logger *logging.Logger) *MessagesModel {
	return &MessagesModel{db, logger}
}

func (m *MessagesModel) GetDialogsByUserId(userId int) ([]entity.Dialog, error) {
	var dialogs []entity.Dialog
	q := `select d.peer, d.last_message_datetime, u.fullname, u.image
		  from dialogs d join users u on d.peer = u.id
		  where d.uid = ? order by d.last_message_datetime desc`
	err := m.db.Select(&dialogs, q, userId)
	return dialogs, err
}

func (m *MessagesModel) CreateMessage(msg *entity.Message) (message *entity.Message, err error) {
	msg.Datetime = uint32(time.Now().Unix())
	result, err := m.db.Exec(`insert into messages (datetime, text, author_id) values (?, ?, ?)`, msg.Datetime, msg.Text, msg.AuthorId)
	messageID, err := result.LastInsertId()
	if err != nil {
		return
	}

	query := `insert into message_entries (uid, peer, message_id) values (?, ?, ?)`
	msg.Id = uint32(messageID)

	var args []interface{}
	args = append(args, msg.AuthorId, msg.PeerId, msg.Id)
	if msg.AuthorId != msg.PeerId {
		query += `, (?, ?, ?)`
		args = append(args, msg.PeerId, msg.AuthorId, msg.Id)
	}

	_, err = m.db.Exec(query, args...)

	message = msg

	return
}
