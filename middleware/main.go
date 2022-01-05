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
		AdminUser.UserId = data.UserId
		AdminUser.DeviceId = data.DeviceId
		AdminUser.AccessToken = data.AccessToken
	} else {
		AdminUser = Admin{
			AccessToken: "syt_YWRtaW4_GtBdQzjSyZpmtedULOCQ_2GBY8L",
			UserId:      "@admin:localhost",
			DeviceId:    "VHQMRYRTSQ",
		}
	}
}

func main() {
	roomId, err := AdminUser.CreateRoom("syndicate", "testing admin room")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("room created successfully: ", roomId)
}

//http://localhost:30009/_matrix/client/r0/rooms/!QdsVwVmKhqRhiOGOhY%3Alocalhost/state/m.room.power_levels/
//users: {
//ban: 50
//events: {
//im.vector.modular.widgets: 50
//m.room.avatar: 50
//m.room.canonical_alias: 50
//m.room.encryption: 100
//m.room.history_visibility: 100
//m.room.name: 50
//m.room.power_levels: 100
//m.room.server_acl: 100
//m.room.tombstone: 100
//m.room.topic: 50}
//events_default: 50
//historical: 100
//invite: 0
//kick: 50
//redact: 50
//state_default: 50
//users: {@alan:localhost: 100}
//users_default: 0
//}
