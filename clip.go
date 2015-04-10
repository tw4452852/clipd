package main

import (
	"log"
	"time"

	"github.com/atotto/clipboard"
)

type clip struct {
	current      string
	informRemote chan string
	updateLocal  chan req
}

type req struct {
	content string
	err     chan error
}

func NewClip() *clip {
	cur, err := clipboard.ReadAll()
	if err != nil {
		log.Println(err)
	}
	c := &clip{
		current:      cur,
		informRemote: make(chan string),
		updateLocal:  make(chan req),
	}
	go c.loop()
	return c
}

func (c *clip) loop() {
	fresh := time.Tick(1 * time.Second)
	for {
		select {
		case <-fresh:
			c.fresh()
		case req := <-c.updateLocal:
			c.update(req)
		}
	}
}

func (c *clip) update(req req) {
	var err error
	defer func() {
		if err == nil {
			c.current = req.content
		}
		req.err <- err
	}()

	if c.current == req.content {
		if debug {
			log.Printf("update: same content[%q], skip\n", req.content)
		}
		return
	}

	if debug {
		log.Printf("update: [%q] -> [%q]\n", c.current, req.content)
	}
	err = clipboard.WriteAll(req.content)
}

func (c *clip) fresh() {
	now, err := clipboard.ReadAll()
	if err != nil {
		log.Printf("fresh: read failed: %s\n", err)
		return
	}
	if now == c.current {
		if debug {
			log.Printf("fresh: no change[%q]\n", now)
		}
		return
	}
	if debug {
		log.Printf("fresh: new change: [%q] -> [%q]\n", c.current, now)
	}
	c.current = now
	c.informRemote <- now
}

func (c *clip) Inform() <-chan string {
	return c.informRemote
}

func (c *clip) Update(s string) error {
	req := req{
		content: s,
		err:     make(chan error, 1),
	}
	c.updateLocal <- req
	return <-req.err
}
