package service

import (
	"engine/container"
)

var registered bool = registerAll()
func registerAll() bool{
	registerCommand("container-start",func(args []string)string{
		err := container.Start(args[0])
		if err == nil{
			return "ok"
		} else {
			return err.Error()
		}
	})
	registerCommand("container-stop",func(args []string)string{
		err := container.Stop(args[0])
		if err == nil{
			return "ok"
		} else {
			return err.Error()
		}
	})
	return true
}