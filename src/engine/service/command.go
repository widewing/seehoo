package service

import (
	"engine/container"
)

var registered bool = registerAll()
func registerAll() bool{
	registerCommand("container-start",func(args []string){
		container.Start(args[0])
	})
	registerCommand("container-stop",func(args []string){
		container.Stop(args[0])
	})
	return true
}