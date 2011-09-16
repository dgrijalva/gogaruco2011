package main

import (
	"http"
	"os"
	"fmt"
	"flag"
	"time"
)

// To get a sample data set:
// curl https://stream.twitter.com/1/statuses/sample.json -uUser:Pass > sample_data.json
// Make sure you don't have a half-complete event at the end of your file or
// it won't loop correctly.

var http_port = flag.Int("port", 8001, "Port to listen on")
var rate = flag.Int64("rate", 30, "entries per second")

func main() {
	flag.Parse()

	// Start server
	http.HandleFunc("/", TweetServer)
	err := http.ListenAndServe(fmt.Sprintf(":%v", *http_port), nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err.String())
	}

}

func TweetServer(w http.ResponseWriter, req *http.Request) {
	if file, err := os.Open("sample_data.json"); err == nil {
		defer file.Close()
		buf := make([]byte, 1024)
		var bytes int
		for err == nil {
			bytes, err = file.Read(buf)
			if err == os.EOF {
				// repeat the file forever
				file.Seek(0, 0)
				err = nil
			}
			_, err = w.Write(buf[0:bytes])
			// send rate chunks per second.  this doesn't 
			// map 1:1 with entries, but it doesn't matter
			// for this demo
			time.Sleep(1e9 / *rate)
		}
	} else {
		w.WriteHeader(500)
		w.Write([]byte("Could not open data file"))
	}
}
