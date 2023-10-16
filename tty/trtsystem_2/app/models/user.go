package models

type User struct {
	ID            uint   `json:"id"`
	UserName      string `json:"user_name"`
	Account       string `json:"account"`
	Password      []byte `json:"-"`
	Sex           int    `json:"sex"`
	Email         string `json:"email"`
	Phonenumber   string `json:"phonenumber"`
	Type          int    `json:"type"`
	Teambelonging string `json:"teambelonging"`
}
