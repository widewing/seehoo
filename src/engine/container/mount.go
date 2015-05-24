package container

import (
	"os"
	"engine/util"
	log "github.com/cihub/seelog"
)

func mountImageFs(image *image) string {
	mountPath := "/mnt/"+image.Hashtag
	imagePath := imageHome + "/" + image.Filename
	os.MkdirAll(mountPath,0755)
	log.Debug("Mounting %s on %s",imagePath,mountPath)
	switch image.ImageType {
		case "squashfs":mountSquashFs(imagePath,mountPath)
		default: return ""
	}
	image.mountPath = mountPath
	return mountPath
}

func mountSquashFs(imagePath string,mountPath string) error {
	log.Debug("mounting SquashFS")
	util.ExecuteDefaultLogger(
		"/bin/squashfuse","-o","allow_other",imagePath, mountPath)
	return nil
}

func mountConfigFs(containerId string, config *config) string {
	
	return ""
}

func mountUserFs(containerId string) string {
	mountPath := containerHome+"/"+containerId+"/data"
	os.MkdirAll(mountPath,0755)
	return mountPath
}

func mountOverlays(containerId string,paths []string) string {
	mountPath := containerHome+"/"+containerId+"/root"
	os.MkdirAll(mountPath,0755)
	args := paths[0]+"=rw"
	for _,path := range paths[1:] {
		if path=="" { continue }
		args += ":"+path
	}
	log.Debug("unionfs-fuse "+args)
	util.ExecuteDefaultLogger(
		"/bin/unionfs-fuse","-o","cow,allow_other,exec,dev",args,mountPath)

	return mountPath
}

func mountContainer(container *container) string {
	log.Info("Mounting container "+container.Id)
	lvls := len(container.images)
	paths := make([]string, lvls*2+1)
	paths[0] = mountUserFs(container.Id)
	for i,image := range container.images {
		if image.mountPath=="" {
			log.Info("Mounting image: "+image.Filename)
			paths[i*2+1] = mountImageFs(image)
		} else {
			log.Info("Image already mounted: "+image.Filename)
			paths[i*2+1] = image.mountPath
		}
		log.Info("Mounting config for "+image.Filename)
		paths[i*2+2] = mountConfigFs(container.Id,container.configs[i])
	}
	log.Info("Mounting unionfs for container "+container.Id)
	mountPath := mountOverlays(container.Id,paths)
	mountMisc(mountPath)
	return mountPath
}

func mountMisc(root string) {
	log.Info("Mounting /dev,/sys,/proc at "+root)
	os.MkdirAll(root+"/dev",0755)
	util.ExecuteDefaultLogger("/bin/busybox","mount",
		"-o","bind","/dev",root+"/dev")
	os.MkdirAll(root+"/sys",0755)
	util.ExecuteDefaultLogger("/bin/busybox","mount",
		"-t","sysfs","sysfs",root+"/sys")
	os.MkdirAll(root+"/proc",0755)
	util.ExecuteDefaultLogger("/bin/busybox","mount",
		"-t","proc","proc",root+"/proc")
	util.ExecuteDefaultLogger("/bin/busybox","mount",
		"-t","devpts","devpts",root+"/dev/pts")
}

func umountContainer(id string) {
	root := containerHome+"/"+id+"/root"
	log.Info("Unmounting /dev,/sys,/proc at "+root)
	util.Umount("/bin/busybox",root+"/dev/pts")
	util.Umount("/bin/busybox",root+"/dev")
	util.Umount("/bin/busybox",root+"/proc")
	util.Umount("/bin/busybox",root+"/sys")
	util.Umount("/bin/busybox",root)
}