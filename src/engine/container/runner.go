package container

import (

)

func runScript(container *container,script string,shell string){
	
}

func runStartScripts(container *container){
	for _,image := range container.images {
		runScript(container,image.startScript,image.shell)
	}
}