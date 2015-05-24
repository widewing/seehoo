package container
import (
	"time"
)

var imageHome string = "/images"
var containerHome string = "/containers"

type container struct {
	Id string	`json:"id"`
	Info string	`json:"info"`
	CreateTime time.Time	`json:"create_time"`
	LastStartTime time.Time	`json:"last_start_time"`
	TopImageHashtag string	`json:"top_image_hashtag"`
	AllConfigs []configItem	`json:"configs"`
	images []*image
	configs []*config
	status int
}

type image struct {
	Name string	`json:"name"`
	Hashtag string	`json:"hashtag"`
	Filename string	`json:"filename"`
	ImageType string `json:"type"`
	ParentHashTag string	`json:"parent_hashtag"`
	Shell string	`json:"shell"`
	ConfigScript string	`json:"config_script"`
	StartScript string	`json:"start_script"`
	ConfigKeys []string	`json:"config_keys"`
	mountPath string
}

type config struct {
	image *image
	items []configItem
	files []fileInfo
}

type configItem struct {
	key string
	value string
}

type fileInfo struct {
	content string
	path string
	mode string
	uid int
	gid int
}
