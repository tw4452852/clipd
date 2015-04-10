package main

import (
	"log"
	"syscall"
	"unsafe"
)

const (
	CF_UNICODETEXT = 13
)

var (
	user32           = syscall.MustLoadDLL("user32")
	openClipboard    = user32.MustFindProc("OpenClipboard")
	getClipboardData = user32.MustFindProc("GetClipboardData")
	setClipboardData = user32.MustFindProc("SetClipboardData")
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
		c.inform <- "windows"
	}
}
