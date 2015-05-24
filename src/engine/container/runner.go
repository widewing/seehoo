package container

import (
	log "github.com/cihub/seelog"
)

func runScript(container *container,script string,shell string){
	
}

func runStartScripts(container *container){
	log.Info("Running start scripts for container "+container.Id)
	for _,image := range container.images {
		runScript(container,image.StartScript,image.Shell)
	}
}