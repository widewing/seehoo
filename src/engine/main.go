package main

import (
    "os"
    "engine/service"
    "syscall"
    "engine/container"
    "engine/config"
    "engine/util"
    log "github.com/cihub/seelog"
)

var stopChan chan string = make(chan string)

func main() {
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")
	if err == nil {
		log.ReplaceLogger(logger)
	} else {
		log.Warn("Logger not configured. Using default logger")
	}
	defer log.Flush()
	log.Info("Seehoo engine started")
	jailSelf()
	defer func(){
		cleanup()
	}()
	service.Start()
}

func jailSelf() {
	os.MkdirAll(config.CONFIG_HOMEDIR+"/dev",0755)
	os.MkdirAll(config.CONFIG_HOMEDIR+"/proc",0755)
	os.MkdirAll(config.CONFIG_HOMEDIR+"/sys",0755)
	util.ExecuteLogger(config.CONFIG_HOMEDIR+"/bin/busybox",
		"mount","-o","bind","/dev",config.CONFIG_HOMEDIR+"/dev")
	util.ExecuteLogger(config.CONFIG_HOMEDIR+"/bin/busybox",
		"mount","-t","sysfs","sysfs",config.CONFIG_HOMEDIR+"/sys")
	util.ExecuteLogger(config.CONFIG_HOMEDIR+"/bin/busybox",
		"mount","-t","proc","proc",config.CONFIG_HOMEDIR+"/proc")
	util.ExecuteLogger(config.CONFIG_HOMEDIR+"/bin/busybox",
		"mount","-t","devpts","devpts",config.CONFIG_HOMEDIR+"/dev/pts")
	os.Symlink("/proc/self/fd/0",config.CONFIG_HOMEDIR+"/dev/stdin")
	os.Symlink("/proc/self/fd/1",config.CONFIG_HOMEDIR+"/dev/stdout")
	os.Symlink("/proc/self/fd/2",config.CONFIG_HOMEDIR+"/dev/stderr")
	syscall.Chroot(config.CONFIG_HOMEDIR)
}

func cleanup() {
	log.Info("Cleaning up...")
	container.StopAll()
	util.Umount("/sys")
	util.Umount("/proc")
	util.Umount("/dev/pts")
	util.Umount("/dev")
}