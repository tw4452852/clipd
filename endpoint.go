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
	log.Printf("update from remote %s\n", content)
	buf <- content
	return nil
}

func StartEndpoint(serverAddr string, port int, informRemote <-chan string, updateFromRemote chan<- string) {
	err := startServer(port)
	if err != nil {
		return
	}

	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", serverAddr, port))
	if err != nil {
		log.Printf("connect to remote server[%s] failed: %s\n",
			fmt.Sprintf("%s:%d", serverAddr, port), err)
		return
	}

	var current string
	for {
		select {
		case current = <-buf:
			updateFromRemote <- current
		case s := <-informRemote:
			if s == current {
				log.Printf("remote is same as request[%q], skip it\n", s)
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
