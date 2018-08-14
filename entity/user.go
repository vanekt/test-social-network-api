package entity

type User struct {
	Id       uint32 `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	FullName string `db:"fullname"`
	Created  string `db:"created"`
}
