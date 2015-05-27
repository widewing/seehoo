package container
import (
	"time"
	"os"
)

var imageHome string = "/images"
var containerHome string = "/containers"

type container struct {
	Id string	`json:"id"`
	Info string	`json:"info"`
	CreateTime time.Time	`json:"create_time"`
	LastStartTime time.Time	`json:"last_start_time"`
	TopImageHashtag string	`json:"top_image_hashtag"`
	AllConfigs map[string]string	`json:"configs"`
	images []*image
	configs []*config
	home string
	rootPath string
	dataPath string
	status int
}

type image struct {
	Name string	`json:"name"`
	Hashtag string	`json:"hashtag"`
	Filename string	`json:"filename"`
	ImageType string	`json:"type"`
	ParentHashTag string	`json:"parent_hashtag"`
	Shell string	`json:"shell"`
	ConfigItems []configItem	`json:"configs"`
	home string
	configScript string
	startScript string
	stopScript string
	mountPath string
}

type configItem struct {
	Key string 	`json:"key"`	
	Type string `json:"type"`
	Default string `json:"default"`
	Choices []string `json:"choices"`
}

type config struct {
	image *image
	mountPath string
	items map[string]string
	files map[string]*fileInfo
}

type fileInfo struct {
	content []byte
	path string
	mode os.FileMode
	uid int
	gid int
}
