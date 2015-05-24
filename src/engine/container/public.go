package container

import (
	log "github.com/cihub/seelog"
)

func Start(id string) error {
	err := setupContainer(id)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func Stop(id string) {
	teardownContainer(id)
}

func Attach(id string,shell string) (in chan string,out chan string,err chan string) {
	return nil,nil,nil
}

func New(imageHashTag string) string {
	return ""
}

func Clone(id string, withdata bool) string {
	return ""
}