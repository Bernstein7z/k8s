package main

import (
	"fmt"
	"log"
)

const DeviceName = "Golang Service"

var Server = HomeServer{"http://localhost:30009"}

func main() {

	user1 := User{
		Username: "alan.bernstein", // alan.bernstein@deutschepost.de
		Password: "Bernstein123!",
	}
	if ok, err := user1.IsAvailable(); !ok {
		log.Fatal("check availability: ", err)
	}

	userData, err := user1.Register()
	if err != nil {
		log.Fatal("register: ", err)
	}

	fmt.Println("matrix: ", userData)
}
