package twitter

import (
	"json"
	"http"
	"os"
	"io"
)

type Update struct {
	Id       uint64
	Username string
	Text     string
	ImageURL string
}

type Stream interface {
	C()(<-chan Update)
	Close()
}

type RawStream struct {
	Updates chan Update
	body io.ReadCloser
}

// var FEED_URL = "https://stream.twitter.com/1/statuses/sample.json"
var FEED_URL = "http://localhost:8001/"

func NewStream(username, password string)(*RawStream, os.Error) {
	var s = &RawStream{Updates: make(chan Update, 100)}
	client := new(http.Client)
	req, _ := http.NewRequest("GET", FEED_URL, nil)
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

type rawTweet struct {
	Id uint64 "id"
	Text string "text"
	User struct{
		Id int64 "id"
		ScreenName string `json:"screen_name"`
		Name string "name"
		ImageURL string `json:"profile_image_url"`
	} "user"
}

func (s *RawStream) process() {
	// loop until body is closed
	decoder := json.NewDecoder(s.body)
	// var nextUpdate map[string]interface{}
	var nextUpdate rawTweet
	var err os.Error
	for {
		if err = decoder.Decode(&nextUpdate); err == nil {
			s.Updates <- Update{
				Id: nextUpdate.Id, 
				Text: nextUpdate.Text, 
				Username: nextUpdate.User.ScreenName, 
				ImageURL: nextUpdate.User.ImageURL,
			}
		} else {
			// some error happened. make sure the body is closed
			s.body.Close()
			// close the output chan so anything waiting on the next entry will return
			close(s.Updates)
			return
		}
	}
}

func (s *RawStream) C()(<-chan Update) {
	return s.Updates
}

func (s *RawStream) Close() {
	s.body.Close()
}