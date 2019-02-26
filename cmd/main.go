package main

import (
	"ev3brick"
	"log"

	"github.com/ev3go/ev3"
)

func main() {
	ev3.LCD.Init(true)
	defer ev3.LCD.Close()

	err := ev3brick.Brick.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", *ev3brick.Brick)
}
