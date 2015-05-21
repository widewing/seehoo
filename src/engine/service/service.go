package service

import (
    "fmt"
    "log"
    "net"
    "bufio"
)

var stopServiceChan chan string = make(chan string)

func Start() {
    socket, err := net.Listen("tcp", "localhost:7777")
    defer func(){
    	socket.Close()
    	log.Println("TCP server stopped")
    }()
    if err!=nil {
    	log.Fatal("Cannot start TCP server")
    	return
    }
    log.Println("TCP server started")
    var acceptChan chan bool = make(chan bool)
    go func(){
	    for {
	        conn, err := socket.Accept()
	        acceptChan <- true
		    if err!=nil {
		    	log.Println("Accpet connection failed")
		    	continue
		    }
		    log.Println("New incoming connection")
	        go handleSession(conn)
	    }
	}()
    for {
    	stopped := false
    	select {
    		case <-acceptChan:
    		case <-stopServiceChan:
    			log.Println("Stopping TCP connetion")
    			stopped = true
    	}
    	if stopped {
    		break
    	}
    }
}

func handleSession(conn net.Conn) {
    defer func() {
    	conn.Close()
    }()
    scanner := bufio.NewScanner(conn)
    for {
    	scanner.Scan()
    	line := scanner.Text()
    	fmt.Println(line)
    	if line=="exit" {
    		return
    	}
    	if line=="stop" {
    		stopServiceChan<-"stop"
    		return
    	}
    }
}