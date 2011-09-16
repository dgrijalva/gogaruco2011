package main

import (
	"fmt"
	"json"
	"io/ioutil"
	"container/list"
	"flag"
	"http"
	"websocket"
	"time"

	// Don't do this.  Ugly shortcut for presentation.
	"../twitter/_go_"
)

type account struct {
	Username string "username"
	Password string "password"
}

var http_port int
var ws_channels *list.List
var template []byte

func main() {
	// Parse args
	flag.IntVar(&http_port, "port", 3100, "port to listen on for http")
	flag.Parse()

	// Prep channel list
	ws_channels = list.New()

	var auth = new(account)
	// read config file
	if confData, err := ioutil.ReadFile("../twitter/account.json"); err == nil {
		if err := json.Unmarshal(confData, auth); err != nil {
			fmt.Println("Error parsing config file:", err)
		}
	} else {
		fmt.Println("Error reading config file:", err)
		return
	}

	// print stream
	if t, err := twitter.NewStream(auth.Username, auth.Password); err == nil {
		// Start processing events
		go process(t)

		// Start server
		http.HandleFunc("/", RootServer)
		http.Handle("/events", websocket.Handler(SocketServer))
		err := http.ListenAndServe(fmt.Sprintf(":%v", http_port), nil)
		if err != nil {
			fmt.Println("ListenAndServe: ", err.String())
		}

	} else {
		fmt.Println(err)
	}
}

func RootServer(w http.ResponseWriter, req *http.Request) {
	if template, err := ioutil.ReadFile("index.html"); err == nil {
		w.Write(template)
	} else {
		w.WriteHeader(500)
		w.Write([]byte("Could not open template file"))
	}
}

func process(t twitter.Stream) {
	// t = twitter.NewDeduper(t)
	for {
		if u, ok := <-t.C(); ok {
			go notify(u)
		} else {
			return
		}
	}
}

func notify(data twitter.Update) {
	// multiplex events out to all connected users
	for e := ws_channels.Front(); e != nil; e = e.Next() {
		if enc, err := json.Marshal(data); err == nil {
			e.Value.(chan []byte) <- enc
		}
	}
}

func SocketServer(ws *websocket.Conn) {
	c := make(chan []byte, 100)
	e := ws_channels.PushBack(c)
	fmt.Printf("New connection:    %v total\n", ws_channels.Len())
	var data []byte
	for {
		select {
		case data = <-c:
		case <-time.After(5e9): // make sure we're still connected
			data = []byte("")
		}

		if _, err := ws.Write(data); err != nil {
			// fmt.Println("Closing")
			ws.Close()
			break
		}
	}
	ws_channels.Remove(e)
	fmt.Printf("Closed connection: %v total\n", ws_channels.Len())
}
