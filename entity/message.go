package entity

type Message struct {
	Id       uint32 `db:"id" json:"id"`
	Datetime uint32 `db:"datetime" json:"datetime"`
	Text     string `db:"text" json:"text"`
	AuthorId uint32 `db:"author_id" json:"authorId"`
	PeerId   uint32 `db:"peer" json:"peerId"`
}

type Dialog struct {
	PeerId   uint32 `db:"peer" json:"peerId"`
	Datetime uint32 `db:"last_message_datetime" json:"datetime"`
	Username string `db:"fullname" json:"username"`
	Image    string `db:"image" json:"image"`
}
