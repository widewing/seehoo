package container

import (
	"os"
	"path"
	"math/rand"
	"io"
	"errors"
	"io/ioutil"
	"bufio"
	"strconv"
	"engine/util"
	"os/exec"
	"encoding/base64"
	log "github.com/cihub/seelog"
)

func mountImageFs(image *image) error {
	mountPath := "/mnt/"+image.Hashtag
	imagePath := image.home + "/" + image.Filename
	os.MkdirAll(mountPath,0755)
	log.Debugf("Mounting %s on %s\n",imagePath,mountPath)
	switch image.ImageType {
		case "squashfs":mountSquashFs(imagePath,mountPath)
		default: return errors.New("Unrecognized image type: "+image.ImageType)
	}
	image.mountPath = mountPath
	return nil
}

func mountSquashFs(imagePath string,mountPath string) error {
	log.Debug("mounting SquashFS")
	return util.ExecuteLogger(
		"/bin/squashfuse","-o","allow_other",imagePath, mountPath)
}

func getConfigFiles(config *config, image *image) error {
	log.Info("Generating config files from config script")
	log.Debug("Mounting /bin/busybox")
	err:=os.Mkdir(config.mountPath+"/bin",0755)
	if err!=nil {return err}
	defer os.Remove(config.mountPath+"/bin")
	err=util.ExecuteLogger("/bin/busybox","mount",
		"-o","bind,ro","/bin",config.mountPath+"/bin")
	if err!=nil {return err}
	defer util.Umount(config.mountPath+"/bin")
	pipe_in,pipe_out,err := os.Pipe()
	if err!=nil {return err}
	defer pipe_out.Close()
	log.Debug("Executing config script within temp env")
	p := exec.Command("/bin/busybox","chroot",config.mountPath,"/bin/busybox","sh","-")
	stdin,err := p.StdinPipe()
	if err!=nil {return err}
	stdout,err := p.StdoutPipe()
	if err!=nil {return err}
	stderr,err := p.StderrPipe()
	if err!=nil {return err}
	p.ExtraFiles = []*os.File{pipe_in}
	err = p.Start()
	if err!=nil {return err}
	token := "--"+string(rand.Uint32())
	begin_token := "begin"+token
	end_token := "end"+token
	log.Debug("Begin to write config script")
	_,err = stdin.Write([]byte("echo config script started\n"))
//	_,err = stdin.Write([]byte("for c in `/bin/busybox --list`;do alias $c=\"/bin/busybox $c\";done\n"))
//	if err!=nil {return}
//	_,err = stdin.Write([]byte(`
//		basecat() {
//			cat /base$1
//		}
//		newfile() {
//			echo `+begin_token+` >&3
//			echo new >&3
//			echo $1 >&3
//			echo $2 >&3
//			echo $3 >&3
//			echo $4 >&3
//			echo >&3
//			base64 >&3
//			echo `+end_token+` >&3
//		}
//	`))
//	if err!=nil {return}
//	file,err := os.Open(image.configScript)
//	if err!=nil {return}
//	_,err = io.Copy(stdin,file)
//	if err!=nil {return}
	stdin.Close()
	util.Wait(
		func(){util.LogStream(stdout,util.LogStdout)},
		func(){util.LogStream(stderr,util.LogStderr)},
		func(){
			scanner := bufio.NewScanner(pipe_out)
			for scanner.Scan() {
				if scanner.Text()!=begin_token { continue }
				log.Debug("Received a new file from config script")
				var fileInfo fileInfo
				scanner.Scan(); fileInfo.path = scanner.Text()
				log.Debug("filename: "+fileInfo.path)
				scanner.Scan(); mode,err := strconv.ParseInt(scanner.Text(),8,0)
				if err!=nil {continue}
				fileInfo.mode = os.FileMode(mode)
				scanner.Scan(); fileInfo.uid,err = strconv.Atoi(scanner.Text())
				if err!=nil {continue}
				scanner.Scan(); fileInfo.gid,err = strconv.Atoi(scanner.Text())
				if err!=nil {continue}
				scanner.Scan()
				base64reader,base64writer := io.Pipe()
				dec := base64.NewDecoder(base64.StdEncoding,base64reader)
				for scanner.Scan() {
					line := scanner.Text()
					if line == end_token { break }
					_,err=base64writer.Write([]byte(line))
					if err!=nil {break}
				}
				fileInfo.content,err = ioutil.ReadAll(dec)
				if err!=nil {continue}
				config.files[fileInfo.path] = &fileInfo
			}
		},
	)

	return p.Wait()
}

func mountConfigFs(container *container, config *config, image *image) error {
	log.Debug("mounting Config for "+config.image.Hashtag)
	mountPath := container.home+"/config_"+config.image.Hashtag
	os.MkdirAll(mountPath,0755)
	err:=util.ExecuteLogger("/bin/busybox","mount","-t","tmpfs","tmpfs",mountPath)
	if err!=nil {return err}
	config.mountPath = mountPath
	err = getConfigFiles(config,image)
	if err!=nil {return err}
	for _,f:=range config.files{
		os.MkdirAll(path.Dir(mountPath+f.path),0755)
		ioutil.WriteFile(mountPath+f.path,f.content,f.mode)
		os.Chown(mountPath+f.path,f.uid,f.gid)
	}
	return util.ExecuteLogger("/bin/busybox","mount","-o","remount,ro",mountPath,mountPath)
}

func mountUserFs(container *container) error {
	mountPath := container.home+"/data"
	err := os.MkdirAll(mountPath,0755)
	if err!=nil {return err}
	container.dataPath = mountPath
	return nil
}

func mountOverlays(container *container) error{
	mountPath := container.home+"/root"
	err:=os.MkdirAll(mountPath,0755)
	if err!=nil {return err}
	args := container.dataPath+"=rw"
	for i,image := range container.images {
		args+=":"+container.configs[i].mountPath
		args+=":"+image.mountPath
	}
	log.Debug("unionfs-fuse "+args)
	err = util.ExecuteLogger(
		"/bin/unionfs-fuse","-o","cow,allow_other,exec,dev",args,mountPath)
	if err!=nil {return err}
	container.rootPath = mountPath
	return nil
}

func mountMisc(container *container) error {
	log.Info("Mounting /dev,/sys,/proc for container "+container.Id)
	err:=os.MkdirAll(container.rootPath+"/dev",0755)
	if err!=nil {return err}
	err=util.ExecuteLogger("/bin/busybox","mount",
		"-o","bind","/dev",container.rootPath+"/dev")
	if err!=nil {return err}
	err=os.MkdirAll(container.rootPath+"/sys",0755)
	if err!=nil {return err}
	err=util.ExecuteLogger("/bin/busybox","mount",
		"-t","sysfs","sysfs",container.rootPath+"/sys")
	if err!=nil {return err}
	err=os.MkdirAll(container.rootPath+"/proc",0755)
	if err!=nil {return err}
	err=util.ExecuteLogger("/bin/busybox","mount",
		"-t","proc","proc",container.rootPath+"/proc")
	if err!=nil {return err}
	err=util.ExecuteLogger("/bin/busybox","mount",
		"-t","devpts","devpts",container.rootPath+"/dev/pts")
	if err!=nil {return err}
	return nil
}

func mountContainer(container *container) error {
	var err error
	log.Info("Mounting container "+container.Id)
	container.home = containerHome + "/" + container.Id
	err=mountUserFs(container)
	if err!=nil {return err}
	defer func(){if err!=nil {util.Umount(container.dataPath)}}()
	for i,image := range container.images {
		config := container.configs[i]
		if image.mountPath=="" {
			log.Info("Mounting image: "+image.Filename)
			err = mountImageFs(image)
			if err!=nil {return err}
			defer func(){if err!=nil {util.Umount(image.mountPath)}}()
		} else {
			log.Info("Image already mounted: "+image.Filename)
		}
		log.Info("Mounting config for "+image.Filename)
		err=mountConfigFs(container,config,image)
		if err!=nil {
			log.Error("Cannot mount config: "+err.Error())
		}
		defer func(){if err!=nil&&config.mountPath!="" {util.Umount(config.mountPath)}}()
	}
	log.Info("Mounting unionfs for container "+container.Id)
	err=mountOverlays(container)
	if err!=nil {return err}
	defer func(){if err!=nil {util.Umount(container.rootPath)}}()
	err=mountMisc(container)
	if err!=nil {return err}
	return nil
}

func umountContainer(container *container) {
	root := containerHome+"/"+container.Id+"/root"
	log.Info("Unmounting /dev,/sys,/proc at "+root)
	util.Umount(root+"/dev/pts")
	util.Umount(root+"/dev")
	util.Umount(root+"/proc")
	util.Umount(root+"/sys")
	util.Umount(root)
	for _,config := range container.configs {
		if config.mountPath == "" { continue }
		util.Umount(config.mountPath)
	}
}