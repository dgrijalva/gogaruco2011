package main

import (
	"fmt"
	"json"
	"io/ioutil"
	"strings"

	// Don't do this.  Ugly shortcut for presentation.
	"../twitter/_go_"
)

type account struct {
	Username string "username"
	Password string "password"
}

func main() {
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
		for {
			if u, ok := <-t.C; ok {
				fmt.Printf("%v: %v\n", u.Username, strings.Replace(u.Text, "\n", "", -1))
			} else {
				return
			}
		}
	} else {
		fmt.Println(err)
	}
}