package engine

import (
    "fmt"
)

var stopChan chan string = make(chan string)

func main() {
	fmt.Println("Hello, world")
	startService()
	
}
