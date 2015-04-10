package main

import (
	"flag"
	"log"
	"os"
	"runtime"
)

var (
	port       int
	serverAddr string
	debug      bool
)

func init() {
	flag.IntVar(&port, "p", 1234, "port")
	flag.StringVar(&serverAddr, "s", "", "server address")
	flag.BoolVar(&debug, "d", false, "debug mode")
}

func main() {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(runtime.NumCPU()))

	flag.Parse()
	if serverAddr == "" {
		flag.Usage()
		os.Exit(2)
	}

	informRemote, updateFromRemote := make(chan string), make(chan string)
	go StartEndpoint(serverAddr, port, informRemote, updateFromRemote)

	localClip := NewClip()
	for {
		select {
		case r := <-updateFromRemote:
			e := localClip.Update(r)
			if e != nil {
				log.Println(e)
			}
		case s := <-localClip.Inform():
			if debug {
				log.Printf("get from local [%q], inform remote\n", s)
			}
			informRemote <- s
		}
	}
}
