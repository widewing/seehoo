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
	bytes,err := ioutil.ReadFile(imageHome+"/"+hashtag+".json")
	if err != nil { return nil,err }
	err = json.Unmarshal(bytes,&image)
	if err != nil {return nil,err}
	return &image,nil
}

func loadConfig(container *container, image *image) (*config,error) {
	return nil,nil
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
	return &container,nil
}

func newContainer(topImageHashtag string) string {
	return ""
}

func cloneContainer(container *container) string {
	return ""
}