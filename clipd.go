package main

import (
	"flag"
	"log"
	"os"
	"runtime"
)

func main() {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(runtime.NumCPU()))

	port := flag.Int("p", 12345, "port")
	serverAddr := flag.String("s", "", "server address")
	flag.Parse()
	if *serverAddr == "" {
		flag.Usage()
		os.Exit(2)
	}

	informRemote, updateFromRemote := make(chan string), make(chan string)
	go StartEndpoint(*serverAddr, *port, informRemote, updateFromRemote)

	localClip := NewClip()
	for {
		select {
		case r := <-updateFromRemote:
			localClip.Update(r)
		case s := <-localClip.Inform():
			log.Printf("get from local [%q], inform remote\n", s)
			informRemote <- s
		}
	}
}
