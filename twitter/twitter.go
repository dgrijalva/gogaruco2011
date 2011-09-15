package twitter

import (
	"json"
	"http"
	"os"
	"io"
)

type Update struct {
	Username string
	Text     string
}

type Stream struct {
	C chan Update
	body io.ReadCloser
}

func NewStream(username, password string)(*Stream, os.Error) {
	var s = &Stream{C: make(chan Update)}
	client := new(http.Client)
	req, _ := http.NewRequest("GET", "https://stream.twitter.com/1/statuses/sample.json", nil)
	req.SetBasicAuth(username, password)
	if res, err := client.Do(req); err == nil {
		if res.StatusCode == 200 {
			s.body = res.Body
			go s.process()
		} else {
			return nil, os.NewError(res.Status)
		}
	} else {
		return nil, err
	}
	
	return s, nil
}

type rawUser struct {
	Id int64 "id"
	ScreenName string `json:"screen_name"`
	Name string "name"
}

type rawTweet struct {
	Text string "text"
	User rawUser "user"
}

func (s *Stream) process() {
	// loop until body is closed
	decoder := json.NewDecoder(s.body)
	// var nextUpdate map[string]interface{}
	var nextUpdate rawTweet
	var err os.Error
	for {
		if err = decoder.Decode(&nextUpdate); err == nil {
			s.C <- Update{Text: nextUpdate.Text, Username: nextUpdate.User.ScreenName}
		} else {
			// some error happened. make sure the body is closed
			s.body.Close()
			// close the output chan so anything waiting on the next entry will return
			close(s.C)
			return
		}
	}
}

func (s *Stream) Close() {
	s.body.Close()
}