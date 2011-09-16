package main

import (
	"http"
	"os"
	"fmt"
	"flag"
)

var http_port = flag.Int("port", 8001, "Port to listen on")

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
			w.Write(buf[0:bytes])
		}
	} else {
		w.WriteHeader(500)
		w.Write([]byte("Could not open data file"))
	}
}
