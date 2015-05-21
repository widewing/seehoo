package container
import (
	"time"
)

type container struct {
	id string
	name string
	createTime time.Time
	lastStartTime time.Time
	status int
	images []image
	configs []config
}

type image struct {
	name string
	hashtag string
	path string
	parentHashTag string
	configScript string
	shell string
	startScript string
	configKeys []string
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