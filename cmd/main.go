package main

import (
	"ev3brick"
	"log"
	"time"

	"github.com/ev3go/ev3"
)

func main() {
	ev3.LCD.Init(true)
	defer ev3.LCD.Close()

	b := ev3brick.Brick
	err := b.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", *ev3brick.Brick)

	log.Println("Run for time")
	err = b.RunForTime(ev3brick.OutPortA, 500, 5000)
	if err != nil {
		log.Println(err)
	}

	log.Println("Run for rotation")
	time.Sleep(5 * time.Second)

	err = b.RunForRotation(ev3brick.OutPortA, 200, 1)
	if err != nil {
		log.Println(err)
	}
	log.Println("MoveTank for rotation")
	time.Sleep(5 * time.Second)
	err = b.MoveTankRotation(ev3brick.OutPortA, ev3brick.OutPortD, 500, 500, 2)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(5 * time.Second)
}
