package main

import (
    "fmt"
    "os"
    "engine/service"
    "syscall"
    "engine/config"
    "engine/util"
)

var stopChan chan string = make(chan string)

func main() {
	fmt.Println("Hello, world")
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
	util.ExecuteDefaultLogger(config.CONFIG_HOMEDIR+"/bin/busybox",
		"mount","-o","bind","/dev",config.CONFIG_HOMEDIR+"/dev")
	util.ExecuteDefaultLogger(config.CONFIG_HOMEDIR+"/bin/busybox",
		"mount","-t","sysfs","sysfs",config.CONFIG_HOMEDIR+"/sys")
	util.ExecuteDefaultLogger(config.CONFIG_HOMEDIR+"/bin/busybox",
		"mount","-t","proc","proc",config.CONFIG_HOMEDIR+"/proc")
	util.ExecuteDefaultLogger(config.CONFIG_HOMEDIR+"/bin/busybox",
		"mount","-t","devpts","devpts",config.CONFIG_HOMEDIR+"/dev/pts")
	os.Symlink("/proc/self/fd/0",config.CONFIG_HOMEDIR+"/dev/stdin")
	os.Symlink("/proc/self/fd/1",config.CONFIG_HOMEDIR+"/dev/stdout")
	os.Symlink("/proc/self/fd/2",config.CONFIG_HOMEDIR+"/dev/stderr")
	syscall.Chroot(config.CONFIG_HOMEDIR)
}

func cleanup() {
	util.Umount("/bin/busybox","/sys")
	util.Umount("/bin/busybox","/proc")
	util.Umount("/bin/busybox","/dev/pts")
	util.Umount("/bin/busybox","/dev")
}