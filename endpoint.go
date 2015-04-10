package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var buf = make(chan string)

type Proxy struct{}

func (Proxy) Update(content string, _ *struct{}) error {
	if debug {
		log.Printf("update from remote %s\n", content)
	}
	buf <- content
	return nil
}

func StartEndpoint(serverAddr string, port int, informRemote <-chan string, updateFromRemote chan<- string) {
	err := startServer(port)
	if err != nil {
		return
	}

	var (
		client  *rpc.Client
		current string
	)
	for {
		select {
		case now := <-buf:
			if client == nil {
				client, err = rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", serverAddr, port))
				if err != nil {
					log.Printf("connect to remote server[%s] failed: %s\n",
						fmt.Sprintf("%s:%d", serverAddr, port), err)
				}
			}
			if now == current {
				log.Printf("same as previous content[%q]\n", now)
				continue
			}
			current = now
			updateFromRemote <- current
		case s := <-informRemote:
			if s == current {
				log.Printf("remote is same as request[%q], skip it\n", s)
				continue
			}
			if client == nil {
				log.Println("remote isn't up, skip it")
				continue
			}
			err := client.Call("Proxy.Update", s, nil)
			if err != nil {
				log.Printf("update remote failed: %s\n", err)
			}
		}
	}
}

func startServer(port int) error {
	rpc.Register(Proxy{})
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("listen on port %d failed: %s\n", port, err)
		return err
	}
	log.Printf("listen on port %d ...\n", port)
	go http.Serve(l, nil)
	return nil
}
