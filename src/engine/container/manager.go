package container

import (
	"engine/util"
	log "github.com/cihub/seelog"
)

type imageRefs struct {
	image *image
	refs map[string]*container
}

var images map[string]imageRefs = make(map[string]imageRefs)
var containers map[string]*container = make(map[string]*container)

func setupContainer(id string) error {
	log.Info("Starting container "+id)
	container,err := loadContainer(id)
	if err != nil {return err}
	mountContainer(container)
	runStartScripts(container)
	containers[id]=container
	for _,image := range container.images {
		assocImage(image,container)
	}
	return nil
}

func teardownContainer(id string) {
	log.Info("Stopping container "+id)
	container,exists:=containers[id]
	if !exists {
		log.Error("Container not started")
		return
	}
	killProcs(container)
	umountContainer(container)
	delete(containers,id)
	for _,image:=range container.images {
		desocImage(image,container)
	}
}

func teardownAll() {
	for _,container := range containers {
		teardownContainer(container.Id)
	}
}

func assocImage(image *image,cont *container){
	imageRef,existed := images[image.Hashtag]
	if !existed {
		imageRef = imageRefs{image:image,refs:make(map[string]*container)}
		images[image.Hashtag] = imageRef
	}
	imageRef.refs[cont.Id] = cont
}

func desocImage(image *image,container *container){
	imageRef,existed := images[image.Hashtag]
	if !existed {
		return
	}
	container,existed = imageRef.refs[container.Id]
	if !existed {
		return
	}
	delete(imageRef.refs,container.Id)
	if len(imageRef.refs) == 0 {
		log.Debug("No container is using image %s, umount %s",image.Filename,image.mountPath)
		util.Umount("/bin/busybox",image.mountPath)
		delete(images,image.Hashtag)
	}
}

func queryImage(hashTag string) *image{
	return nil
}

func queryContainer(id string) *container{
	return nil
}