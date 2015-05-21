package container

import (

)

func Start(id string) {
	container := loadContainer(id)
	mountContainer(container)
	runStartScripts(container)
}

func Stop(id string) {
	
}

func Attach(id string,shell string) (in chan string,out chan string,err chan string) {
	return nil,nil,nil
}