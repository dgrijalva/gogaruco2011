package twitter

import (
	"json"
)

type Update struct {
	Username string
	Post     string
}

type Stream struct {
	C chan Update
}

func NewStream(username, password string) *Stream {

}
