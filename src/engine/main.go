package main

import (
    "fmt"
    "engine/service"
)

var stopChan chan string = make(chan string)

func main() {
	fmt.Println("Hello, world")
	service.Start()
	
}
