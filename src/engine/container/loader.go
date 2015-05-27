package container

import (
	"io/ioutil"
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
)

func loadImage(hashtag string) (*image,error){
	log.Info("Loading image: "+hashtag)
	var image image
	image.home = imageHome+"/"+hashtag
	bytes,err := ioutil.ReadFile(image.home+"/image.json")
	if err != nil { return nil,err }
	err = json.Unmarshal(bytes,&image)
	if err != nil {return nil,err}
	image.configScript = image.home+"/config.sh"
	image.startScript = image.home+"/start.sh"
	image.stopScript = image.home+"/stop.sh"
	return &image,nil
}

func loadConfig(container *container, image *image) (*config,error) {
	var config config
	config.image = image
	config.items = make(map[string]string)
	for _,configItem := range image.ConfigItems {
		value,ok := container.AllConfigs[configItem.Key]
		if !ok { value = configItem.Default }
		config.items[configItem.Key] = value
	}
	return &config,nil
}

func loadContainer(id string) (*container,error){
	log.Info("Loading container "+id)
	var container container
	bytes,err := ioutil.ReadFile(containerHome+"/"+id+"/config.json")
	if err != nil { 
		log.Error("Cannot read config file for container "+id)
		return nil,err
	}
	err = json.Unmarshal(bytes,&container)
	if err != nil {
		log.Error("Cannot unmarshal config")
		return nil,err
	}
	if container.Id != id {
		log.Error("Container id mismatch!")
		return nil,errors.New("Container id mismatch")
	}
	hashtag := container.TopImageHashtag
	container.images = []*image{}
	container.configs = []*config{}
	for hashtag!="" {
		image,_ := loadImage(hashtag)
		container.images = append(container.images,image)
		config,_ := loadConfig(&container,image)
		container.configs = append(container.configs,config)
		hashtag = image.ParentHashTag
	}
	lastShell:="/bin/sh"
	for i:=len(container.images)-1;i>=0;i-- {
		if container.images[i].Shell=="" {
			container.images[i].Shell = lastShell
		}
		lastShell = container.images[i].Shell
	}
	return &container,nil
}

func newContainer(topImageHashtag string) string {
	return ""
}

func cloneContainer(container *container) string {
	return ""
}