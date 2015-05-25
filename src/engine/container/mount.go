package container

import (
	"os"
	"path"
	"io/ioutil"
	"engine/util"
	log "github.com/cihub/seelog"
)

func mountImageFs(image *image) {
	mountPath := "/mnt/"+image.Hashtag
	imagePath := image.home + "/" + image.Filename
	os.MkdirAll(mountPath,0755)
	log.Debugf("Mounting %s on %s\n",imagePath,mountPath)
	switch image.ImageType {
		case "squashfs":mountSquashFs(imagePath,mountPath)
		default: return
	}
	image.mountPath = mountPath
}

func mountSquashFs(imagePath string,mountPath string) error {
	log.Debug("mounting SquashFS")
	util.ExecuteDefaultLogger(
		"/bin/squashfuse","-o","allow_other",imagePath, mountPath)
	return nil
}

func mountConfigFs(container *container, config *config) {
	log.Debug("mounting Config for "+config.image.Hashtag)
	mountPath := container.home+"/config_"+config.image.Hashtag
	os.MkdirAll(mountPath,0755)
	err:=util.ExecuteDefaultLogger("/bin/busybox","mount","-t","tmpfs","tmpfs",mountPath)
	if err!=nil {return}
	for _,f:=range config.files{
		os.MkdirAll(path.Dir(mountPath+f.path),0755)
		ioutil.WriteFile(mountPath+f.path,f.content,f.mode)
		os.Chown(mountPath+f.path,f.uid,f.gid)
	}
	util.ExecuteDefaultLogger("/bin/busybox","mount","-o","remount,ro",mountPath,mountPath)
	config.mountPath = mountPath
}

func mountUserFs(container *container) {
	mountPath := container.home+"/data"
	os.MkdirAll(mountPath,0755)
	container.dataPath = mountPath
}

func mountOverlays(container *container) {
	mountPath := container.home+"/root"
	os.MkdirAll(mountPath,0755)
	args := container.dataPath+"=rw"
	for i,image := range container.images {
		args+=":"+container.configs[i].mountPath
		args+=":"+image.mountPath
	}
	log.Debug("unionfs-fuse "+args)
	util.ExecuteDefaultLogger(
		"/bin/unionfs-fuse","-o","cow,allow_other,exec,dev",args,mountPath)
	container.rootPath = mountPath
}

func mountMisc(container *container) {
	log.Info("Mounting /dev,/sys,/proc for container "+container.Id)
	os.MkdirAll(container.rootPath+"/dev",0755)
	util.ExecuteDefaultLogger("/bin/busybox","mount",
		"-o","bind","/dev",container.rootPath+"/dev")
	os.MkdirAll(container.rootPath+"/sys",0755)
	util.ExecuteDefaultLogger("/bin/busybox","mount",
		"-t","sysfs","sysfs",container.rootPath+"/sys")
	os.MkdirAll(container.rootPath+"/proc",0755)
	util.ExecuteDefaultLogger("/bin/busybox","mount",
		"-t","proc","proc",container.rootPath+"/proc")
	util.ExecuteDefaultLogger("/bin/busybox","mount",
		"-t","devpts","devpts",container.rootPath+"/dev/pts")
}

func mountContainer(container *container) {
	log.Info("Mounting container "+container.Id)
	container.home = containerHome + "/" + container.Id
	mountUserFs(container)
	for i,image := range container.images {
		log.Info("Mounting config for "+image.Filename)
		mountConfigFs(container,container.configs[i])
		if image.mountPath=="" {
			log.Info("Mounting image: "+image.Filename)
			mountImageFs(image)
		} else {
			log.Info("Image already mounted: "+image.Filename)
		}
	}
	log.Info("Mounting unionfs for container "+container.Id)
	mountOverlays(container)
	mountMisc(container)
}

func umountContainer(container *container) {
	root := containerHome+"/"+container.Id+"/root"
	log.Info("Unmounting /dev,/sys,/proc at "+root)
	util.Umount("/bin/busybox",root+"/dev/pts")
	util.Umount("/bin/busybox",root+"/dev")
	util.Umount("/bin/busybox",root+"/proc")
	util.Umount("/bin/busybox",root+"/sys")
	util.Umount("/bin/busybox",root)
	for _,config := range container.configs {
		if config.mountPath == "" { continue }
		util.Umount("/bin/busybox",config.mountPath)
	}
}