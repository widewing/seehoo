package container

import (

)

func mountImageFs(image image) string {
	return ""
}

func mountConfigFs(config config) string {
	return ""
}

func mountUserFs(containerId string) string {
	return ""
}

func mountOverlays(paths []string) string {
	return ""
}

func mountContainer(container *container) string {
	lvls := len(container.images)
	paths := make([]string, lvls*2+1)
	for i,image := range container.images {
		paths[i*2] = mountImageFs(image)
		paths[i*2+1] = mountConfigFs(container.configs[i])
	}
	paths[lvls*2] = mountUserFs(container.id)
	return mountOverlays(paths)
}