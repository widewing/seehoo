package container

import (
	log "github.com/cihub/seelog"
	"os/exec"
	"os"
)

func runScript(container *container,script string,shell string){
	log.Info("Running script: "+script)
	p := exec.Command("/bin/busybox","chroot",container.rootPath,shell,"-")
	p.Stdin,_ = os.Open(script)
	p.Run()
}

func runStartScripts(container *container){
	log.Info("Running start scripts for container "+container.Id)
	for _,image := range container.images {
		runScript(container,image.startScript,image.Shell)
	}
}

func killProcs(container *container) {
	log.Info("Killing processes in container "+container.Id)
}