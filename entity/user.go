package entity

type User struct {
	Id       uint32 `db:"id" json:"id"`
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"-"`
	FullName string `db:"fullname" json:"fullname"`
	Created  string `db:"created" json:"-"`
}
