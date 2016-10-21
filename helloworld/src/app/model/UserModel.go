package model

type UserModel struct {
Table    string  `db:"user" json:"-"`
Connect  string `db:"default" json:"-"`

Id       int64 `name:"id" pk:"auto" json:"id"`
Password string `name:"password" json:"password,omitempty"`
UserName string `name:"username" json:"username"`

CreateAt string `name:"create_at" auto:"time,insert" json:"create_at"`
UpdateAt string `name:"update_at" auto:"time,update|insert" json:"update_at"`

Token    string `json:"token,omitempty"`
}
    