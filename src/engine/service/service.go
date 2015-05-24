package service

import (
    "fmt"
    "log"
    "net"
    "strings"
    "bufio"
)

var stopServiceChan chan string = make(chan string)

func Start() {
    socket, err := net.Listen("tcp", "127.0.0.1:7777")
    if err!=nil {
    	log.Println("Cannot start TCP server")
    	log.Println(err.Error())
    	return
    }
    defer func(){
    	socket.Close()
    	log.Println("TCP server stopped")
    }()
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

var functions map[string]func([]string) = make(map[string]func([]string))

func handleSession(conn net.Conn) {
    defer func() {
    	conn.Close()
    }()
    scanner := bufio.NewScanner(conn)
    for {
    	scanner.Scan()
    	line := scanner.Text()
    	fmt.Println(line)
    	parts := strings.Fields(line)
    	cmd := parts[0]
    	args := parts[1:]
    	if cmd=="exit" {
    		return
    	} else if cmd=="stop" {
    		stopServiceChan<-"stop"
    		return
    	} else {
    		if fn,found := functions[cmd];found{
    			fn(args)
    		} else {
    			log.Println("Undefined command: "+cmd)
    		}
    	}
    }
}

func registerCommand(cmd string,fn func([]string)) bool{
	functions[cmd] = fn
	return true
}