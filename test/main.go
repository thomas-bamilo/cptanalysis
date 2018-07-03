package main

import (
	"log"
	"time"
)

// Profile - is the memory representation of one user profile
type Profile struct {
	Name        string `json: "username"`
	Password    string `json: "password"`
	Age         int    `json: "age"`
	LastUpdated time.Time
}

func main() {

	p := Profile{
		Name:        "test",
		Password:    "pwtest",
		Age:         3,
		LastUpdated: time.Now(),
	}

	isUpdated := p.CreateOrUpdateProfile()
	log.Println(isUpdated)

}
