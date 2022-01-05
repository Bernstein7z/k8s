package main

import (
	"log"
)

const (
	DeviceName = "Golang Service"
)

var (
	Server    = HomeServer{"http://localhost:30009"}
	Password  = "Bernstein123!"
	AdminUser = Admin{}
)

func init() {
	admin := User{
		Username: "admin",
		Password: Password,
	}

	if ok, _ := admin.IsAvailable(); ok {
		data, err := admin.Register()
		if err != nil {
			log.Fatal(err)
		}
		AdminUser = Admin(data)
	} else {
		data, err := admin.Login()
		if err != nil {
			log.Fatal(err)
		}
		AdminUser = Admin(data)
	}
}

func main() {
	var users []string
	users = []string{"@alan:localhost"}
	room, err := AdminUser.CreateRoom("syndicate", "testing admin room", &users)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("room created successfully: ", room)
}

// for update power level of a certain room
//http://localhost:30009/_matrix/client/v3/rooms/{room_id}/state/m.room.power_levels/
