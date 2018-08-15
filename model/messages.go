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

func (m *MessagesModel) GetDialogsByUserId(userId int) ([]uint32, error) {
	var peers []uint32
	err := m.db.Select(&peers, `select peer from dialogs WHERE uid = ? order by last_message_datetime desc`, userId)
	return peers, err
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
