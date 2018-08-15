package entity

type Message struct {
	Id       uint32 `db:"id" json:"id"`
	Datetime uint32 `db:"datetime" json:"datetime"`
	Text     string `db:"text" json:"text"`
	AuthorId uint32 `db:"author_id" json:"author_id"`
	PeerId   uint32 `db:"peer" json:"peer"`
}
