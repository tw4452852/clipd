package main

import (
	"log"
	"time"
)

type clip struct {
	inform chan string
}

func (c *clip) Inform() <-chan string {
	return c.inform
}

func (c *clip) Update(s string) {
	log.Printf("update [%q]\n", s)
}

func NewClip() Clip {
	c := &clip{
		inform: make(chan string),
	}
	go c.loop()
	return c
}

func (c *clip) loop() {
	for range time.Tick(1 * time.Second) {
		c.inform <- "linux"
	}
}
